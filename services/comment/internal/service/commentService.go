package service

import (
	"context"
	"socialnet/pkg/utils"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "socialnet/services/comment/gen"
	"socialnet/services/comment/internal/model"
	"socialnet/services/comment/internal/repos"
)

type CommentService struct {
	repo *repos.CommentRepo
}

func NewCommentService(repo *repos.CommentRepo) *CommentService {
	return &CommentService{repo: repo}
}

// Добавление комментария
func (s *CommentService) AddComment(ctx context.Context, userID string, req *pb.AddCommentRequest) (*pb.Comment, error) {
	comment := &model.Comment{
		PostID:  req.PostId,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.repo.AddComment(comment); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add comment: %v", err)
	}
	id := utils.UintToString(comment.ID)
	return &pb.Comment{
		Id:         id,
		PostId:     comment.PostID,
		UserId:     comment.UserID,
		Content:    comment.Content,
		LikesCount: int32(comment.LikesCount),
		CreatedAt:  comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  comment.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// Получение комментария
func (s *CommentService) GetComment(ctx context.Context, id string) (*pb.Comment, error) {
	c, err := s.repo.GetComment(id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "comment not found")
	}

	idc := utils.UintToString(c.ID)
	return &pb.Comment{
		Id:         idc,
		PostId:     c.PostID,
		UserId:     c.UserID,
		Content:    c.Content,
		LikesCount: int32(c.LikesCount),
		CreatedAt:  c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  c.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// Удаление комментария
func (s *CommentService) DeleteComment(ctx context.Context, id, userID string) error {
	if err := s.repo.DeleteComment(id, userID); err != nil {
		return status.Error(codes.Internal, "failed to delete comment")
	}
	return nil
}

// Список комментариев к посту
func (s *CommentService) ListComments(ctx context.Context, postID string) (*pb.Comments, error) {
	comments, err := s.repo.ListComments(postID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get comments")
	}

	resp := make([]*pb.Comment, 0, len(comments))
	for _, c := range comments {
		id := utils.UintToString(c.ID)
		resp = append(resp, &pb.Comment{
			Id:         id,
			PostId:     c.PostID,
			UserId:     c.UserID,
			Content:    c.Content,
			LikesCount: int32(c.LikesCount),
			CreatedAt:  c.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  c.UpdatedAt.Format(time.RFC3339),
		})
	}
	return &pb.Comments{Comments: resp}, nil
}
