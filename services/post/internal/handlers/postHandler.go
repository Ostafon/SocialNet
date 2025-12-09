package handlers

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	authpb "socialnet/services/auth/gen"
	pb "socialnet/services/post/gen"
	"socialnet/services/post/internal/service"
	userpb "socialnet/services/user/gen"
)

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	service *service.PostService
}

func NewPostHandler(s *service.PostService) *PostHandler {
	return &PostHandler{service: s}
}

// CreatePost — создать пост (контент + изображение)
func (h *PostHandler) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	post, err := h.service.CreatePost(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create post: %v", err)
	}
	return post, nil
}

// GetPost — получить один пост по ID
func (h *PostHandler) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	post, err := h.service.GetPost(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post: %v", err)
	}
	return post, nil
}

// DeletePost — удалить пост
func (h *PostHandler) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*authpb.Confirmation, error) {
	err := h.service.DeletePost(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete post: %v", err)
	}
	return &authpb.Confirmation{Status: "Post deleted successfully"}, nil
}

// ListUserPosts — все посты конкретного пользователя
func (h *PostHandler) ListUserPosts(ctx context.Context, req *pb.UserPostsRequest) (*pb.Posts, error) {
	posts, err := h.service.ListUserPosts(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list user posts: %v", err)
	}
	return posts, nil
}

// GetFeed — лента пользователя (его посты + посты друзей)
func (h *PostHandler) GetFeed(ctx context.Context, req *pb.GetFeedRequest) (*pb.Posts, error) {
	feed, err := h.service.GetFeed(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load feed: %v", err)
	}
	fmt.Printf("Feed loaded successfully: %d posts\n", len(feed.Posts))
	return feed, nil
}

func (h *PostHandler) ListPosts(ctx context.Context, req *userpb.EmptyRequest) (*pb.Posts, error) {
	posts, err := h.service.GetAllPosts()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load feed: %v", err)
	}
	fmt.Printf("Feed loaded successfully: %d posts\n", len(posts.Posts))
	return posts, nil
}
