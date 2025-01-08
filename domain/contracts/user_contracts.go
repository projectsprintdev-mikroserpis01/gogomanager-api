package contracts

import (
	"context"

	"github.com/google/uuid"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type UserRepository interface {
	BeginTransaction(ctx context.Context) error
	CommitTransaction() error
	RollbackTransaction() error

	GetUsers(ctx context.Context, query dto.GetUsersQuery) ([]entity.User, error)
	GetUserByField(ctx context.Context, field, value string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) (uuid.UUID, error)
	UpdateUser(ctx context.Context, user *entity.User) (uuid.UUID, error)
	SoftDeleteUser(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	DeleteUser(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	RestoreUser(ctx context.Context, id uuid.UUID) (uuid.UUID, error)

	CountUsers(ctx context.Context, query dto.GetUsersStatsQuery) (int64, error)
}

type UserService interface {
	GetUsers(ctx context.Context, query dto.GetUsersQuery) (dto.GetUsersResponse, error)
	GetUserByID(ctx context.Context, req dto.GetUserByIDRequest) (dto.GetUserByIDResponse, error)
	GetUsersStats(ctx context.Context) (dto.GetUsersStatsResponse, error)
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error)
	UpdateUser(ctx context.Context, req dto.UpdateUserRequest) (dto.UpdateUserResponse, error)
	SoftDeleteUser(ctx context.Context, req dto.SoftDeleteUserRequest) (dto.SoftDeleteUserResponse, error)
	DeleteUser(ctx context.Context, req dto.DeleteUserRequest) (dto.DeleteUserResponse, error)
	RestoreUser(ctx context.Context, req dto.RestoreUserRequest) (dto.RestoreUserResponse, error)
}
