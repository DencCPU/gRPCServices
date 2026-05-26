package spotclient

import (
	"context"
	"errors"
	"fmt"

	spotservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/spot_service"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/common"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/spot_service"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) ViewEnableMarkets(ctx context.Context, input spotservicedto.Input) ([]spotservicedto.Output, error) {
	result, err := c.breaker.Execute(func() (interface{}, error) {
		req := spot_service.ViewMarketReq{
			UserId:    input.UserID,
			PageSize:  int32(input.PageSize),
			PageToken: input.PageToken,
		}

		switch input.UserRole {
		case "basic":
			req.UserRoles = common.UserRole_USER_ROLE_BASIC_USER
		case "premium":
			req.UserRoles = common.UserRole_USER_ROLE_PREMIUM_USER
		default:
			return nil, errors.New("unknow role")
		}

		resp, err := c.ViewMarket(ctx, &req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return nil, status.Errorf(codes.Unavailable, "service is temporarily unavailable")
		}
		return nil, err
	}

	resp, ok := result.(*spot_service.ViewMarketResp)
	if !ok {
		return nil, fmt.Errorf("Inappropriate result type:%T", result)
	}

	var output = make([]spotservicedto.Output, 0, len(resp.EnableMarkets))

	for _, el := range resp.EnableMarkets {
		var out_el spotservicedto.Output

		out_el.ID = el.MarketId
		out_el.Name = el.MarketName

		output = append(output, out_el)
	}
	return output, nil
}
