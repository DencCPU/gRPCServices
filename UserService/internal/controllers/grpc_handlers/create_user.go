package userhandlers

import (
	"context"
	"errors"
	"fmt"

	user "github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	domainuser "github.com/DencCPU/gRPCServices/UserService/internal/domain/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateUser(ctx context.Context, req *user.CreateUserReq) (*user.CreateUserResp, error) {

	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("incorrect data format:%w", err)
	}
	var role string

	switch req.UserRole {
	case user.UserRole_USER_ROLE_BASIC_USER:
		role = "basic"
	case user.UserRole_USER_ROLE_PREMIUM_USER:
		role = "premium"
	default:
		return nil, errors.New("user type does not exist")
	}

	newUser := domainuser.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     role,
	}
	pairToken, err := h.Service.CreateUser(ctx, newUser)

	resp := &user.CreateUserResp{
		AccessToken:  pairToken.AccessToken,
		RefreshToken: pairToken.RefreshToken,
		ExpireAt:     timestamppb.New(pairToken.Expire_at),
	}
	return resp, nil
}
