package userclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	userdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/user"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) RegistrationUser(ctx context.Context, newUser userdomain.User) (tokens.PairToken, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		req := &user_service.RegistrationUserReq{
			Name:     newUser.Name,
			Email:    newUser.Email,
			Password: newUser.Password,
		}

		resp, err := c.RegistrationNewUser(ctx, req)
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
	resp, ok := result.(*user_service.RegistrationUserResp)
	if !ok {
		return tokens.PairToken{}, fmt.Errorf("Inappropriate result type:%T", result)
	}

	pairToken := tokens.NewPairToken(
		resp.AccessToken,
		resp.RefreshToken,
		resp.ExpireAt.AsTime(),
	)
	return pairToken, nil
}
