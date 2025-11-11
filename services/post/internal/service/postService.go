package service

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"socialnet/pkg/storage"
	"socialnet/pkg/utils"
	pb "socialnet/services/post/gen"
	"socialnet/services/post/internal/model"
	"socialnet/services/post/internal/repos"
)

type PostService struct {
	repo *repos.PostRepo
}

func NewPostService(repo *repos.PostRepo) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	if req.Content == "" && req.FileName == "" {
		return nil, status.Error(codes.InvalidArgument, "content and image cannot be null")
	}
	user_id := ctx.Value("user_id").(string)
	post := &model.Post{}

	reader := bytes.NewReader(req.Image)
	s3, err := storage.NewS3Client()
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot to create client")
	}

	key := fmt.Sprintf("avatars/%s_%s", user_id, req.FileName)

	url, err := s3.UploadFile(reader, key)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload avatar: %v", err)
	}

	post.UserId = user_id
	post.ImageUrl = url
	post.Content = req.Content
	post.LikesCount = 0
	post.CommentsCount = 0

	err = s.repo.SavePost(post)
	if err != nil {
		return nil, err
	}

	return &pb.Post{}

}

func (s *PostService) GetPost(req *pb.GetPostRequest) (*pb.Post, error)          {}
func (s *PostService) DeletePost(req *pb.DeletePostRequest) error                {}
func (s *PostService) ListUserPosts(req *pb.UserPostsRequest) (*pb.Posts, error) {}
