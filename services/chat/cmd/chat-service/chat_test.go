package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"socialnet/services/chat/internal/service"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"socialnet/pkg/config"
	"socialnet/pkg/contextx"
	authpb "socialnet/services/auth/gen"
	notificationpb "socialnet/services/notification/gen"

	pb "socialnet/services/chat/gen"
	"socialnet/services/chat/internal/model"
	"socialnet/services/chat/internal/repos"
)

// =========================================================
// GLOBALS
// =========================================================

var (
	testDB    *gorm.DB
	testRepo  *repos.ChatRepo
	testSvc   *service.ChatService
	testRedis *redis.Client
	testCtx   = context.WithValue(context.Background(), contextx.UserIDKey, "1")
)

// ---------------------------------------------------------
// MOCK Notification Client
// ---------------------------------------------------------

type mockNotif struct {
	Created []*notificationpb.CreateNotificationRequest
}

func (m *mockNotif) CreateNotification(ctx context.Context, in *notificationpb.CreateNotificationRequest, opts ...grpc.CallOption) (*authpb.Confirmation, error) {
	m.Created = append(m.Created, in)
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
	var err error
	dsn := "postgres://ostafon:0000@localhost:5432/test?sslmode=disable"

	fmt.Println("🔵 Running CHAT migrations on TEST DB...")

	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	_ = testDB.Exec(`DROP TABLE IF EXISTS messages CASCADE`)
	_ = testDB.Exec(`DROP TABLE IF EXISTS chats CASCADE`)
	_ = testDB.Exec(`DROP TABLE IF EXISTS participants CASCADE`)

	if err := testDB.AutoMigrate(&model.Chat{}, &model.Message{}, &model.Participant{}); err != nil {
		panic(err)
	}

	testRepo = repos.NewChatRepo(testDB)

	testRedis = redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	mockNotifClient := &mockNotif{}

	clients := &config.GRPCClients{
		NotifClient: mockNotifClient,
	}

	testSvc = service.NewChatService(testRepo, testRedis, clients)

	m.Run()
}

// =========================================================
// TESTS
// =========================================================

func TestCreateChat(t *testing.T) {
	req := &pb.CreateChatRequest{
		Name:         "Test Chat",
		Participants: []string{"2", "3"},
	}

	resp, err := testSvc.CreateChat(testCtx, req)
	assert.NoError(t, err)
	assert.Equal(t, "Test Chat", resp.Name)
	assert.Len(t, resp.Participants, 3)
}

func TestSendMessage(t *testing.T) {
	// создаём чат
	chat := &model.Chat{Name: "X", CreatedAt: time.Now()}
	_ = testDB.Create(chat)

	req := &pb.SendMessageRequest{
		ChatId:      fmt.Sprint(chat.ID),
		Content:     "Hello!",
		ContentType: "text",
	}

	resp, err := testSvc.SendMessage(testCtx, req)
	assert.NoError(t, err)
	assert.Equal(t, "Hello!", resp.Content)

	var cnt int64
	testDB.Model(&model.Message{}).Count(&cnt)
	assert.Equal(t, int64(1), cnt)
}

func TestListMessages(t *testing.T) {
	chat := &model.Chat{Name: "LM", CreatedAt: time.Now()}
	_ = testDB.Create(chat)

	_ = testDB.Create(&model.Message{
		ChatID:    chat.ID,
		SenderID:  "1",
		Content:   "msg1",
		CreatedAt: time.Now(),
	})

	resp, err := testSvc.ListMessages(testCtx, &pb.ListMessagesRequest{
		ChatId: fmt.Sprint(chat.ID),
		Limit:  10,
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Messages, 1)
}

func TestListChats(t *testing.T) {
	chat := &model.Chat{Name: "LC", CreatedAt: time.Now()}
	_ = testDB.Create(chat)
	_ = testDB.Create(&model.Participant{
		ChatID: chat.ID,
		UserID: "1",
	})

	resp, err := testSvc.ListChats(testCtx, &pb.EmptyRequest{})
	assert.NoError(t, err)
	assert.True(t, len(resp.Chats) >= 1)
}

// ---- fake stream ----

type fakeStream struct {
	ctx     context.Context
	recvCh  chan *pb.Message
	sendLog []*pb.Message
}

func newFakeStream() *fakeStream {
	return &fakeStream{
		ctx:    context.Background(),
		recvCh: make(chan *pb.Message, 10),
	}
}

func (f *fakeStream) Context() context.Context {
	return f.ctx
}

func (f *fakeStream) Send(msg *pb.Message) error {
	f.sendLog = append(f.sendLog, msg)
	return nil
}

// интерфейс требует — оставляем пустыми
func (f *fakeStream) SendHeader(md metadata.MD) error { return nil }
func (f *fakeStream) SetHeader(md metadata.MD) error  { return nil }
func (f *fakeStream) SetTrailer(md metadata.MD)       {}
func (f *fakeStream) SendMsg(m interface{}) error     { return nil }
func (f *fakeStream) RecvMsg(m interface{}) error     { return nil }
