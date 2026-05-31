package spotservice

import (
	"context"
	"errors"
	"fmt"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
	"github.com/DencCPU/gRPCServices/Protobuf/gen/common"
	spot "github.com/DencCPU/gRPCServices/Protobuf/gen/spot_service"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Client) GetEnableMarkets(ctx context.Context, userID string, userRole orderdomain.UserRole) ([]orderdomain.Market, error) {
	//Creating a breaker shell over a query
	result, err := s.breaker.Execute(func() (interface{}, error) {
		childCtx, cancel := context.WithTimeout(ctx, s.ConnectionTimeout)
		defer cancel()
		req := spot.ViewMarketReq{
			UserId:    userID,
			UserRoles: common.UserRole(userRole),
			PageSize:  0,
		}
		//Request to service
		resp, err := s.ViewMarket(childCtx, &req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	},
	)

	//Checking the status of the breaker (closed, open)
	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return nil, status.Errorf(codes.Unavailable, "service is temporarily unavailable")
		}
		return nil, err
	}

	//Casting the result of a safe query to a specific type
	resp, ok := result.(*spot.ViewMarketResp)
	if !ok {
		return nil, fmt.Errorf("Inappropriate result type:%T", result)
	}

	//Checking the response from the server
	if len(resp.EnableMarkets) == 0 {
		return []orderdomain.Market{}, nil
	}

	//Create the response
	var output = make([]orderdomain.Market, 0, len(resp.EnableMarkets))
	for _, em := range resp.EnableMarkets {
		market := orderdomain.Market{MarketId: em.MarketId, MarketName: em.MarketName}
		output = append(output, market)
	}
	return output, nil
}
