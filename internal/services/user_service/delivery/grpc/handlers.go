package grpc

import (
	"context"
	"k071123/internal/services/user_service/contracts/pkg/proto"
	"k071123/internal/services/user_service/domain"
	"k071123/internal/services/user_service/domain/cases/user"
	"k071123/internal/services/user_service/domain/props"
	"log"
)

type GrpcHandler struct {
	proto.UnimplementedUserServer
	ctx         domain.Context
	userUseCase *user.UserUseCase
}

func NewHandler(ctx domain.Context, userUseCase *user.UserUseCase) *GrpcHandler {
	return &GrpcHandler{
		ctx:         ctx,
		userUseCase: userUseCase,
	}
}

func (h *GrpcHandler) GetUserByUUID(c context.Context, req *proto.GetUserReq) (response *proto.GetUserResp, err error) {
	resp, err := h.userUseCase.GetUserByUUID(props.GetUserByUUIDReq{
		UUID: req.Uuid,
	})
	response = &proto.GetUserResp{}
	if err != nil {
		return response, err
	}

	if resp.User == nil {
		log.Printf("USER IS NIL ON GRPC HANDLER")
		return response, nil
	}
	if resp.User.FirstName != nil {
		response.FirstName = *resp.User.FirstName
	}
	if resp.User.SecondName != nil {
		response.SecondName = *resp.User.SecondName
	}
	if resp.User.PhoneNumber != nil {
		response.PhoneNumber = *resp.User.PhoneNumber
	}
	if resp.User.BirthDate != nil {
		response.BirthDate = resp.User.BirthDate.String()
	}
	log.Printf("full response 1: %+v", resp.User)
	response.Status = string(resp.User.Status)
	response.Uuid = resp.User.UUID.String()
	response.Email = resp.User.Email
	response.Role = string(resp.User.Role)
	log.Printf("full response 2: %+v", resp.User)
	return response, nil
}
