package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"socialnet/pkg/config"
	"socialnet/pkg/utils"
	pb "socialnet/services/like/gen"
	"socialnet/services/like/internal/repos"
	notificationpb "socialnet/services/notification/gen"
	postpb "socialnet/services/post/gen"
)

type LikeService struct {
	repo    *repos.LikeRepo
	clients *config.GRPCClients
}

func NewLikeService(repo *repos.LikeRepo, cl *config.GRPCClients) *LikeService {
	return &LikeService{repo: repo, clients: cl}
}

// POST LIKES
func (s *LikeService) LikePost(ctx context.Context, userID, postID string) (*pb.LikePostResponse, error) {
	exists, err := s.repo.HasUserLikedPost(userID, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check like: %v", err)
	}
	if exists {
		if err := s.repo.UnlikePost(userID, postID); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to remove like: %v", err)
		}
		count, _ := s.repo.CountPostLikes(postID)
		return &pb.LikePostResponse{Status: "unliked", LikesCount: int32(count)}, nil
	}

	if err := s.repo.LikePost(userID, postID); err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "post already liked or invalid: %v", err)
	}
	postClient, err := s.clients.GetPostClient("localhost:50053")
	if err == nil {
		postResp, err := postClient.GetPost(ctx, &postpb.GetPostRequest{Id: postID})
		if err == nil {
			postOwnerID := postResp.UserId
			notifClient, err := s.clients.GetNotifClient("localhost:50057")
			if err == nil {

				md := metadata.New(map[string]string{"user-id": userID})
				ctxWithUser := metadata.NewOutgoingContext(ctx, md)

				_, _ = notifClient.CreateNotification(ctxWithUser,
					&notificationpb.CreateNotificationRequest{
						UserId:      postOwnerID,
						Type:        "post_liked",
						ReferenceId: postID,
						Content:     fmt.Sprintf("Your post was liked by %s", userID),
					})
			}
		}
	}

	count, _ := s.repo.CountPostLikes(postID)
	return &pb.LikePostResponse{Status: "liked", LikesCount: int32(count)}, nil
}

func (s *LikeService) UnlikePost(ctx context.Context, userID, postID string) (*pb.LikePostResponse, error) {
	if err := s.repo.UnlikePost(userID, postID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unlike: %v", err)
	}

	count, _ := s.repo.CountPostLikes(postID)
	return &pb.LikePostResponse{Status: "unliked", LikesCount: int32(count)}, nil
}

func (s *LikeService) ListPostLikes(ctx context.Context, postID string) (*pb.ListLikesResponse, error) {
	likes, err := s.repo.ListPostLikes(postID)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Like, 0, len(likes))
	for _, l := range likes {
		id := utils.UintToString(l.ID)
		res = append(res, &pb.Like{
			Id:        id,
			UserId:    l.UserID,
			PostId:    *l.PostID,
			CreatedAt: l.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return &pb.ListLikesResponse{Likes: res}, nil
}

// COMMENT LIKES
func (s *LikeService) LikeComment(ctx context.Context, userID, commentID string) (*pb.LikeCommentResponse, error) {
	if err := s.repo.LikeComment(userID, commentID); err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "comment already liked or invalid: %v", err)
	}
	count, _ := s.repo.CountCommentLikes(commentID)
	return &pb.LikeCommentResponse{Status: "liked", LikesCount: int32(count)}, nil
}

func (s *LikeService) UnlikeComment(ctx context.Context, userID, commentID string) (*pb.LikeCommentResponse, error) {
	if err := s.repo.UnlikeComment(userID, commentID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unlike comment: %v", err)
	}
	count, _ := s.repo.CountCommentLikes(commentID)
	return &pb.LikeCommentResponse{Status: "unliked", LikesCount: int32(count)}, nil
}

func (s *LikeService) ListCommentLikes(ctx context.Context, commentID string) (*pb.ListLikesResponse, error) {
	likes, err := s.repo.ListCommentLikes(commentID)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Like, 0, len(likes))
	for _, l := range likes {
		id := utils.UintToString(l.ID)
		res = append(res, &pb.Like{
			Id:        id,
			UserId:    l.UserID,
			CommentId: *l.CommentID,
			CreatedAt: l.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return &pb.ListLikesResponse{Likes: res}, nil
}
