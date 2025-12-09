package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"socialnet/pkg/contextx"
	authpb "socialnet/services/auth/gen"
	pb "socialnet/services/notification/gen"
	"socialnet/services/notification/internal/model"
	"socialnet/services/notification/internal/repos"
	"time"
)

type NotificationService struct {
	repo  *repos.NotificationRepo
	redis *redis.Client
}

func NewNotificationService(repo *repos.NotificationRepo, redis *redis.Client) *NotificationService {
	return &NotificationService{repo: repo, redis: redis}
}

// 📩 Получить список уведомлений
func (s *NotificationService) ListNotifications(ctx context.Context, req *pb.ListNotificationsRequest) (*pb.Notifications, error) {
	userID := contextx.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}

	notes, err := s.repo.List(userID, req.Filter, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get notifications: %v", err)
	}

	var pbNotes []*pb.Notification
	for _, n := range notes {
		pbNotes = append(pbNotes, &pb.Notification{
			Id:          fmt.Sprint(n.ID),
			UserId:      n.UserID,
			Type:        n.Type,
			ReferenceId: n.ReferenceID,
			Content:     n.Content,
			Read:        n.Read,
			CreatedAt:   n.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.Notifications{Notifications: pbNotes}, nil
}

// ✅ Отметить уведомление как прочитанное
func (s *NotificationService) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*authpb.Confirmation, error) {
	if err := s.repo.MarkAsRead(req.Id); err != nil {
		return nil, status.Error(codes.Internal, "failed to mark as read")
	}
	return &authpb.Confirmation{Status: "Marked as read"}, nil
}

// ✅ Отметить все как прочитанные
func (s *NotificationService) MarkAllAsRead(ctx context.Context, _ *pb.EmptyRequest) (*authpb.Confirmation, error) {
	userID := contextx.GetUserID(ctx)
	if err := s.repo.MarkAllAsRead(userID); err != nil {
		return nil, status.Error(codes.Internal, "failed to mark all as read")
	}
	return &authpb.Confirmation{Status: "All marked as read"}, nil
}

// 🗑️ Удалить одно уведомление
func (s *NotificationService) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*authpb.Confirmation, error) {
	if err := s.repo.Delete(req.Id); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete notification")
	}
	return &authpb.Confirmation{Status: "Deleted"}, nil
}

// 🗑️ Очистить все уведомления
func (s *NotificationService) ClearAll(ctx context.Context, _ *pb.EmptyRequest) (*authpb.Confirmation, error) {
	userID := contextx.GetUserID(ctx)
	if err := s.repo.ClearAll(userID); err != nil {
		return nil, status.Error(codes.Internal, "failed to clear notifications")
	}
	return &authpb.Confirmation{Status: "Cleared all"}, nil
}

// 🔔 Поток уведомлений в реальном времени (через Redis Pub/Sub)
func (s *NotificationService) StreamNotifications(req *pb.StreamRequest, stream pb.NotificationService_StreamNotificationsServer) error {
	ctx := stream.Context()
	pubsub := s.redis.Subscribe(ctx, fmt.Sprintf("notifications:%s", req.UserId))
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-ch:
			var notif pb.Notification
			if err := json.Unmarshal([]byte(msg.Payload), &notif); err == nil {
				_ = stream.Send(&notif)
			}
		}
	}
}

// 📢 Вспомогательная функция для публикации нового уведомления
func (s *NotificationService) PublishNotification(ctx context.Context, n *model.Notification) {
	data, _ := json.Marshal(&pb.Notification{
		Id:          fmt.Sprint(n.ID),
		UserId:      n.UserID,
		Type:        n.Type,
		ReferenceId: n.ReferenceID,
		Content:     n.Content,
		Read:        n.Read,
		CreatedAt:   n.CreatedAt.Format(time.RFC3339),
	})
	s.redis.Publish(ctx, fmt.Sprintf("notifications:%s", n.UserID), data)
}
