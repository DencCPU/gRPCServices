package userclient

import (
	"context"
	"errors"

	userservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/user_service"
	user "github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
)

func (c *Client) Validation(ctx context.Context, accessToken string) (userservicedto.Output, error) {
	req := user.ValidationReq{AccessToken: accessToken}
	resp, err := c.ValidationTokens(ctx, &req)
	if err != nil {
		return userservicedto.Output{}, err
	}

	err = resp.Validate()
	if err != nil {
		return userservicedto.Output{}, errors.New("Invalid server respons format")
	}

	var output userservicedto.Output
	output.User_id = resp.UserId
	switch resp.Role {
	case user.UserRole_USER_ROLE_BASIC_USER:
		output.Role = "basic"
	case user.UserRole_USER_ROLE_PREMIUM_USER:
		output.Role = "premium"
	}
	return output, nil
}
