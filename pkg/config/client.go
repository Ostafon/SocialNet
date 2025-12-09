package config

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	commentpb "socialnet/services/comment/gen"
	likepb "socialnet/services/like/gen"
	postpb "socialnet/services/post/gen"
	userpb "socialnet/services/user/gen"
)

type GRPCClients struct {
	mu sync.Mutex

	likeConn    *grpc.ClientConn
	commentConn *grpc.ClientConn
	userConn    *grpc.ClientConn
	postConn    *grpc.ClientConn

	LikeClient    likepb.LikeServiceClient
	CommentClient commentpb.CommentServiceClient
	UserClient    userpb.UserServiceClient
	PostClient    postpb.PostServiceClient
}

// 🔹 Универсальная функция подключения
func dial(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to %s: %v", addr, err)
	}
	return conn, nil
}

// 🔹 Ленивое подключение к LikeService
func (c *GRPCClients) GetLikeClient(addr string) (likepb.LikeServiceClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.LikeClient == nil {
		conn, err := dial(addr)
		if err != nil {
			return nil, err
		}
		c.likeConn = conn
		c.LikeClient = likepb.NewLikeServiceClient(conn)
	}
	return c.LikeClient, nil
}

// 🔹 Ленивое подключение к CommentService
func (c *GRPCClients) GetCommentClient(addr string) (commentpb.CommentServiceClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.CommentClient == nil {
		conn, err := dial(addr)
		if err != nil {
			return nil, err
		}
		c.commentConn = conn
		c.CommentClient = commentpb.NewCommentServiceClient(conn)
	}
	return c.CommentClient, nil
}

// 🔹 Ленивое подключение к UserService
func (c *GRPCClients) GetUserClient(addr string) (userpb.UserServiceClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.UserClient == nil {
		conn, err := dial(addr)
		if err != nil {
			return nil, err
		}
		c.userConn = conn
		c.UserClient = userpb.NewUserServiceClient(conn)
	}
	return c.UserClient, nil
}

// 🔹 Ленивое подключение к PostService
func (c *GRPCClients) GetPostClient(addr string) (postpb.PostServiceClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.PostClient == nil {
		conn, err := dial(addr)
		if err != nil {
			return nil, err
		}
		c.postConn = conn
		c.PostClient = postpb.NewPostServiceClient(conn)
	}
	return c.PostClient, nil
}

// 🔹 Закрыть все соединения (при завершении приложения)
func (c *GRPCClients) CloseAll() {
	if c.likeConn != nil {
		c.likeConn.Close()
	}
	if c.commentConn != nil {
		c.commentConn.Close()
	}
	if c.userConn != nil {
		c.userConn.Close()
	}
	if c.postConn != nil {
		c.postConn.Close()
	}
}
