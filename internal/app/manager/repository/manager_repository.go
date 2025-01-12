package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type ManagerRepository interface {
	EmailExists(ctx context.Context, email string) (bool, error)
	CreateManager(ctx context.Context, req dto.AuthRequest) error
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

func (r *managerRepository) CreateManager(ctx context.Context, req dto.AuthRequest) error {
	_, err := r.db.Exec("INSERT INTO managers (email, password) VALUES ($1, $2, $3)", req.Email, req.Password)
	return err
}

func (r *managerRepository) GetManagerByEmail(ctx context.Context, email string) (entity.Manager, error) {
	var manager entity.Manager
	err := r.db.Get(&manager, "SELECT * FROM managers WHERE email=$1", email)
	return manager, err
}
