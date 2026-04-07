package tokensdto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshToken struct {
	ID         int
	Token      string
	Expire_at  time.Time
	IsRevoked  bool
	UserId     uuid.UUID
	Created_at time.Time
	Update_at  time.Time
}

type PairToken struct {
	AccessToken  string
	RefreshToken string
	Expire_at    time.Time
}

type AccessClaim struct {
	User_id string `json:"user_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

type InputTokens struct {
	AccsesToken  string
	RefreshToken string
}

func NewRefreshToken(user_id uuid.UUID) *RefreshToken {
	return &RefreshToken{
		Token:     uuid.NewString(),
		Expire_at: time.Now().Add(5 * 24 * time.Hour),
		IsRevoked: false,
		UserId:    user_id,
	}
}

func NewPairToken(accessToken string, refreshToken string, expire_at time.Time) PairToken {
	return PairToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expire_at:    expire_at,
	}
}

func NewAccessClaim(user_id, email, role string) *AccessClaim {
	return &AccessClaim{
		User_id: user_id,
		Email:   email,
		Role:    role,
	}
}

func NewInputTokens(accessToken string, refreshToken string) InputTokens {
	return InputTokens{
		AccsesToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
