package usecase

import (
	"context"
	"errors"

	"github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/user"
	"go.uber.org/zap"
)

func (s *Service) Validation(ctx context.Context, accessToken string) (user.Output, error) {
	ctx, span := s.tracer.Start(ctx, "validation token")
	defer span.End()
	output, err := s.JWT.Validation(accessToken)
	if err != nil {
		s.logger.Error("token validation error:",
			zap.Error(err),
		)
		return user.Output{}, err
	}
	if output.Role == "" || output.User_id == "" {
		s.logger.Error("data from token was not received")
		return user.Output{}, errors.New("data from token was not received")
	}
	s.logger.Info("token validation successful")
	return output, nil
}
