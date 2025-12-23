package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"os"
	"socialnet/pkg/config"
	authpb "socialnet/services/auth/gen"
	notificationpb "socialnet/services/notification/gen"
	"socialnet/services/user/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pb "socialnet/services/user/gen"
	"socialnet/services/user/internal/model"
	"socialnet/services/user/internal/repos"
)

// ---- GLOBALS ----
var testDB *gorm.DB
var testRepo *repos.UserRepo
var testSvc *service.UserService
var ctx context.Context

// ========== FULL MOCK for NotificationServiceClient ==========
type mockNotif struct {
	Created []*notificationpb.CreateNotificationRequest
}

func (m *mockNotif) ListNotifications(ctx context.Context, in *notificationpb.ListNotificationsRequest, opts ...grpc.CallOption) (*notificationpb.Notifications, error) {
	return &notificationpb.Notifications{}, nil
}

func (m *mockNotif) MarkAsRead(ctx context.Context, in *notificationpb.MarkAsReadRequest, opts ...grpc.CallOption) (*authpb.Confirmation, error) {
	return &authpb.Confirmation{Status: "ok"}, nil
}

func (m *mockNotif) MarkAllAsRead(ctx context.Context, in *notificationpb.EmptyRequest, opts ...grpc.CallOption) (*authpb.Confirmation, error) {
	return &authpb.Confirmation{Status: "ok"}, nil
}

func (m *mockNotif) DeleteNotification(ctx context.Context, in *notificationpb.DeleteNotificationRequest, opts ...grpc.CallOption) (*authpb.Confirmation, error) {
	return &authpb.Confirmation{Status: "ok"}, nil
}

func (m *mockNotif) ClearAll(ctx context.Context, in *notificationpb.EmptyRequest, opts ...grpc.CallOption) (*authpb.Confirmation, error) {
	return &authpb.Confirmation{Status: "ok"}, nil
}

func (m *mockNotif) StreamNotifications(ctx context.Context, in *notificationpb.StreamRequest, opts ...grpc.CallOption) (notificationpb.NotificationService_StreamNotificationsClient, error) {
	return nil, nil // стрим нам не нужен
}

func (m *mockNotif) CreateNotification(ctx context.Context, req *notificationpb.CreateNotificationRequest, opts ...grpc.CallOption) (*authpb.Confirmation, error) {
	return &authpb.Confirmation{Status: "ok"}, nil
}

// =====================================================
//
//	TestMain (как в AuthService)
//
// =====================================================
func TestMain(m *testing.M) {

	dsn := "postgres://ostafon:0000@localhost:5432/test?sslmode=disable"

	fmt.Println("🔵 Running USER migrations on TEST DB...")

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Чистим таблицы
	_ = testDB.Exec(`DROP TABLE IF EXISTS follows CASCADE`)
	_ = testDB.Exec(`DROP TABLE IF EXISTS users CASCADE`)

	// Миграции
	if err := testDB.AutoMigrate(&model.User{}, &model.Follow{}); err != nil {
		panic(err)
	}

	// Репозиторий
	testRepo = repos.NewUserRepo(testDB)

	// 🔥 СОЗДАЁМ МОК ДО СОЗДАНИЯ СЕРВИСА!
	mockNotifClient := &mockNotif{}

	// 🔥 СОЗДАЁМ GRPCClients ВРУЧНУЮ
	mockClients := &config.GRPCClients{
		NotifClient: mockNotifClient,
	}

	// 🔥 ТЕПЕРЬ СОЗДАЁМ СЕРВИС С mockClients
	testSvc = service.NewUserService(testRepo, mockClients)

	// чтобы FollowUser не паниковал на контексте
	ctx = context.Background()

	// Запускаем тесты
	code := m.Run()
	os.Exit(code)
}

// =======================================================================
//                           TESTS
// =======================================================================

