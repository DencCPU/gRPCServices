package userhandlers

import (
	"context"
	"fmt"

	user "github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	domainuser "github.com/DencCPU/gRPCServices/UserService/internal/domain/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateUser(ctx context.Context, req *user.RegistrationUserReq) (*user.RegistrationUserResp, error) {

	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("incorrect data format:%w", err)
	}

	newUser := domainuser.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     domainuser.BasicRole,
	}

	pairToken, err := h.Service.CreateUser(ctx, newUser)

	resp := &user.RegistrationUserResp{
		AccessToken:  pairToken.AccessToken,
		RefreshToken: pairToken.RefreshToken,
		ExpireAt:     timestamppb.New(pairToken.ExpireAt),
	}
	return resp, nil
}
