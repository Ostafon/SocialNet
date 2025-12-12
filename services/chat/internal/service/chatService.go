package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"socialnet/pkg/config"
	"socialnet/pkg/contextx"
	pb "socialnet/services/chat/gen"
	"socialnet/services/chat/internal/model"
	"socialnet/services/chat/internal/repos"
	notificationpb "socialnet/services/notification/gen"
	"time"
)

type ChatService struct {
	repo    *repos.ChatRepo
	redis   *redis.Client
	clients *config.GRPCClients
}

func NewChatService(repo *repos.ChatRepo, redis *redis.Client, clients *config.GRPCClients) *ChatService {
	return &ChatService{repo: repo, redis: redis, clients: clients}
}

func (s *ChatService) CreateChat(ctx context.Context, req *pb.CreateChatRequest) (*pb.Chat, error) {
	userID := contextx.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}

	if len(req.Participants) == 0 {
		return nil, status.Error(codes.InvalidArgument, "participants required")
	}

	// =============== PRIVATE CHAT (1-on-1) ===============
	if len(req.Participants) == 1 {
		other := req.Participants[0]

		existing, err := s.repo.FindPrivateChat(userID, other)
		if err == nil && existing != nil && existing.ID != 0 {
			// Приватный чат уже существует
			return &pb.Chat{
				Id:           fmt.Sprint(existing.ID),
				Name:         existing.Name,
				IsGroup:      false,
				Participants: []string{userID, other},
				CreatedAt:    existing.CreatedAt.Format(time.RFC3339),
				UpdatedAt:    existing.UpdatedAt.Format(time.RFC3339),
			}, nil
		}
	}

	// =============== CREATE NEW CHAT ===============
	chat := &model.Chat{
		Name:      req.Name,
		IsGroup:   len(req.Participants) > 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateChat(chat); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create chat: %v", err)
	}

	// Участники: все + автор
	allParticipants := make([]string, 0, len(req.Participants)+1)
	allParticipants = append(allParticipants, req.Participants...)

	// Добавляем userID, но ТОЛЬКО если его нет
	exists := false
	for _, p := range allParticipants {
		if p == userID {
			exists = true
			break
		}
	}
	if !exists {
		allParticipants = append(allParticipants, userID)
	}

	// =============== SAVE PARTICIPANTS ===============
	for _, id := range allParticipants {
		if err := s.repo.AddParticipant(chat.ID, id); err != nil {
			log.Printf("⚠ failed AddParticipant for user %s: %v", id, err)
		}
	}

	// =============== RESPONSE ===============
	return &pb.Chat{
		Id:           fmt.Sprint(chat.ID),
		Name:         chat.Name,
		IsGroup:      chat.IsGroup,
		Participants: allParticipants,
		CreatedAt:    chat.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    chat.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *ChatService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.Message, error) {
	userID := contextx.GetUserID(ctx)

	msg := &model.Message{
		ChatID:      parseUint(req.ChatId),
		SenderID:    userID,
		Content:     req.Content,
		ContentType: req.ContentType,
		MediaURL:    req.MediaUrl,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.SaveMessage(msg); err != nil {
		return nil, status.Error(codes.Internal, "failed to save message")
	}

	event := pb.Message{
		Id:          fmt.Sprint(msg.ID),
		ChatId:      req.ChatId,
		SenderId:    userID,
		Content:     msg.Content,
		ContentType: msg.ContentType,
		MediaUrl:    msg.MediaURL,
		CreatedAt:   msg.CreatedAt.Format(time.RFC3339),
	}

	// Redis pubsub
	data, _ := json.Marshal(event)
	_ = s.redis.Publish(ctx, fmt.Sprintf("chat:%s", req.ChatId), data).Err()

	// ==== УВЕДОМЛЕНИЕ ====

	notifClient, err := s.clients.GetNotifClient("localhost:50057")
	if err == nil {

		// создаём metadata
		md := metadata.New(map[string]string{
			"user-id": userID,
		})
		ctxWithUser := metadata.NewOutgoingContext(ctx, md)

		participants, _ := s.repo.GetChatParticipants(parseUint(req.ChatId))
		for _, p := range participants {
			if p.UserID == userID {
				continue
			}

			_, _ = notifClient.CreateNotification(ctxWithUser,
				&notificationpb.CreateNotificationRequest{
					UserId:      p.UserID,
					Type:        "new_message",
					ReferenceId: fmt.Sprint(msg.ID),
					Content:     fmt.Sprintf("New message from %s", userID),
				})
		}
	}

	return &event, nil
}

func (s *ChatService) ListMessages(ctx context.Context, req *pb.ListMessagesRequest) (*pb.Messages, error) {
	msgs, err := s.repo.GetMessages(parseUint(req.ChatId), int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load messages: %v", err)
	}

	var pbMsgs []*pb.Message
	for _, m := range msgs {
		pbMsgs = append(pbMsgs, &pb.Message{
			Id:          fmt.Sprint(m.ID),
			ChatId:      fmt.Sprint(m.ChatID),
			SenderId:    m.SenderID,
			Content:     m.Content,
			ContentType: m.ContentType,
			MediaUrl:    m.MediaURL,
			Read:        m.Read,
			CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.Messages{Messages: pbMsgs}, nil
}

func (s *ChatService) SubscribeMessages(
	req *pb.SubscribeRequest,
	stream pb.ChatService_SubscribeMessagesServer,
) error {

	ctx := stream.Context()
	userID := contextx.GetUserID(ctx)

	log.Printf("🔵 [CHAT-STREAM-START] User %s subscribed to chats: %v", userID, req.ChatIds)

	// Формируем список каналов Redis
	var channels []string
	for _, id := range req.ChatIds {
		ch := fmt.Sprintf("chat:%s", id)
		log.Printf("🔔 [REDIS] Subscribing to channel: %s", ch)
		channels = append(channels, ch)
	}

	pubsub := s.redis.Subscribe(ctx, channels...)
	defer func() {
		_ = pubsub.Close()
		log.Printf("🟡 [CHAT-STREAM-END] User %s disconnected", userID)
	}()

	ch := pubsub.Channel()

	for {
		select {

		case <-ctx.Done():
			log.Printf("🟡 [CHAT-STREAM-END] Context closed for user %s", userID)
			return nil

		case msg, ok := <-ch:
			if !ok {
				log.Printf("🔴 [CHAT-STREAM-ERROR] Redis channel closed for user %s", userID)
				return nil
			}

			log.Printf("🔥 [REDIS → CHAT-SERVICE] Raw message: %s", msg.Payload)

			// Разбираем json → protobuf сообщения
			var m pb.Message
			if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
				log.Printf("❌ [ERROR] Failed to unmarshal redis message: %v", err)
				continue
			}

			log.Printf("📤 [CHAT-STREAM → CLIENT] Sending message %s from chat %s to user %s",
				m.Id, m.ChatId, userID)

			// Отправляем в stream
			if err := stream.Send(&m); err != nil {
				log.Printf("❌ [CHAT-STREAM-ERROR] Failed to send to client %s: %v", userID, err)
				return err
			}
		}
	}
}

func parseUint(s string) uint {
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}

func (s *ChatService) ListChats(ctx context.Context, req *pb.EmptyRequest) (*pb.Chats, error) {
	userID := contextx.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}

	chats, err := s.repo.GetChatsByUser(userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load chats: %v", err)
	}

	var pbChats []*pb.Chat

	for _, c := range chats {

		// собрать участников
		participants := make([]string, 0, len(c.Participants))
		for _, p := range c.Participants {
			participants = append(participants, p.UserID)
		}

		// получить последнее сообщение
		lastMsg, err := s.repo.GetLastMessage(c.ID)

		var lastText string
		if err == nil && lastMsg != nil {
			lastText = lastMsg.Content
		} else {
			lastText = "" // если сообщений нет
		}

		pbChats = append(pbChats, &pb.Chat{
			Id:           fmt.Sprint(c.ID),
			Name:         c.Name,
			IsGroup:      c.IsGroup,
			Participants: participants,
			LastMessage:  lastText,
			CreatedAt:    c.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    c.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.Chats{Chats: pbChats}, nil
}
