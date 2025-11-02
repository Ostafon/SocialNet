package handlers

import (
	"context"
	pb "socialnet/services/like/gen"
	"socialnet/services/like/internal/service"
)

type LikeHandler struct {
	pb.UnimplementedLikeServiceServer
	service *service.LikeService
}

func NewLikeHandler(s *service.LikeService) *LikeHandler {
	return &LikeHandler{service: s}
}

func (h *LikeHandler) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostResponse, error) {
	userID := ctx.Value("user_id").(string)
	return h.service.LikePost(ctx, userID, req.Id)
}

func (h *LikeHandler) UnlikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostResponse, error) {
	userID := ctx.Value("user_id").(string)
	return h.service.UnlikePost(ctx, userID, req.Id)
}

func (h *LikeHandler) LikeComment(ctx context.Context, req *pb.LikeCommentRequest) (*pb.LikeCommentResponse, error) {
	userID := ctx.Value("user_id").(string)
	return h.service.LikeComment(ctx, userID, req.Id)
}

func (h *LikeHandler) UnlikeComment(ctx context.Context, req *pb.LikeCommentRequest) (*pb.LikeCommentResponse, error) {
	userID := ctx.Value("user_id").(string)
	return h.service.UnlikeComment(ctx, userID, req.Id)
}

func (h *LikeHandler) ListPostLikes(ctx context.Context, req *pb.LikePostRequest) (*pb.ListLikesResponse, error) {
	return h.service.ListPostLikes(ctx, req.Id)
}

func (h *LikeHandler) ListCommentLikes(ctx context.Context, req *pb.LikeCommentRequest) (*pb.ListLikesResponse, error) {
	return h.service.ListCommentLikes(ctx, req.Id)
}
