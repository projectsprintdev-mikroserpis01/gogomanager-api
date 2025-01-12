package service

import (
	"context"
	"errors"

	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/manager/repository"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/bcrypt"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/jwt"
)

type ManagerService interface {
	Authenticate(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error)
}

type managerService struct {
	repo   repository.ManagerRepository
	jwt    jwt.JwtManagerInterface
	bcrypt bcrypt.BcryptInterface
}

func NewManagerService(repo repository.ManagerRepository, jwt jwt.JwtManagerInterface, bcrypt bcrypt.BcryptInterface) ManagerService {
	return &managerService{
		repo:   repo,
		jwt:    jwt,
		bcrypt: bcrypt,
	}
}

func (s *managerService) Authenticate(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error) {
	switch req.Action {
	case "create":
		exists, err := s.repo.EmailExists(ctx, req.Email)
		if err != nil {
			return dto.AuthResponse{}, err
		}

		if exists {
			return dto.AuthResponse{}, errors.New("email already exists")
		}

		hashedPassword, err := s.bcrypt.Hash(req.Password)
		if err != nil {
			return dto.AuthResponse{}, err
		}

		req.Password = hashedPassword
		err = s.repo.CreateManager(ctx, req)
		if err != nil {
			return dto.AuthResponse{}, err
		}

		token, err := s.jwt.CreateManager(req.Email)
		if err != nil {
			return dto.AuthResponse{}, err
		}

		return dto.AuthResponse{Email: req.Email, Token: token}, nil

	case "login":
		manager, err := s.repo.GetManagerByEmail(ctx, req.Email)
		if err != nil {
			return dto.AuthResponse{}, errors.New("email not found")
		}

		isValid := s.bcrypt.Compare(req.Password, manager.Password)
		if !isValid {
			return dto.AuthResponse{}, domain.ErrCredentialsNotMatch
		}

		token, err := s.jwt.CreateManager(req.Email)
		if err != nil {
			return dto.AuthResponse{}, err
		}

		return dto.AuthResponse{Email: req.Email, Token: token}, nil

	default:
		return dto.AuthResponse{}, errors.New("invalid action")
	}
}
