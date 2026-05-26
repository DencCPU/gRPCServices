package userclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	user "github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) UpdateAccessToken(ctx context.Context, accessToken, refreshToken string) (tokens.PairToken, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		req := user.UpdateTokensReq{AccessToken: accessToken, RefreshToken: refreshToken}

		resp, err := c.UpdateTokens(ctx, &req)
		if err != nil {
			return tokens.PairToken{}, err
		}

		if resp.AccessToken == "" || resp.RefreshToken == "" {
			return tokens.PairToken{}, errors.New("incorrect response from the server")
		}

		return resp, nil
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return tokens.PairToken{}, status.Errorf(codes.Unavailable, "service is temporarily unavailable")
		}

		return tokens.PairToken{}, err
	}

	resp, ok := result.(*user.UpdateTokensResp)
	if !ok {
		return tokens.PairToken{}, fmt.Errorf("Inappropriate result type:%T", result)
	}

	pairToken := tokens.NewPairToken(resp.AccessToken, resp.RefreshToken, resp.ExpireAt.AsTime())
	return pairToken, nil
}
