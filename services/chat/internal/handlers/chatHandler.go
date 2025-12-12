package handlers

import (
	"context"
	authpb "socialnet/services/auth/gen"
	pb "socialnet/services/chat/gen"
	"socialnet/services/chat/internal/service"
)

type ChatHandler struct {
	pb.UnimplementedChatServiceServer
	s *service.ChatService
}

func NewChatHandler(s *service.ChatService) *ChatHandler {
	return &ChatHandler{s: s}
}

func (h *ChatHandler) CreateChat(ctx context.Context, req *pb.CreateChatRequest) (*pb.Chat, error) {
	return h.s.CreateChat(ctx, req)
}

func (h *ChatHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.Message, error) {
	return h.s.SendMessage(ctx, req)
}

func (h *ChatHandler) ListMessages(ctx context.Context, req *pb.ListMessagesRequest) (*pb.Messages, error) {
	return h.s.ListMessages(ctx, req)
}

func (h *ChatHandler) SubscribeMessages(req *pb.SubscribeRequest, stream pb.ChatService_SubscribeMessagesServer) error {
	return h.s.SubscribeMessages(req, stream)
}

func (h *ChatHandler) ListChats(ctx context.Context, req *pb.EmptyRequest) (*pb.Chats, error) {
	return h.s.ListChats(ctx, req)
}

func (h *ChatHandler) GetChat(ctx context.Context, req *pb.GetChatRequest) (*pb.Chat, error) {
	return h.s.GetChat(ctx, req)
}

func (h *ChatHandler) DeleteChat(ctx context.Context, req *pb.DeleteChatRequest) (*authpb.Confirmation, error) {
	return h.s.DeleteChat(ctx, req)
}

func (h *ChatHandler) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*authpb.Confirmation, error) {
	return h.s.MarkAsRead(ctx, req)
}
