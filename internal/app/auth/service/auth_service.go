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
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/jwt"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/uuid"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/validator"
)

type authService struct {
	authRepo  contracts.AuthRepository
	validator validator.ValidatorInterface
	uuid      uuid.UUIDInterface
	jwt       jwt.JwtInterface
	bcrypt    bcrypt.BcryptInterface
}

func NewAuthService(
	authRepo contracts.AuthRepository,
	validator validator.ValidatorInterface,
	uuid uuid.UUIDInterface,
	jwt jwt.JwtInterface,
	bcrypt bcrypt.BcryptInterface,
) contracts.AuthService {
	return &authService{
		authRepo:  authRepo,
		validator: validator,
		uuid:      uuid,
		jwt:       jwt,
		bcrypt:    bcrypt,
	}
}

func (s *authService) RegisterUser(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return dto.RegisterResponse{}, valErr
	}

	_, err := s.authRepo.GetUserByEmail(ctx, req.Email)
	if err == nil { // successfully found a user with the same email
		return dto.RegisterResponse{}, nil
	}

	if !errors.Is(err, sql.ErrNoRows) { // some other error occurred
		return dto.RegisterResponse{}, err
	}

	uuid, err := s.uuid.NewV7()
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	hashedPassword, err := s.bcrypt.Hash(req.Password)
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	user := entity.User{
		ID:       uuid,
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	id, err := s.authRepo.RegisterUser(ctx, user)
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	res := dto.RegisterResponse{
		ID: id,
	}

	return res, nil
}

func (s *authService) LoginUser(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return dto.LoginResponse{}, valErr
	}

	user, err := s.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.LoginResponse{}, domain.ErrEmailNotFound
		}

		return dto.LoginResponse{}, err
	}

	isValid := s.bcrypt.Compare(req.Password, user.Password)
	if !isValid {
		return dto.LoginResponse{}, domain.ErrCredentialsNotMatch
	}

	accessToken, err := s.jwt.Create(user.ID, user.Role.Name)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	res := dto.LoginResponse{
		AccessToken: accessToken,
	}

	return res, nil
}
