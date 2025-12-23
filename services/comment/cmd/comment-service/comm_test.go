package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"os"
	"socialnet/services/comment/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"socialnet/pkg/config"
	authpb "socialnet/services/auth/gen"
	pb "socialnet/services/comment/gen"
	"socialnet/services/comment/internal/model"
	"socialnet/services/comment/internal/repos"
	notificationpb "socialnet/services/notification/gen"
)

// ------------------- GLOBAL -------------------

var (
	testDB   *gorm.DB
	testRepo *repos.CommentRepo
	testSvc  *service.CommentService

	ctx = context.Background()
)

// ------------------- MOCK NOTIF -------------------

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

// ------------------- TEST MAIN -------------------

func TestMain(m *testing.M) {
	dsn := "postgres://ostafon:0000@localhost:5432/test?sslmode=disable"

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// очищаем таблицы
	_ = testDB.Exec(`DROP TABLE IF EXISTS comments CASCADE`)
	if err := testDB.AutoMigrate(&model.Comment{}); err != nil {
		panic(err)
	}

	testRepo = repos.NewCommentRepo(testDB)

	// подменяем gRPC клиентов
	clients := &config.GRPCClients{
		NotifClient: &mockNotif{},
	}

	testSvc = service.NewCommentService(testRepo, clients)

	os.Exit(m.Run())
}

// ------------------- TESTS -------------------

func TestAddComment(t *testing.T) {
	req := &pb.AddCommentRequest{
		PostId:  "10",
		Content: "Hello comment",
	}

	resp, err := testSvc.AddComment(ctx, "user1", req)
	assert.NoError(t, err)
	assert.Equal(t, "user1", resp.UserId)
	assert.Equal(t, "Hello comment", resp.Content)
}

func TestGetComment(t *testing.T) {
	c := &model.Comment{
		PostID:  "20",
		UserID:  "u2",
		Content: "get me",
	}
	_ = testDB.Create(c)

	resp, err := testSvc.GetComment(ctx, fmt.Sprint(c.ID))
	assert.NoError(t, err)
	assert.Equal(t, "get me", resp.Content)
	assert.Equal(t, "u2", resp.UserId)
}

func TestDeleteComment(t *testing.T) {
	c := &model.Comment{
		PostID:  "30",
		UserID:  "u3",
		Content: "delete me",
	}
	_ = testDB.Create(c)

	err := testSvc.DeleteComment(ctx, fmt.Sprint(c.ID), "u3")
	assert.NoError(t, err)

	var count int64
	testDB.Model(&model.Comment{}).Where("id = ?", c.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestListComments(t *testing.T) {
	_ = testDB.Create(&model.Comment{PostID: "40", UserID: "u1", Content: "A"})
	_ = testDB.Create(&model.Comment{PostID: "40", UserID: "u2", Content: "B"})

	resp, err := testSvc.ListComments(ctx, "40")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Comments))
}
