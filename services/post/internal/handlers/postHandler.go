package handlers

import (
	pb "socialnet/services/post/gen"
	"socialnet/services/post/internal/service"
)

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	s *service.PostService
}

func NewPostHandler(s *service.PostService) *PostHandler {
	return &PostHandler{s: s}
}
