package tokens

import "time"

type PairToken struct {
	AccessToken  string
	RefreshToken string
	ExpireAt     time.Time
}

func NewPairToken(accessToken string, refreshToken string, expire_at time.Time) PairToken {
	return PairToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpireAt:     expire_at,
	}
}
