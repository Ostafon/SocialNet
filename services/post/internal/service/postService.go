package service

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"socialnet/pkg/contextx"
	"socialnet/pkg/storage"
	commentpb "socialnet/services/comment/gen"
	likepb "socialnet/services/like/gen"
	pb "socialnet/services/post/gen"
	"socialnet/services/post/internal/model"
	"socialnet/services/post/internal/repos"
	userpb "socialnet/services/user/gen"
)

type PostService struct {
	repo       *repos.PostRepo
	commClient commentpb.CommentServiceClient
	likeClient likepb.LikeServiceClient
	userClient userpb.UserServiceClient
}

func NewPostService(repo *repos.PostRepo, commClient commentpb.CommentServiceClient,
	likeClient likepb.LikeServiceClient, userClient userpb.UserServiceClient) *PostService {
	return &PostService{repo: repo, commClient: commClient, likeClient: likeClient, userClient: userClient}
}

func (s *PostService) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	if req.Content == "" && req.FileName == "" {
		return nil, status.Error(codes.InvalidArgument, "content and image cannot be null")
	}
	userID := contextx.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}
	post := &model.Post{
		UserId:        userID,
		Content:       req.Content,
		LikesCount:    0,
		CommentsCount: 0,
	}

	reader := bytes.NewReader(req.Image)
	s3, err := storage.NewS3Client()
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot to create client")
	}

	key := fmt.Sprintf("avatars/%s_%s", userID, req.FileName)

	url, err := s3.UploadFile(reader, key)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload avatar: %v", err)
	}

	post.ImageUrl = url

	err = s.repo.SavePost(post)
	if err != nil {
		return nil, err
	}

	return &pb.Post{
		Id:            fmt.Sprint(post.ID),
		UserId:        post.UserId,
		Content:       post.Content,
		ImageUrl:      post.ImageUrl,
		LikesCount:    post.LikesCount,
		CommentsCount: post.CommentsCount,
		CreatedAt:     post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *PostService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	post, err := s.repo.GetPostByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post not found: %v", err)
	}

	// 🔹 Получаем количество лайков
	likesResp, err := s.likeClient.LikePost(ctx, &likepb.LikePostRequest{
		Id: req.Id,
	})
	if err == nil {
		post.LikesCount = likesResp.LikesCount
	}

	// 🔹 Получаем количество комментариев
	commentsResp, err := s.commClient.ListComments(ctx, &commentpb.ListCommentsRequest{
		PostId: req.Id,
	})
	if err == nil {
		post.CommentsCount = int32(len(commentsResp.Comments))
	}

	return &pb.Post{
		Id:            fmt.Sprint(post.ID),
		UserId:        post.UserId,
		Content:       post.Content,
		ImageUrl:      post.ImageUrl,
		LikesCount:    post.LikesCount,
		CommentsCount: post.CommentsCount,
		CreatedAt:     post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *pb.DeletePostRequest) error {
	userID := contextx.GetUserID(ctx)
	if userID == "" {
		return status.Error(codes.Unauthenticated, "missing user id")
	}

	post, err := s.repo.GetPostByID(req.Id)
	if err != nil {
		return status.Error(codes.NotFound, "post not found")
	}

	// только владелец может удалить
	if post.UserId != userID {
		return status.Error(codes.PermissionDenied, "not your post")
	}

	if err := s.repo.DeletePostByID(req.Id); err != nil {
		return status.Errorf(codes.Internal, "failed to delete post: %v", err)
	}

	return nil
}

func (s *PostService) ListUserPosts(ctx context.Context, req *pb.UserPostsRequest) (*pb.Posts, error) {
	posts, err := s.repo.GetUserPosts(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load user posts: %v", err)
	}

	var pbPosts []*pb.Post
	for _, p := range posts {
		likesResp, _ := s.likeClient.LikePost(ctx, &likepb.LikePostRequest{
			Id: fmt.Sprint(p.ID),
		})
		commentsResp, _ := s.commClient.ListComments(ctx, &commentpb.ListCommentsRequest{
			PostId: fmt.Sprint(p.ID),
		})

		pbPosts = append(pbPosts, &pb.Post{
			Id:            fmt.Sprint(p.ID),
			UserId:        p.UserId,
			Content:       p.Content,
			ImageUrl:      p.ImageUrl,
			LikesCount:    likesResp.GetLikesCount(),
			CommentsCount: int32(len(commentsResp.GetComments())),
			CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &pb.Posts{Posts: pbPosts}, nil
}

func (s *PostService) GetAllPosts() (*pb.Posts, error) {
	posts, err := s.repo.GetAllPosts()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load posts: %v", err)
	}

	var pbPosts []*pb.Post
	for _, p := range posts {
		pbPosts = append(pbPosts, &pb.Post{
			Id:            fmt.Sprint(p.ID),
			UserId:        p.UserId,
			Content:       p.Content,
			ImageUrl:      p.ImageUrl,
			LikesCount:    p.LikesCount,
			CommentsCount: p.CommentsCount,
			CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return &pb.Posts{Posts: pbPosts}, nil
}

func (s *PostService) GetFeed(ctx context.Context, req *pb.GetFeedRequest) (*pb.Posts, error) {
	userID := contextx.GetUserID(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing user id")
	}

	followingResp, err := s.userClient.GetFollowing(ctx, &userpb.GetFollowingRequest{
		Id: userID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load following list: %v", err)
	}

	//  Собираем все ID: мой + друзей
	userIDs := []string{userID}
	for _, u := range followingResp.Users {
		userIDs = append(userIDs, u.Id)
	}

	//  Достаём посты этих пользователей
	posts, err := s.repo.GetPostsByUsers(userIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load posts: %v", err)
	}

	//  Формируем ответ
	var pbPosts []*pb.Post
	for _, p := range posts {
		pbPosts = append(pbPosts, &pb.Post{
			Id:            fmt.Sprint(p.ID),
			UserId:        p.UserId,
			Content:       p.Content,
			ImageUrl:      p.ImageUrl,
			LikesCount:    p.LikesCount,
			CommentsCount: p.CommentsCount,
			CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &pb.Posts{Posts: pbPosts}, nil
}
