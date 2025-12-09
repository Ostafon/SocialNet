package handlers

import (
	"context"
	authpb "socialnet/services/auth/gen"
	pb "socialnet/services/notification/gen"
	"socialnet/services/notification/internal/service"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	svc *service.NotificationService
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) ListNotifications(ctx context.Context, req *pb.ListNotificationsRequest) (*pb.Notifications, error) {
	return h.svc.ListNotifications(ctx, req)
}

func (h *NotificationHandler) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*authpb.Confirmation, error) {
	return h.svc.MarkAsRead(ctx, req)
}

func (h *NotificationHandler) MarkAllAsRead(ctx context.Context, req *pb.EmptyRequest) (*authpb.Confirmation, error) {
	return h.svc.MarkAllAsRead(ctx, req)
}

func (h *NotificationHandler) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*authpb.Confirmation, error) {
	return h.svc.DeleteNotification(ctx, req)
}

func (h *NotificationHandler) ClearAll(ctx context.Context, req *pb.EmptyRequest) (*authpb.Confirmation, error) {
	return h.svc.ClearAll(ctx, req)
}

func (h *NotificationHandler) StreamNotifications(req *pb.StreamRequest, stream pb.NotificationService_StreamNotificationsServer) error {
	return h.svc.StreamNotifications(req, stream)
}
