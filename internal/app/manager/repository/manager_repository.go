package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type ManagerRepository interface {
	EmailExists(ctx context.Context, email string) (bool, error)
	CreateManager(ctx context.Context, req dto.AuthRequest) (entity.Manager, error)
	GetManagerByEmail(ctx context.Context, email string) (entity.Manager, error)
}

type managerRepository struct {
	db *sqlx.DB
}

func NewManagerRepository(db *sqlx.DB) ManagerRepository {
	return &managerRepository{
		db: db,
	}
}

func (r *managerRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM managers WHERE email=$1)", email)
	return exists, err
}

func (r *managerRepository) CreateManager(ctx context.Context, req dto.AuthRequest) (entity.Manager, error) {
	_, err := r.db.Exec("INSERT INTO managers (email, password) VALUES ($1, $2)", req.Email, req.Password)
	if err != nil {
		return entity.Manager{}, err
	}

	return r.GetManagerByEmail(ctx, req.Email)
}

func (r *managerRepository) UpdateManager(ctx context.Context, req dto.ManagerProfile) (entity.Manager, error) {
	_, err := r.db.Exec("UPDATE managers SET name = $1, user_image_uri = $2, company_name = $3, company_image_uri = $4 WHERE email = $5", req.Name, req.UserImageUri, req.CompanyName, req.CompanyImageUri, req.Email)
	if err != nil {
		return entity.Manager{}, err
	}

	return r.GetManagerByEmail(ctx, req.Email)
}

func (r *managerRepository) GetManagerByEmail(ctx context.Context, email string) (entity.Manager, error) {
	var manager entity.Manager
	err := r.db.Get(&manager, "SELECT * FROM managers WHERE email=$1", email)
	return manager, err
}
