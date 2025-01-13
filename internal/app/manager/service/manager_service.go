package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/manager/repository"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/bcrypt"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/jwt"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/validator"
)

type ManagerService interface {
	Authenticate(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error)
	GetManagerById(ctx context.Context, id int) (*dto.GetCurrentManagerResponse, error)
	UpdateManagerById(ctx context.Context, id int, req dto.UpdateManagerRequest) (*dto.UpdateManagerResponse, error)
}

type managerService struct {
	repo      repository.ManagerRepository
	jwt       jwt.JwtManagerInterface
	bcrypt    bcrypt.BcryptInterface
	validator validator.ValidatorInterface
}

func NewManagerService(
	repo repository.ManagerRepository,
	jwt jwt.JwtManagerInterface,
	bcrypt bcrypt.BcryptInterface,
	validator validator.ValidatorInterface,
) ManagerService {
	return &managerService{
		repo:      repo,
		jwt:       jwt,
		bcrypt:    bcrypt,
		validator: validator,
	}
}

func (s *managerService) Authenticate(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return dto.AuthResponse{}, valErr
	}

	switch req.Action {
	case "create":
		exists, err := s.repo.EmailExists(ctx, req.Email)
		if err != nil {
			return dto.AuthResponse{}, err
		}

		if exists {
			return dto.AuthResponse{}, fiber.NewError(fiber.StatusConflict, "email already exists")
		}

		hashedPassword, err := s.bcrypt.Hash(req.Password)
		if err != nil {
			return dto.AuthResponse{}, err
		}

		req.Password = hashedPassword
		manager, err := s.repo.CreateManager(ctx, req)
		if err != nil {
			return dto.AuthResponse{}, err
		}

		token, err := s.jwt.CreateManager(manager.ID, manager.Email)
		if err != nil {
			return dto.AuthResponse{}, err
		}

		return dto.AuthResponse{Email: req.Email, Token: token}, nil

	case "login":
		manager, err := s.repo.GetManagerByEmail(ctx, req.Email)
		if err != nil {
			return dto.AuthResponse{}, fiber.NewError(fiber.StatusNotFound, "manager not found")
		}

		isValid := s.bcrypt.Compare(req.Password, manager.Password)
		if !isValid {
			return dto.AuthResponse{}, domain.ErrCredentialsNotMatch
		}

		token, err := s.jwt.CreateManager(manager.ID, req.Email)
		if err != nil {
			return dto.AuthResponse{}, err
		}

		return dto.AuthResponse{Email: req.Email, Token: token}, nil

	default:
		return dto.AuthResponse{}, errors.New("invalid action")
	}

}

func (s *managerService) GetManagerById(ctx context.Context, id int) (*dto.GetCurrentManagerResponse, error) {
	manager, err := s.repo.GetManagerById(ctx, id)
	if err != nil {
		return nil, err
	}

	ret := dto.GetCurrentManagerResponse{Email: manager.Email, Name: *manager.Name, UserImageUri: *manager.UserImageURI, CompanyName: *manager.CompanyName, CompanyImageUri: *manager.CompanyImageURI}
	return &ret, nil
}

func (s *managerService) UpdateManagerById(ctx context.Context, id int, req dto.UpdateManagerRequest) (*dto.UpdateManagerResponse, error) {
	rowsAffected, err := s.repo.UpdateManagerById(ctx, id, req.Email, req.Name, req.UserImageUri, req.CompanyName, req.CompanyImageUri)
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("manager with ID %d not found", id)
	}
	ret := dto.UpdateManagerResponse{Email: req.Email, Name: req.Name, UserImageUri: req.UserImageUri, CompanyName: req.CompanyName, CompanyImageUri: req.CompanyImageUri}
	return &ret, nil
}
