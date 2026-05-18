package usecase

import (
	"context"

	spotservicedto "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/spot_service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (s *Service) ViewEnableMarkets(ctx context.Context, input spotservicedto.Input) ([]spotservicedto.Output, error) {
	ctx, span := s.tracer.Start(ctx, "View markets")
	defer span.End()
	span.SetAttributes(
		attribute.String("userID", input.UserID),
		attribute.Int("PageSize", input.PageSize),
		attribute.String("PageToken", input.PageToken),
	)

	markets, err := s.spotClient.ViewEnableMarkets(ctx, input)
	if err != nil {
		s.logger.Error("error getting list of markets",
			zap.String("spanID:", span.SpanContext().SpanID().String()),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "view enable market failed")
		return nil, err
	}

	s.logger.Info("list of markets received")

	span.SetStatus(codes.Ok, "view enable markets successfuly")
	return markets, nil
}
