package usecase

import (
	"context"

	"github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/dto/tokens"
	userdomain "github.com/DencCPU/gRPCServices/APIGetway/internal/domain/user"
	"go.uber.org/zap"
)

func (s *Service) RegistrationUser(ctx context.Context, newUser userdomain.User) (tokens.PairToken, error) {
	pairToken, err := s.user_client.RegistrationUser(ctx, newUser)
	if err != nil {
		s.logger.Error("ошибка регистрации нового пользователя на user-сервисе:",
			zap.Error(err),
		)
		return tokens.PairToken{}, err
	}
	return pairToken, nil
}
