package postgres

import (
	"context"
	"time"

	postgresdto "github.com/DencCPU/gRPCServices/UserService/internal/adapters/dto/postgres"
	userhash "github.com/DencCPU/gRPCServices/UserService/internal/adapters/hash"
	domainuser "github.com/DencCPU/gRPCServices/UserService/internal/domain/user"
	"github.com/google/uuid"
)

// Добавление нового пользоваткля
func (p *PostgresDB) AddUser(ctx context.Context, newUser domainuser.User) (string, string, error) {

	//Start of transaction
	tx, err := p.Begin(ctx)
	if err != nil {
		return "", "", err
	}

	//Create New DTO
	dto, err := postgresdto.NewUserDTO(newUser.Name, newUser.Email, newUser.Role)
	if err != nil {
		return "", "", err
	}
	hashPassword, err := userhash.HashPassword(newUser.Password)
	if err != nil {
		return "", "", err
	}
	dto.HashPassword = hashPassword
	dto.CreatedAt = time.Now()

	//Adding a new record to the database
	id := uuid.New()
	err = tx.QueryRow(ctx, `
	INSERT INTO users(id,name,email,hash_password,role,created_at)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id
	`,
		id,
		dto.Name,
		dto.Email,
		dto.HashPassword,
		dto.Role,
		dto.CreatedAt,
	).Scan(&dto.ID)

	if err != nil {
		return "", "", err
	}

	//Create a new refresh token
	token, err := p.AddRefreshToken(tx, ctx, dto.ID)
	if err != nil {

		return "", "", err
	}

	defer tx.Rollback(ctx)
	return dto.ID.String(), token, tx.Commit(ctx)
}

// Update password
func (p *PostgresDB) UpdatePassword(ctx context.Context, email, password string) error {
	dto := postgresdto.NewUpdatePassord(email, password)
	dto.UpdateAt = time.Now()

	_, err := p.Exec(ctx, `
	UPDATE users SET hash_password = $1
	WHERE email = $2
	`, dto.Email, dto.Password)

	if err != nil {
		return err
	}
	return nil
}

// Authentication
func (p *PostgresDB) Authentication(ctx context.Context, email, password string) (postgresdto.AuthUser, error) {
	var (
		hashPassword string
		output       postgresdto.AuthUser
	)

	err := p.QueryRow(ctx, `
	SELECT id,hash_password,role
	FROM users
	WHERE email = $1
	`, email).Scan(
		&output.ID,
		&hashPassword,
		&output.Role,
	)

	if err != nil {
		return postgresdto.AuthUser{}, err
	}

	if !userhash.CheckPasswordHash(password, hashPassword) {
		return postgresdto.AuthUser{}, err
	}

	return output, nil
}
