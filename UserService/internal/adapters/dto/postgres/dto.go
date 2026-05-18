package postgresdto

import (
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID           uuid.UUID
	Name         string
	Email        string
	HashPassword string
	Role         string
	CreatedAt    time.Time
}

type UpdatePassword struct {
	Email    string
	Password string
	UpdateAt time.Time
}

type RefreshToken struct {
	ID         int
	Token      string
	Expires_at time.Time
	IsRevoked  bool
	UserId     int
	CreatedAt  time.Time
}

type AuthUser struct {
	ID   string
	Role string
}

func NewUserDTO(name, email, role string) (*UserDTO, error) {

	dto := UserDTO{
		Name:  name,
		Email: email,
		Role:  role,
	}
	return &dto, nil
}

func NewUpdatePassord(email, password string) *UpdatePassword {
	return &UpdatePassword{Email: email, Password: password}
}
