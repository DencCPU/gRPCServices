package userclient

import (
	"context"
	"errors"
	"fmt"

	userservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/user_service"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/common"
	user "github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) Validation(ctx context.Context, accessToken string) (userservicedto.Output, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		childCtx, cancel := context.WithTimeout(ctx, c.ConnectionTimeout)
		defer cancel()

		req := user.ValidationReq{AccessToken: accessToken}

		resp, err := c.ValidationTokens(childCtx, &req)
		if err != nil {
			return userservicedto.Output{}, err
		}

		err = resp.Validate()
		if err != nil {
			return userservicedto.Output{}, errors.New("Invalid server respons format")
		}
		return resp, nil
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return userservicedto.Output{}, status.Errorf(codes.Unavailable, "service is temporarily unavailable")
		}
		return userservicedto.Output{}, err
	}

	resp, ok := result.(*user.ValidationResp)
	if !ok {
		return userservicedto.Output{}, fmt.Errorf("Inappropriate result type:%T", result)
	}

	var output userservicedto.Output
	output.User_id = resp.UserId
	switch resp.Role {
	case common.UserRole_USER_ROLE_BASIC_USER:
		output.Role = "basic"
	case common.UserRole_USER_ROLE_PREMIUM_USER:
		output.Role = "premium"
	}
	return output, nil
}
