package userclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) AuthenticationUser(ctx context.Context, email, password string) (tokens.PairToken, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		childCtx, cancel := context.WithTimeout(ctx, c.ConnectionTimeout)
		defer cancel()
		req := user_service.AuthReq{
			Email:    email,
			Password: password,
		}
		resp, err := c.Authentication(childCtx, &req)
		if err != nil {
			return tokens.PairToken{}, err
		}
		return resp, nil
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return tokens.PairToken{}, status.Errorf(codes.Unavailable, "service is temporarily unavailable")
		}
		return tokens.PairToken{}, err
	}

	resp, ok := result.(*user_service.AuthResp)
	if !ok {
		return tokens.PairToken{}, fmt.Errorf("Inappropriate result type:%T", result)
	}

	pairToken := tokens.NewPairToken(resp.AccessToken, resp.RefreshToken, resp.ExpireAt.AsTime())
	return pairToken, nil

}
