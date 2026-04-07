package userclient

import (
	user "github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	"github.com/DencCPU/gRPCServices/UserService/pkg/userclient"
)

type Client struct {
	user.UserServiceClient
}

func NewClient() (*Client, error) {
	client, err := userclient.NewClient()
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}
