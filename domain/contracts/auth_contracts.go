package contracts

import (
	"context"

	"github.com/google/uuid"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	RegisterUser(ctx context.Context, user entity.User) (uuid.UUID, error)
}

type AuthService interface {
	RegisterUser(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
	LoginUser(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
}
