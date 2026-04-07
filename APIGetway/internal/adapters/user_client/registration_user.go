package userclient

import (
	"context"

	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	userdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/user"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
)

func (c *Client) RegistrationUser(ctx context.Context, newUser userdomain.User) (tokens.PairToken, error) {
	req := &user_service.CreateUserReq{
		Name:     newUser.Name,
		Email:    newUser.Email,
		Password: newUser.Password,
		UserRole: user_service.UserRole_USER_ROLE_BASIC_USER,
	}
	resp, err := c.CreateUser(ctx, req)
	if err != nil {
		return tokens.PairToken{}, err
	}
	pairToken := tokens.NewPairToken(
		resp.AccessToken,
		resp.RefreshToken,
		resp.ExpireAt.AsTime(),
	)
	return pairToken, nil
}
