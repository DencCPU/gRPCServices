package spothandlers

import (
	"context"
	"fmt"

	spot "github.com/DencCPU/gRPCServices/Protobuf/gen/spot_service"
	domainusers "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/users"
)

func (h *Handlers) ViewMarket(ctx context.Context, req *spot.ViewReq) (*spot.ViewResp, error) {
	//Валидация
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("неправильный формат запроса:%w", err)
	}

	//Преобразование запроса в доменную структуру User
	user := domainusers.NewUser(domainusers.UserType(req.UserRoles))

	//Запрос сервиса на получение доступных рынков
	output, err := h.Service.ViewMarket(ctx, user)
	if err != nil {
		return &spot.ViewResp{}, err
	}

	//Формирование ответа
	resp := &spot.ViewResp{}
	resp.EnableMarkets = make([]*spot.Markets, 0, len(output))

	for _, el := range output {
		market := spot.Markets{MarketId: el.ID, MarketName: el.Name}
		resp.EnableMarkets = append(resp.EnableMarkets, &market)
	}
	return resp, nil
}
