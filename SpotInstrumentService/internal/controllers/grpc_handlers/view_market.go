package spothandlers

import (
	spot "Academy/gRPCServices/Protobuf/gen/spot_service"
	domainusers "Academy/gRPCServices/SpotInstrumentService/internal/domain/users"
	"context"
)

func (h *Handlers) ViewMarket(ctx context.Context, req *spot.ViewReq) (*spot.ViewResp, error) {
	user := domainusers.NewUser(domainusers.UserType(req.UserRoles)) //Преобразование запроса в доменную структуру User

	output, err := h.Service.ViewMarket(ctx, user) //Запрос сервиса на получение доступных рынков
	if err != nil {
		return &spot.ViewResp{}, err
	}

	resp := &spot.ViewResp{}
	resp.EnableMarkets = make([]*spot.Markets, 0, len(output))

	for _, el := range output {
		market := spot.Markets{MarketId: el.ID, MarketName: el.Name}
		resp.EnableMarkets = append(resp.EnableMarkets, &market)
	}
	//Формирование ответа сервера
	return resp, nil
}
