package postgres

import (
	"context"
	"errors"
	"time"

	tokensdto "github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/tokens"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Adding a new refresh token to database
func (p *PostgresDB) AddRefreshToken(tx pgx.Tx, ctx context.Context, user_id uuid.UUID) (string, error) {

	//Create a new refresh token
	refreshToken := tokensdto.NewRefreshToken(user_id)

	var token = refreshToken.Token
	refreshToken.Created_at = time.Now()

	_, err := tx.Exec(ctx, `
	INSERT INTO refresh_tokens(token,expire_at,is_revoked,user_id,created_at)
	VALUES($1,$2,$3,$4,$5) 
	RETURNING token 
	`,
		refreshToken.Token,
		refreshToken.Expire_at,
		refreshToken.IsRevoked,
		refreshToken.UserId,
		refreshToken.Created_at,
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Update token data
func (p *PostgresDB) UpdateRefreshToken(ctx context.Context, token string) (string, error) {

	tx, err := p.Begin(ctx)
	if err != nil {
		return "", err
	}

	var rToken tokensdto.RefreshToken
	err = tx.QueryRow(ctx, `
	  SELECT id, token, expire_at, is_revoked 
        FROM refresh_tokens
        WHERE token = $1
	`, token).Scan(
		&rToken.ID,
		&rToken.Token,
		&rToken.Expire_at,
		&rToken.IsRevoked)
	if err != nil {
		return "", err
	}

	if rToken.IsRevoked {
		return "", errors.New("user is blocked")
	}

	if time.Now().After(rToken.Expire_at) {
		return "", errors.New("re-autentification required")
	}

	rToken.Token = uuid.NewString()
	rToken.Expire_at = time.Now().Add(5 * 24 * time.Hour)
	rToken.Update_at = time.Now()

	_, err = tx.Exec(ctx, `
	UPDATE refresh_tokens SET 
	token=$1,
	expire_at=$2,
	update_at =$3
	WHERE id = $4
	`,
		rToken.Token,
		rToken.Expire_at,
		rToken.Update_at,
		rToken.ID)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	return rToken.Token, tx.Commit(ctx)
}
