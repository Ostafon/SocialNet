package main

import (
	"context"
	"fmt"
	"os"
	"socialnet/pkg/contextx"
	"socialnet/services/notification/internal/service"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pb "socialnet/services/notification/gen"
	"socialnet/services/notification/internal/model"
	"socialnet/services/notification/internal/repos"
)

var (
	testDB   *gorm.DB
	testRepo *repos.NotificationRepo
	testSvc  *service.NotificationService
	rdb      *redis.Client
)

//
// ------------------------- HELPERS -------------------------
//

var ctx = context.WithValue(context.Background(), contextx.UserIDKey, "1")

//
// ------------------------- FAKE STREAM -------------------------
//

type fakeStream struct {
	pb.NotificationService_StreamNotificationsServer
	ctx   context.Context
	recvC chan *pb.Notification
}

func newFakeStream() *fakeStream {
	return &fakeStream{
		ctx:   context.Background(),
		recvC: make(chan *pb.Notification, 10),
	}
}

func (s *fakeStream) Context() context.Context {
	return s.ctx
}

func (s *fakeStream) Send(n *pb.Notification) error {
	s.recvC <- n
	return nil
}

//
// ------------------------- TEST MAIN -------------------------
//

func TestMain(m *testing.M) {
	var err error

	// PostgreSQL
	dsn := "postgres://ostafon:0000@localhost:5432/test?sslmode=disable"
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	_ = testDB.Exec(`DROP TABLE IF EXISTS notifications CASCADE`)
	if err := testDB.AutoMigrate(&model.Notification{}); err != nil {
		panic(err)
	}

	testRepo = repos.NewNotificationRepo(testDB)

	// Redis
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	testSvc = service.NewNotificationService(testRepo, rdb)

	os.Exit(m.Run())
}

//
// ------------------------- TESTS -------------------------
//

func TestCreateNotification(t *testing.T) {

	_, err := testSvc.CreateNotification(ctx, &pb.CreateNotificationRequest{
		UserId:      "1",
		Type:        "like",
		ReferenceId: "post1",
		Content:     "Hello",
	})
	assert.NoError(t, err)

	var count int64
	testDB.Model(&model.Notification{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestListNotifications(t *testing.T) {
	_ = testDB.Create(&model.Notification{
		UserID:    "1",
		Type:      "msg",
		Content:   "test",
		CreatedAt: time.Now(),
	})
	ctx = context.WithValue(context.Background(), contextx.UserIDKey, "1")

	resp, err := testSvc.ListNotifications(ctx, &pb.ListNotificationsRequest{
		Limit:  1,
		Offset: 1,
		Filter: "unread",
	})
	assert.NoError(t, err)
	assert.True(t, len(resp.Notifications) >= 1)
}

func TestMarkAsRead(t *testing.T) {
	n := &model.Notification{UserID: "1", Type: "mark"}
	_ = testDB.Create(n)

	_, err := testSvc.MarkAsRead(ctx, &pb.MarkAsReadRequest{Id: fmt.Sprint(n.ID)})
	assert.NoError(t, err)

	var nn model.Notification
	testDB.First(&nn, n.ID)
	assert.True(t, nn.Read)
}

func TestMarkAllAsRead(t *testing.T) {
	testDB.Create(&model.Notification{UserID: "55"})
	testDB.Create(&model.Notification{UserID: "55"})

	ctx = context.WithValue(context.Background(), contextx.UserIDKey, "55")

	_, err := testSvc.MarkAllAsRead(ctx, &pb.EmptyRequest{})
	assert.NoError(t, err)

	var count int64
	testDB.Model(&model.Notification{}).
		Where("user_id=? AND read=true", "55").
		Count(&count)

	assert.Equal(t, int64(2), count)
}

func TestDeleteNotification(t *testing.T) {
	n := &model.Notification{UserID: "1"}
	testDB.Create(n)

	ctx = context.WithValue(context.Background(), contextx.UserIDKey, "1")

	_, err := testSvc.DeleteNotification(ctx, &pb.DeleteNotificationRequest{Id: fmt.Sprint(n.ID)})
	assert.NoError(t, err)

	var count int64
	testDB.Model(&model.Notification{}).Where("id=?", n.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestClearAll(t *testing.T) {
	testDB.Create(&model.Notification{UserID: "777"})
	testDB.Create(&model.Notification{UserID: "777"})

	ctx = context.WithValue(context.Background(), contextx.UserIDKey, "777")

	_, err := testSvc.ClearAll(ctx, &pb.EmptyRequest{})
	assert.NoError(t, err)

	var count int64
	testDB.Model(&model.Notification{}).Where("user_id=?", "777").Count(&count)
	assert.Equal(t, int64(0), count)
}
