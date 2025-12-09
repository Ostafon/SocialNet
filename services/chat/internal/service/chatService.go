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
	pb "socialnet/services/chat/gen"
	"socialnet/services/chat/internal/model"
	"socialnet/services/chat/internal/repos"
	"time"
)

type ChatService struct {
	repo  *repos.ChatRepo
	redis *redis.Client
}

func NewChatService(repo *repos.ChatRepo, redis *redis.Client) *ChatService {
	return &ChatService{repo: repo, redis: redis}
}

func (s *ChatService) CreateChat(ctx context.Context, req *pb.CreateChatRequest) (*pb.Chat, error) {
	userID := contextx.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}

	if len(req.Participants) == 0 {
		return nil, status.Error(codes.InvalidArgument, "participants required")
	}

	chat := &model.Chat{
		Name:      req.Name,
		IsGroup:   len(req.Participants) > 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateChat(chat); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create chat: %v", err)
	}

	all := append(req.Participants, userID)
	for _, id := range all {
		_ = s.repo.AddParticipant(chat.ID, id)
	}

	return &pb.Chat{
		Id:           fmt.Sprint(chat.ID),
		Name:         chat.Name,
		IsGroup:      chat.IsGroup,
		Participants: all,
		CreatedAt:    chat.CreatedAt.Format(time.RFC3339),
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

	data, _ := json.Marshal(event)
	_ = s.redis.Publish(ctx, fmt.Sprintf("chat:%s", req.ChatId), data).Err()

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

func (s *ChatService) SubscribeMessages(req *pb.SubscribeRequest, stream pb.ChatService_SubscribeMessagesServer) error {
	ctx := stream.Context()
	var channels []string
	for _, id := range req.ChatIds {
		channels = append(channels, fmt.Sprintf("chat:%s", id))
		log.Printf("✅ User subscribed to chat channels: %v\n", req.ChatIds)

	}
	pubsub := s.redis.Subscribe(ctx, channels...)

	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-ch:
			var m pb.Message
			if err := json.Unmarshal([]byte(msg.Payload), &m); err == nil {
				if err := stream.Send(&m); err != nil {
					return err
				}
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
		// преобразуем []Participant → []string
		var participants []string
		for _, p := range c.Participants {
			participants = append(participants, p.UserID)
		}
		lastMsg, _ := s.repo.GetLastMessage(c.ID)
		pbChats = append(pbChats, &pb.Chat{
			Id:           fmt.Sprint(c.ID),
			Participants: participants,    // ✅ теперь []string
			LastMessage:  lastMsg.Content, // см. ниже
			CreatedAt:    c.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    c.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.Chats{Chats: pbChats}, nil
}
