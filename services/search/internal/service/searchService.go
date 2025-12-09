package service

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	postpb "socialnet/services/post/gen"
	pb "socialnet/services/search/gen"
	userpb "socialnet/services/user/gen"
)

type SearchService struct {
	userClient userpb.UserServiceClient
	postClient postpb.PostServiceClient
}

func NewSearchService(userClient userpb.UserServiceClient, postClient postpb.PostServiceClient) *SearchService {
	return &SearchService{
		userClient: userClient,
		postClient: postClient,
	}
}

// SearchUsers 🔍 Поиск пользователей
func (s *SearchService) SearchUsers(ctx context.Context, req *pb.SearchRequest) (*pb.SearchUsersResponse, error) {
	if strings.TrimSpace(req.Query) == "" {
		return nil, status.Error(codes.InvalidArgument, "query cannot be empty")
	}

	usersResp, err := s.userClient.ListUsers(ctx, &userpb.EmptyRequest{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch users: %v", err)
	}

	query := strings.ToLower(req.Query)
	var results []*userpb.User

	for _, u := range usersResp.Users {
		if strings.Contains(strings.ToLower(u.FirstName), query) ||
			strings.Contains(strings.ToLower(u.LastName), query) {
			results = append(results, u)
		}
		if req.Limit > 0 && int32(len(results)) >= req.Limit {
			break
		}
	}

	return &pb.SearchUsersResponse{Users: results}, nil
}

// SearchPosts 🔍 Поиск постов
func (s *SearchService) SearchPosts(ctx context.Context, req *pb.SearchRequest) (*pb.SearchPostsResponse, error) {
	if strings.TrimSpace(req.Query) == "" {
		return nil, status.Error(codes.InvalidArgument, "query cannot be empty")
	}

	postsResp, err := s.postClient.ListPosts(ctx, &userpb.EmptyRequest{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch posts: %v", err)
	}

	query := strings.ToLower(req.Query)
	var results []*postpb.Post

	for _, p := range postsResp.Posts {
		if strings.Contains(strings.ToLower(p.Content), query) {
			results = append(results, p)
		}
		if req.Limit > 0 && int32(len(results)) >= req.Limit {
			break
		}
	}

	return &pb.SearchPostsResponse{Posts: results}, nil
}

func (s *SearchService) SearchAll(ctx context.Context, req *pb.SearchRequest) (*pb.SearchAllResponse, error) {
	if s.userClient == nil || s.postClient == nil {
		return nil, status.Error(codes.Internal, "missing user or post gRPC client")
	}

	users, errU := s.SearchUsers(ctx, req)
	posts, errP := s.SearchPosts(ctx, req)

	if errU != nil && errP != nil {
		return nil, status.Error(codes.Internal, "both user and post searches failed")
	}

	return &pb.SearchAllResponse{
		Users: users.GetUsers(),
		Posts: posts.GetPosts(),
	}, nil
}
