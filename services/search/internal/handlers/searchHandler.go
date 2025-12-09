package handlers

import (
	"context"

	pb "socialnet/services/search/gen"
	"socialnet/services/search/internal/service"
)

type SearchHandler struct {
	svc *service.SearchService
	pb.UnimplementedSearchServiceServer
}

func NewSearchHandler(svc *service.SearchService) *SearchHandler {
	return &SearchHandler{svc: svc}
}

func (h *SearchHandler) SearchUsers(ctx context.Context, req *pb.SearchRequest) (*pb.SearchUsersResponse, error) {
	return h.svc.SearchUsers(ctx, req)
}

func (h *SearchHandler) SearchPosts(ctx context.Context, req *pb.SearchRequest) (*pb.SearchPostsResponse, error) {
	return h.svc.SearchPosts(ctx, req)
}

func (h *SearchHandler) SearchAll(ctx context.Context, req *pb.SearchRequest) (*pb.SearchAllResponse, error) {
	return h.svc.SearchAll(ctx, req)
}
