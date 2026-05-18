package usecase

import (
	"context"

	viewdto "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/adapters/dto"
	spoterrors "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/errors"
	domainmarket "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/market"
	domainusers "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/users"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

// Получение доступных рынков
func (s *SpotService) ViewMarket(ctx context.Context, input domainusers.Input) ([]viewdto.Output, string, error) {

	var enableMarkets []*domainmarket.Market

	ctx, span := s.tracer.Start(ctx, "View markets")
	defer span.End()
	span.SetAttributes(
		attribute.String("userID", input.UserId),
		attribute.Int("PageSize", input.PageSize),
		attribute.String("PageToken", input.PageToken),
	)
	enableMarkets, pageToken := s.GetEnableMarkets(input)

	if len(enableMarkets) == 0 {
		s.logger.Error("no markets available")
		span.RecordError(spoterrors.Avalible_markets)
		span.SetStatus(codes.Error, "no markets available")
		return nil, "", spoterrors.Avalible_markets

	}

	s.logger.Info("List of available markets received",
		zap.String("spanID:", span.SpanContext().SpanID().String()),
	)
	span.SetStatus(codes.Ok, "view markets successfuly")
	return Mapper(enableMarkets), pageToken, nil
}

// Маппер для ViewMarket
func Mapper(em []*domainmarket.Market) []viewdto.Output {
	var resp []viewdto.Output
	for _, el := range em {
		resp = append(resp, viewdto.Output{ID: el.ID, Name: el.Name})
	}
	return resp
}
