package userclient

import (
	"context"

	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	user "github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
)

func (c *Client) UpdateAccessToken(ctx context.Context, accessToken, refreshToken string) (tokens.PairToken, error) {
	req := user.UpdateTokensReq{AccessToken: accessToken, RefreshToken: refreshToken}
	resp, err := c.UpdateTokens(ctx, &req)
	if err != nil {
		return tokens.PairToken{}, err
	}
	pairToken := tokens.NewPairToken(resp.AccessToken, resp.RefreshToken, resp.ExpireAt.AsTime())
	return pairToken, nil
}
