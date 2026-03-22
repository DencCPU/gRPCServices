package spothandlers

import (
	spot "Academy/gRPCServices/Protobuf/gen/spot_service"
	viewdto "Academy/gRPCServices/SpotInstrumentService/internal/adapters/dto"
	domainusers "Academy/gRPCServices/SpotInstrumentService/internal/domain/users"
	"Academy/gRPCServices/SpotInstrumentService/internal/usecase"
	"context"
)

type Service interface {
	ViewMarket(context.Context, *domainusers.User) ([]viewdto.Output, error)
}

type Handlers struct {
	spot.UnimplementedSpotInstrumentServiceServer
	Service Service //Функционал обработчиков
}

// Конструктор для SpotInstrument
func NewHandlers(spotService *usecase.SpotService) *Handlers {
	return &Handlers{Service: spotService}
}