// ---------------------------------------------------------
// Create test user helper
// ---------------------------------------------------------
func createUser(first, last string) (*model.User, error) {
	u := &model.User{
		Firstname: first,
		Lastname:  last,
	}
	if err := testDB.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// ---------------------------------------------------------
// GET USER SUCCESS
// ---------------------------------------------------------
func TestGetUser_Success(t *testing.T) {
	u, err := createUser("Test", "User")
	assert.NoError(t, err)

	resp, err := testSvc.GetUser(&pb.GetUserRequest{Id: fmt.Sprint(u.Id)})

	assert.NoError(t, err)
	assert.Equal(t, "Test", resp.Firstname)
}

// ---------------------------------------------------------
// GET USER NOT FOUND
// ---------------------------------------------------------
func TestGetUser_NotFound(t *testing.T) {
	_, err := testSvc.GetUser(&pb.GetUserRequest{Id: "999999"})
	assert.Error(t, err)
}

// ---------------------------------------------------------
// UpdateUser SUCCESS
// ---------------------------------------------------------
func TestUpdateUser_Success(t *testing.T) {
	u, _ := createUser("Old", "Name")

	err := testSvc.UpdateUser(&pb.UpdateUserRequest{
		Id:        fmt.Sprint(u.Id),
		FirstName: "New",
		LastName:  "Name",
		Bio:       "Updated bio",
		BirthDate: "2020-01-01",
	})
	assert.NoError(t, err)

	updated, _ := testRepo.GetUser(u.Id)
	assert.Equal(t, "New", updated.Firstname)
	assert.Equal(t, "Updated bio", updated.Bio)
}

// ---------------------------------------------------------
// FOLLOW SUCCESS
// ---------------------------------------------------------
func TestFollowUser_Success(t *testing.T) {
	u1, _ := createUser("A", "One")
	u2, _ := createUser("B", "Two")

	err := testSvc.FollowUser(ctx, fmt.Sprint(u1.Id), fmt.Sprint(u2.Id))
	assert.NoError(t, err)

	var count int64
	testDB.Model(&model.Follow{}).Where("follower_id=? AND following_id=?", u1.Id, u2.Id).Count(&count)
	assert.Equal(t, int64(1), count)
}

// ---------------------------------------------------------
// FOLLOW YOURSELF ERROR
// ---------------------------------------------------------
func TestFollowUser_CannotFollowYourself(t *testing.T) {
	u, _ := createUser("Self", "Test")

	err := testSvc.FollowUser(nil, fmt.Sprint(u.Id), fmt.Sprint(u.Id))
	assert.Error(t, err)
}

// ---------------------------------------------------------
// UNFOLLOW SUCCESS
// ---------------------------------------------------------
func TestUnfollowUser_Success(t *testing.T) {
	u1, _ := createUser("A", "One")
	u2, _ := createUser("B", "Two")

	// follow
	_ = testSvc.FollowUser(ctx, fmt.Sprint(u1.Id), fmt.Sprint(u2.Id))

	// unfollow
	err := testSvc.UnfollowUser(fmt.Sprint(u1.Id), fmt.Sprint(u2.Id))
	assert.NoError(t, err)

	var count int64
	testDB.Model(&model.Follow{}).Where("follower_id=? AND following_id=?", u1.Id, u2.Id).Count(&count)
	assert.Equal(t, int64(0), count)
}

// ---------------------------------------------------------
// GET FOLLOWING SUCCESS
// ---------------------------------------------------------
func TestGetFollowing_Success(t *testing.T) {
	follower, _ := createUser("Follower", "X")
	following, _ := createUser("Following", "Y")

	_ = testSvc.FollowUser(ctx, fmt.Sprint(follower.Id), fmt.Sprint(following.Id))

	users, err := testSvc.GetFollowing(fmt.Sprint(follower.Id))
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, fmt.Sprint(following.Id), users[0].Id)
}

// ---------------------------------------------------------
// GET ALL USERS
// ---------------------------------------------------------
func TestGetAllUsers_Success(t *testing.T) {
	_, _ = createUser("U1", "L1")
	_, _ = createUser("U2", "L2")

	res, err := testSvc.GetAllUsers()
	assert.NoError(t, err)
	assert.True(t, len(res.Users) >= 2)
}
