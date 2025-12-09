package handlers

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"socialnet/pkg/contextx"
	pbauth "socialnet/services/auth/gen"
	pb "socialnet/services/comment/gen"
	"socialnet/services/comment/internal/service"
)

type CommentHandler struct {
	pb.UnimplementedCommentServiceServer
	service *service.CommentService
}

func NewCommentHandler(s *service.CommentService) *CommentHandler {
	return &CommentHandler{service: s}
}

func (h *CommentHandler) AddComment(ctx context.Context, req *pb.AddCommentRequest) (*pb.Comment, error) {
	userID := contextx.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}
	return h.service.AddComment(ctx, userID, req)
}

func (h *CommentHandler) GetComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.Comment, error) {
	return h.service.GetComment(ctx, req.Id)
}

func (h *CommentHandler) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pbauth.Confirmation, error) {
	userID := contextx.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}
	if err := h.service.DeleteComment(ctx, req.Id, userID); err != nil {
		return nil, err
	}
	return &pbauth.Confirmation{Status: "deleted"}, nil
}

func (h *CommentHandler) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.Comments, error) {
	return h.service.ListComments(ctx, req.PostId)
}
