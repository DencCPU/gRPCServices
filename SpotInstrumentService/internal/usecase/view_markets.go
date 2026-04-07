package usecase

import (
	"context"

	viewdto "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/adapters/dto"
	spoterrors "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/errors"
	domainmarket "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/market"
	domainusers "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/users"
	"go.uber.org/zap"
)

// Получение доступных рынков
func (s *SpotService) ViewMarket(ctx context.Context, user *domainusers.User) ([]viewdto.Output, error) {

	var enableMarkets []*domainmarket.Market
	ctx, span := s.tracer.Start(ctx, "view markets")
	defer span.End()

	enableMarkets = s.GetEnableMarkets()

	if len(enableMarkets) == 0 {
		s.logger.Error("no markets available")
		return nil, spoterrors.Avalible_markets
	}
	s.logger.Info("List of available markets received",
		zap.String("get enable markets span:", span.SpanContext().TraceID().String()),
	)
	return Mapper(enableMarkets), nil
}

// Маппер для ViewMarket
func Mapper(em []*domainmarket.Market) []viewdto.Output {
	var resp []viewdto.Output
	for _, el := range em {
		resp = append(resp, viewdto.Output{ID: el.ID, Name: el.Name})
	}
	return resp
}
