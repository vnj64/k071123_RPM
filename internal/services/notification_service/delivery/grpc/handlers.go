package grpc

import (
	"context"
	"fmt"
	"k071123/internal/services/notification_service/contracts/pkg/proto"
	"k071123/internal/services/notification_service/domain"
	"k071123/internal/services/notification_service/domain/cases/sender"
	"k071123/internal/services/notification_service/domain/props"
)

type GrpcHandler struct {
	proto.UnimplementedNotificationServer
	ctx           domain.Context
	senderUseCase *sender.EmailSenderUseCase
}

func NewHandler(ctx domain.Context, senderUseCase *sender.EmailSenderUseCase) *GrpcHandler {
	return &GrpcHandler{
		ctx:           ctx,
		senderUseCase: senderUseCase,
	}
}

func (h *GrpcHandler) SendEmail(c context.Context, req *proto.SendEmailReq) (*proto.SendEmailResp, error) {
	resp, err := h.senderUseCase.SendEmail(&props.SendEmailRequest{
		Subject: req.Subject,
		To:      req.To,
		Data:    req.Data,
	})
	if err != nil {
		return nil, err
	}

	if resp.Status == "failed" {
		return nil, fmt.Errorf("failed to send email. email status: %s", resp.Status)
	}

	return &proto.SendEmailResp{
		Response: resp.Status,
	}, nil
}
