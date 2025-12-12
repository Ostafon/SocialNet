package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
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

// ListNotifications  Получить список уведомлений
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

func (s *NotificationService) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*authpb.Confirmation, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id required")
	}

	n := &model.Notification{
		UserID:      req.UserId,
		Type:        req.Type,
		ReferenceID: req.ReferenceId,
		Content:     req.Content,
		Read:        false,
		CreatedAt:   time.Now(),
	}

	// сохраняем в БД
	if err := s.repo.Save(n); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save notification: %v", err)
	}

	// отправляем в реальном времени
	s.PublishNotification(ctx, n)

	return &authpb.Confirmation{Status: "Notification created"}, nil
}

// 🔔 Поток уведомлений в реальном времени (через Redis Pub/Sub)
func (s *NotificationService) StreamNotifications(
	req *pb.StreamRequest,
	stream pb.NotificationService_StreamNotificationsServer,
) error {

	ctx := stream.Context()
	userID := req.UserId

	log.Printf("🔵 [STREAM-START] User %s connected to notifications stream", userID)

	// Подписываемся на Redis канал
	channel := fmt.Sprintf("notifications:%s", userID)
	pubsub := s.redis.Subscribe(ctx, channel)
	defer pubsub.Close()

	log.Printf("🔵 [REDIS] Subscribed to channel: %s", channel)

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			log.Printf("🟡 [STREAM-END] User %s disconnected", userID)
			return nil

		case msg, ok := <-ch:
			if !ok {
				log.Printf("🔴 [ERROR] Redis channel closed for user %s", userID)
				return nil
			}

			log.Printf("🔥 [REDIS → SERVICE] Raw message for %s: %s", userID, msg.Payload)

			// 1️⃣ Парсим JSON в map
			var raw model.RedisNotification

			if err := json.Unmarshal([]byte(msg.Payload), &raw); err != nil {
				log.Printf("❌ Failed to unmarshal JSON: %v", err)
				continue
			}

			// 2️⃣ Создаём правильный protobuf объект вручную
			notif := &pb.Notification{
				Id:          raw.Id,
				UserId:      raw.UserID,
				Type:        raw.Type,
				ReferenceId: raw.ReferenceID,
				Content:     raw.Content,
				Read:        false,
				CreatedAt:   raw.CreatedAt,
			}

			log.Printf("📤 [STREAM → CLIENT] Sending notification to user %s: %+v", userID, notif)

			// 3️⃣ Отправляем как protobuf-байты → ЭТО ВАЖНО!
			if err := stream.Send(notif); err != nil {
				log.Printf("❌ [STREAM-ERROR] Failed send to %s: %v", userID, err)
				return err
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
