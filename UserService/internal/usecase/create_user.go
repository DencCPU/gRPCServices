package usecase

import (
	"context"

	tokensdto "github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/tokens"
	domainuser "github.com/DencCPU/gRPCServices/UserService/internal/domain/user"
	"go.uber.org/zap"
)

func (s *Service) CreateUser(ctx context.Context, user domainuser.User) (tokensdto.PairToken, error) {

	//Add user to database and generate refresh token
	ctx, span := s.tracer.Start(ctx, "Add user:")
	defer span.End()
	user_id, refreshToken, err := s.AddUser(ctx, user)
	if err != nil {
		s.logger.Error("error adding user to database:",
			zap.Error(err),
		)
		return tokensdto.PairToken{}, err
	}
	s.logger.Info("adding user succefully")

	//Create accses token
	ctx, span = s.tracer.Start(ctx, "Create accsess token:")
	defer span.End()
	accsesToken, ttl, err := s.CreateAccessToken(user_id, user.Email, user.Role)
	if err != nil {
		s.logger.Error("jwt generation error:",
			zap.Error(err),
		)
		return tokensdto.PairToken{}, err
	}

	pairToken := tokensdto.NewPairToken(accsesToken, refreshToken, ttl)
	s.logger.Info("token creation succeful")
	return pairToken, nil
}
