package spotclient

import (
	"context"
	"errors"

	spotservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/spot_service"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/spot_service"
)

func (c *Client) ViewEnableMarkets(ctx context.Context, role string) ([]spotservicedto.Output, error) {
	req := spot_service.ViewReq{}

	switch role {
	case "basic":
		req.UserRoles = spot_service.UserRole_USER_ROLE_BASIC_USER
	case "premium":
		req.UserRoles = spot_service.UserRole_USER_ROLE_PREMIUM_USER
	default:
		return nil, errors.New("unknow role")
	}

	resp, err := c.ViewMarket(ctx, &req)
	if err != nil {
		return nil, err
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
