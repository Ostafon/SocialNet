package main

import (
	"context"
	"google.golang.org/grpc"
	"os"
	authpb "socialnet/services/auth/gen"
	"socialnet/services/like/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"socialnet/pkg/config"
	"socialnet/services/like/internal/model"
	"socialnet/services/like/internal/repos"
	notificationpb "socialnet/services/notification/gen"
)

// =========================================================
// GLOBALS
// =========================================================

var (
	testDB   *gorm.DB
	testRepo *repos.LikeRepo
	testSvc  *service.LikeService

	ctx = context.Background()
)

// =========================================================
// MOCK NOTIFICATION CLIENT
// =========================================================

type mockNotif struct{}

func (m *mockNotif) CreateNotification(ctx context.Context, in *notificationpb.CreateNotificationRequest, opts ...grpc.CallOption) (*authpb.Confirmation, error) {
	return &authpb.Confirmation{Status: "ok"}, nil
}
func (m *mockNotif) ListNotifications(context.Context, *notificationpb.ListNotificationsRequest, ...grpc.CallOption) (*notificationpb.Notifications, error) {
	return &notificationpb.Notifications{}, nil
}
func (m *mockNotif) MarkAsRead(context.Context, *notificationpb.MarkAsReadRequest, ...grpc.CallOption) (*authpb.Confirmation, error) {
	return &authpb.Confirmation{Status: "ok"}, nil
}
func (m *mockNotif) MarkAllAsRead(context.Context, *notificationpb.EmptyRequest, ...grpc.CallOption) (*authpb.Confirmation, error) {
	return &authpb.Confirmation{Status: "ok"}, nil
}
func (m *mockNotif) DeleteNotification(context.Context, *notificationpb.DeleteNotificationRequest, ...grpc.CallOption) (*authpb.Confirmation, error) {
	return &authpb.Confirmation{Status: "ok"}, nil
}
func (m *mockNotif) ClearAll(context.Context, *notificationpb.EmptyRequest, ...grpc.CallOption) (*authpb.Confirmation, error) {
	return &authpb.Confirmation{Status: "ok"}, nil
}
func (m *mockNotif) StreamNotifications(context.Context, *notificationpb.StreamRequest, ...grpc.CallOption) (notificationpb.NotificationService_StreamNotificationsClient, error) {
	return nil, nil
}

// =========================================================
// TEST MAIN
// =========================================================

func TestMain(m *testing.M) {
	dsn := "postgres://ostafon:0000@localhost:5432/test?sslmode=disable"

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Пересоздаём таблицы
	_ = testDB.Exec(`DROP TABLE IF EXISTS likes CASCADE`)
	if err := testDB.AutoMigrate(&model.Like{}); err != nil {
		panic(err)
	}

	testRepo = repos.NewLikeRepo(testDB)

	clients := &config.GRPCClients{
		NotifClient: &mockNotif{},
	}

	testSvc = service.NewLikeService(testRepo, clients)

	os.Exit(m.Run())
}

// =========================================================
// TESTS
// =========================================================

// ---------------- TEST LIKE POST ----------------

func TestLikePost(t *testing.T) {
	// лайк — первый раз
	resp, err := testSvc.LikePost(ctx, "user1", "post123")
	assert.NoError(t, err)
	assert.Equal(t, "liked", resp.Status)
	assert.Equal(t, int32(1), resp.LikesCount)

	// лайк — повторно → дизлайк
	resp2, err := testSvc.LikePost(ctx, "user1", "post123")
	assert.NoError(t, err)
	assert.Equal(t, "unliked", resp2.Status)
	assert.Equal(t, int32(0), resp2.LikesCount)
}

// ---------------- TEST LIST POST LIKES ----------------

func TestListPostLikes(t *testing.T) {
	// добавим 2 лайка вручную
	_ = testDB.Create(&model.Like{UserID: "u1", PostID: strPtr("p1")})
	_ = testDB.Create(&model.Like{UserID: "u2", PostID: strPtr("p1")})

	resp, err := testSvc.ListPostLikes(ctx, "p1")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Likes))
}

// ---------------- TEST UNLIKE COMMENT ----------------

func TestUnlikeComment(t *testing.T) {
	// подготавливаем
	_ = testDB.Create(&model.Like{UserID: "u3", CommentID: strPtr("c1")})

	resp, err := testSvc.UnlikeComment(ctx, "u3", "c1")
	assert.NoError(t, err)
	assert.Equal(t, "unliked", resp.Status)
	assert.Equal(t, int32(0), resp.LikesCount)
}

// ---------------- TEST LIST COMMENT LIKES ----------------

func TestListCommentLikes(t *testing.T) {
	_ = testDB.Create(&model.Like{UserID: "u1", CommentID: strPtr("cx")})
	_ = testDB.Create(&model.Like{UserID: "u2", CommentID: strPtr("cx")})

	resp, err := testSvc.ListCommentLikes(ctx, "cx")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Likes))
}

// =========================================================
// HELPERS
// =========================================================

func strPtr(s string) *string {
	return &s
}
