package jwt

import (
	"errors"
	"time"

	sharederrors "github.com/DencCPU/gRPCServices/Shared/errors"
	userconfig "github.com/DencCPU/gRPCServices/UserService/config"
	tokensdto "github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/tokens"
	"github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/user"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret         []byte
	AccessTokenTTL time.Duration
}

func NewJWT(cfg userconfig.JWT) *JWT {
	secret := cfg.Secret
	ttl := cfg.TTL
	return &JWT{secret: []byte(secret), AccessTokenTTL: time.Minute * time.Duration(ttl)}
}

// Create a new JWT
func (j *JWT) CreateAccessToken(user_id, email, role string) (string, time.Time, error) {
	accessClaim := tokensdto.NewAccessClaim(user_id, email, role)
	accessClaim.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.AccessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaim).SignedString(j.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return accessToken, accessClaim.ExpiresAt.Time, nil
}

// Update JWT
func (j *JWT) UpdateAccessToken(oldAccessToken string) (string, time.Time, error) {
	claims := &tokensdto.AccessClaim{}
	_, err := jwt.ParseWithClaims(oldAccessToken, claims, func(t *jwt.Token) (any, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unknow signature creation method")
		}
		return j.secret, nil
	})
	if err != nil {
		if !errors.Is(err, jwt.ErrTokenExpired) {
			return "", time.Time{}, err
		}
	}

	newAccessClaim := tokensdto.NewAccessClaim(claims.User_id, claims.Email, claims.Role)
	newAccessClaim.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.AccessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	accsessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaim).SignedString(j.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return accsessToken, newAccessClaim.ExpiresAt.Time, nil
}

// Validation JWT
func (j *JWT) Validation(accessToken string) (user.Output, error) {
	claims := &tokensdto.AccessClaim{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (any, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unknow signature creation method")
		}
		return j.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return user.Output{}, sharederrors.EXPIRED_TOKEN
		}
		return user.Output{}, err
	}

	if !token.Valid {
		return user.Output{}, errors.New("invalid token")
	}

	user := user.Output{
		User_id: claims.User_id,
		Role:    claims.Role,
	}

	return user, nil
}
