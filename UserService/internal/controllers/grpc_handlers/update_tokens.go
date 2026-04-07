package userhandlers

import (
	"context"
	"fmt"

	user "github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	tokensdto "github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/tokens"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) UpdateTokens(ctx context.Context, req *user.UpdateTokensReq) (*user.UpdateTokensResp, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("incorrect data format:%w", err)
	}

	input := tokensdto.NewInputTokens(req.AccessToken, req.RefreshToken)
	pairToken, err := h.Service.UpdateTokens(ctx, input)
	if err != nil {
		return &user.UpdateTokensResp{}, err
	}

	resp := &user.UpdateTokensResp{
		AccessToken:  pairToken.AccessToken,
		RefreshToken: pairToken.RefreshToken,
		ExpireAt:     timestamppb.New(pairToken.Expire_at),
	}
	return resp, nil
}
