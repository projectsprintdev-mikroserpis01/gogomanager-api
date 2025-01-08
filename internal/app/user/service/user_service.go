package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/bcrypt"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/uuid"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/validator"
)

type userService struct {
	userRepo  contracts.UserRepository
	validator validator.ValidatorInterface
	uuid      uuid.UUIDInterface
	bcrypt    bcrypt.BcryptInterface
}

func NewUserService(
	userRepo contracts.UserRepository,
	validator validator.ValidatorInterface,
	uuid uuid.UUIDInterface,
	bcrypt bcrypt.BcryptInterface,
) contracts.UserService {
	return &userService{
		userRepo:  userRepo,
		validator: validator,
		uuid:      uuid,
		bcrypt:    bcrypt,
	}
}

func (s *userService) GetUsers(ctx context.Context, query dto.GetUsersQuery) (dto.GetUsersResponse, error) {
	valErr := s.validator.Validate(query)
	if valErr != nil {
		return dto.GetUsersResponse{}, valErr
	}

	if query.Limit == 0 {
		query.Limit = 10
	}

	if query.Page == 0 {
		query.Page = 1
	}

	if query.SortBy == "" {
		query.SortBy = "created_at"
	}

	if query.Order == "" {
		query.Order = "desc"
	}

	users, err := s.userRepo.GetUsers(ctx, query)
	if err != nil {
		return dto.GetUsersResponse{}, err
	}

	return dto.GetUsersResponse{
		Users: users,
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, req dto.GetUserByIDRequest) (dto.GetUserByIDResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return dto.GetUserByIDResponse{}, valErr
	}

	user, err := s.userRepo.GetUserByField(ctx, "id", req.ID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.GetUserByIDResponse{}, domain.ErrUserNotFound
		}

		return dto.GetUserByIDResponse{}, err
	}

	return dto.GetUserByIDResponse{
		User: *user,
	}, nil
}

func (s *userService) GetUsersStats(ctx context.Context) (dto.GetUsersStatsResponse, error) {
	resultCh := make(chan int64, 2)
	errCh := make(chan error, 2)

	go func() {
		totalUsers, err := s.userRepo.CountUsers(ctx, dto.GetUsersStatsQuery{
			IncludeDeleted: true,
		})
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- totalUsers
	}()

	go func() {
		totalNonDeletedUsers, err := s.userRepo.CountUsers(ctx, dto.GetUsersStatsQuery{
			IncludeDeleted: false,
		})
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- totalNonDeletedUsers
	}()

	var totalUsers, totalNonDeletedUsers int64
	for i := 0; i < 2; i++ {
		select {
		case res := <-resultCh:
			if totalUsers == 0 {
				totalUsers = res
			} else {
				totalNonDeletedUsers = res
			}
		case err := <-errCh:
			return dto.GetUsersStatsResponse{}, err
		}
	}

	totalDeletedUsers := totalUsers - totalNonDeletedUsers

	return dto.GetUsersStatsResponse{
		TotalUsers:           totalUsers,
		TotalNonDeletedUsers: totalNonDeletedUsers,
		TotalDeletedUsers:    totalDeletedUsers,
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return dto.CreateUserResponse{}, valErr
	}

	uuid, err := s.uuid.NewV7()
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	hashedPassword, err := s.bcrypt.Hash(req.Password)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	user := &entity.User{
		ID:       uuid,
		Name:     req.Name,
		Password: hashedPassword,
		Email:    req.Email,
	}

	_, err = s.userRepo.GetUserByField(ctx, "email", user.Email)
	if err == nil { // successfully found a user with the same email
		return dto.CreateUserResponse{}, domain.ErrUserEmailAlreadyExists
	}

	if !errors.Is(err, sql.ErrNoRows) { // some other error occurred
		return dto.CreateUserResponse{}, err
	}

	id, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	return dto.CreateUserResponse{
		ID: id,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, req dto.UpdateUserRequest) (dto.UpdateUserResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return dto.UpdateUserResponse{}, valErr
	}

	user := &entity.User{
		ID:       req.ID,
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
	}

	_, err := s.userRepo.GetUserByField(ctx, "id", user.ID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.UpdateUserResponse{}, domain.ErrUserNotFound
		}

		return dto.UpdateUserResponse{}, err
	}

	_, err = s.userRepo.GetUserByField(ctx, "email", user.Email)
	if err == nil { // successfully found a user with the same email
		return dto.UpdateUserResponse{}, domain.ErrUserEmailAlreadyExists
	}

	if !errors.Is(err, sql.ErrNoRows) { // some other error occurred
		return dto.UpdateUserResponse{}, err
	}

	hashedPassword, err := s.bcrypt.Hash(req.Password)
	if err != nil {
		return dto.UpdateUserResponse{}, err
	}

	user.Password = hashedPassword

	id, err := s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return dto.UpdateUserResponse{}, err
	}

	return dto.UpdateUserResponse{
		ID: id,
	}, nil
}

func (s *userService) SoftDeleteUser(ctx context.Context, req dto.SoftDeleteUserRequest) (dto.SoftDeleteUserResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return dto.SoftDeleteUserResponse{}, valErr
	}

	id, err := s.userRepo.SoftDeleteUser(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.SoftDeleteUserResponse{}, domain.ErrUserNotFound
		}

		return dto.SoftDeleteUserResponse{}, err
	}

	return dto.SoftDeleteUserResponse{
		ID: id,
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, req dto.DeleteUserRequest) (dto.DeleteUserResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return dto.DeleteUserResponse{}, valErr
	}

	id, err := s.userRepo.DeleteUser(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.DeleteUserResponse{}, domain.ErrUserNotFound
		}

		return dto.DeleteUserResponse{}, err
	}

	return dto.DeleteUserResponse{
		ID: id,
	}, nil
}

func (s *userService) RestoreUser(ctx context.Context, req dto.RestoreUserRequest) (dto.RestoreUserResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return dto.RestoreUserResponse{}, valErr
	}

	id, err := s.userRepo.RestoreUser(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.RestoreUserResponse{}, domain.ErrUserNotFound
		}

		return dto.RestoreUserResponse{}, err
	}

	return dto.RestoreUserResponse{
		ID: id,
	}, nil
}
