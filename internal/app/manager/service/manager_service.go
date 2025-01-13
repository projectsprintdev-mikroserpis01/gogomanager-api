package service

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/manager/repository"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/bcrypt"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/jwt"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/log"
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

	ret := dto.GetCurrentManagerResponse{Email: manager.Email, Name: manager.Name, UserImageUri: manager.UserImageURI, CompanyName: manager.CompanyName, CompanyImageUri: manager.CompanyImageURI}
	return &ret, nil
}

func (s *managerService) UpdateManagerById(ctx context.Context, id int, req dto.UpdateManagerRequest) (*dto.UpdateManagerResponse, error) {
	valErr := s.validator.Validate(req)
	if valErr != nil {
		return nil, valErr
	}

	if req.Email != nil {
		manager, err := s.repo.GetManagerByEmail(ctx, *req.Email)
		if err == nil { // found a manager with the same email
			if manager.ID != id {
				log.Info(log.LogInfo{
					"manager id": manager.ID,
					"id":         id,
				}, "[managerService.UpdateManagerById] id")

				return nil, domain.ErrUserEmailAlreadyExists
			}
		}

		if err != nil && !errors.Is(err, sql.ErrNoRows) { // some other error occurred
			return nil, err
		}
	}

	log.Info(log.LogInfo{
		"is req user image uri nil": req.UserImageUri == nil,
	}, "[managerService.UpdateManagerById] id")

	if req.UserImageUri != nil {
		u, err := url.ParseRequestURI(*req.UserImageUri)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "invalid user image uri")
		}

		if u.Scheme == "" || u.Host == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "invalid user image uri")
		}

		// Additional validation: Check if the host contains a domain or is not empty
		if !strings.Contains(u.Host, ".") {
			return nil, fiber.NewError(fiber.StatusBadRequest, "invalid company image uri")
		}
	}

	if req.CompanyImageUri != nil {
		u, err := url.ParseRequestURI(*req.CompanyImageUri)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "invalid company image uri")
		}

		log.Info(log.LogInfo{
			"url": u,
		}, "[managerService.UpdateManagerById] company image uri")

		if u.Scheme == "" || u.Host == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "invalid company image uri")
		}

		// Additional validation: Check if the host contains a domain or is not empty
		if !strings.Contains(u.Host, ".") {
			return nil, fiber.NewError(fiber.StatusBadRequest, "invalid company image uri")
		}
	}

	fields := []string{}
	args := []interface{}{}
	if req.Email != nil {
		fields = append(fields, "email")
		args = append(args, *req.Email)
	}
	if req.Name != nil {
		fields = append(fields, "name")
		args = append(args, *req.Name)
	}
	if req.UserImageUri != nil {
		fields = append(fields, "user_image_uri")
		args = append(args, *req.UserImageUri)
	}
	if req.CompanyName != nil {
		fields = append(fields, "company_name")
		args = append(args, *req.CompanyName)
	}
	if req.CompanyImageUri != nil {
		fields = append(fields, "company_image_uri")
		args = append(args, *req.CompanyImageUri)
	}

	_, err := s.repo.UpdateManagerByIDSomeFields(ctx, id, fields, args)
	if err != nil {
		return nil, err
	}

	ret := dto.UpdateManagerResponse{}

	return &ret, nil
}
