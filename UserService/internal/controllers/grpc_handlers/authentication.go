package userhandlers

import (
	"context"
	"fmt"

	"github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) Authentication(ctx context.Context, req *user_service.AuthReq) (*user_service.AuthResp, error) {

	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("incorrect data format:%w", err)
	}
	var (
		email    string = req.Email
		password string = req.Password
	)

	pairToken, err := h.Service.AuthenticationUser(ctx, email, password)
	if err != nil {
		return nil, err
	}
	resp := user_service.AuthResp{
		AccessToken:  pairToken.AccessToken,
		RefreshToken: pairToken.RefreshToken,
		ExpireAt:     timestamppb.New(pairToken.ExpireAt),
	}
	return &resp, nil
}
