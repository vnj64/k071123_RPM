package grpc

import (
	"context"
	"k071123/internal/services/order_service/contracts/pkg/proto"
	"k071123/internal/services/order_service/domain"
	"k071123/internal/services/order_service/domain/cases/card"
	"k071123/internal/services/order_service/domain/props"
	"time"
)

type GrpcHandler struct {
	proto.UnimplementedOrderServer
	ctx         domain.Context
	cardUseCase *card.CardUseCase
}

func NewHandler(ctx domain.Context, cardUseCase *card.CardUseCase) *GrpcHandler {
	return &GrpcHandler{
		ctx:         ctx,
		cardUseCase: cardUseCase,
	}
}

func (h *GrpcHandler) SaveCard(c context.Context, req *proto.SaveCardReq) (*proto.SaveCardResp, error) {
	resp, err := h.cardUseCase.SaveCard(props.SaveCardReq{
		UserUUID:      req.UserUuid,
		CardNumber:    req.CardNumber,
		Date:          FromUnixToTime(req.Date),
		CVC:           req.Cvc,
		PaymentSystem: req.PaymentSystem,
		IsPreferred:   req.IsPreferred,
		Email:         req.Email,
	})
	if err != nil {
		return nil, err
	}
	return &proto.SaveCardResp{
		Message: resp.Message,
	}, nil
}

func (h *GrpcHandler) GetPreferredByUserUUID(c context.Context, req *proto.GetPreferredCardReq) (*proto.GetPreferredCardResp, error) {
	ucReq := props.GetPreferredByUserUUIDReq{
		UserUUID: req.UserUuid,
	}
	ucResp, err := h.cardUseCase.GetPreferredByUserUUID(ucReq)
	if err != nil {
		return nil, err
	}
	resp := &proto.GetPreferredCardResp{
		Card: &proto.Card{
			Uuid:          ucResp.Card.UUID.String(),
			Last_4Digits:  ucResp.Card.Last4Digits,
			PaymentSystem: ucResp.Card.PaymentSystem,
			UserUuid:      ucResp.Card.UserUUID,
			IsPreferred:   ucResp.Card.IsPreferred,
			IsActive:      ucResp.Card.IsActive,
		},
	}
	if ucResp.Card.Token != nil {
		resp.Card.Token = *ucResp.Card.Token
	}
	return resp, nil
}

func FromUnixToTime(duration int64) time.Time {
	return time.Unix(duration, 0)
}
