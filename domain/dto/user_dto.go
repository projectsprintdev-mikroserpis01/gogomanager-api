package dto

import (
	"github.com/google/uuid"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type GetUsersRequest struct {
}

type GetUsersQuery struct {
	Limit          int    `query:"limit" validate:"omitempty,number,gte=1,lte=100"`
	Page           int    `query:"page" validate:"omitempty,number,gte=1"`
	SortBy         string `query:"sort_by" validate:"omitempty,oneof=created_at updated_at name email id"`
	Order          string `query:"order" validate:"omitempty,oneof=asc desc"`
	IncludeDeleted bool   `query:"include_deleted"`
	Search         string `query:"search"`
}

type GetUsersResponse struct {
	Users []entity.User `json:"users"`
}

type GetUserByIDRequest struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

type GetUserByIDResponse struct {
	User entity.User `json:"user"`
}

type GetUsersStatsRequest struct {
}

type GetUsersStatsQuery struct {
	IncludeDeleted bool   `json:"include_deleted"`
	Search         string `json:"search"`
}

type GetUsersStatsResponse struct {
	TotalNonDeletedUsers int64 `json:"total_non_deleted_users"`
	TotalDeletedUsers    int64 `json:"total_deleted_users"`
	TotalUsers           int64 `json:"total_users"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100,ascii"`
	Password string `json:"password" validate:"required,min=8,max=100,ascii"`
	Email    string `json:"email" validate:"required,email"`
	RoleID   int    `json:"role_id" validate:"number,gte=1"`
}

type CreateUserResponse struct {
	ID uuid.UUID `json:"id"`
}

type UpdateUserRequest struct {
	ID       uuid.UUID `param:"id" validate:"required,uuid"`
	Name     string    `json:"name" validate:"required,min=3,max=100,ascii"`
	Password string    `json:"password" validate:"required,min=8,max=100,ascii"`
	Email    string    `json:"email" validate:"required,email"`
}

type UpdateUserResponse struct {
	ID uuid.UUID `json:"id"`
}

type SoftDeleteUserRequest struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

type SoftDeleteUserResponse struct {
	ID uuid.UUID `json:"id"`
}

type DeleteUserRequest struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

type DeleteUserResponse struct {
	ID uuid.UUID `json:"id"`
}

type RestoreUserRequest struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

type RestoreUserResponse struct {
	ID uuid.UUID `json:"id"`
}
