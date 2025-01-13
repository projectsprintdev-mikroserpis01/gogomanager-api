package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type ManagerRepository interface {
	EmailExists(ctx context.Context, email string) (bool, error)
	CreateManager(ctx context.Context, req dto.AuthRequest) (entity.Manager, error)
	GetManagerByEmail(ctx context.Context, email string) (entity.Manager, error)
	GetManagerById(ctx context.Context, id int) (*entity.Manager, error)
	UpdateManagerById(ctx context.Context, id int, email string, name string, userImageUri string, companyName string, companyImageUri string) (int, error)
	UpdateManagerByIDSomeFields(ctx context.Context, id int, fields []string, args []interface{}) (int, error)
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

func (r *managerRepository) GetManagerById(ctx context.Context, id int) (*entity.Manager, error) {
	var manager entity.Manager
	err := r.db.Get(&manager, "SELECT * FROM managers WHERE id=$1", id)
	return &manager, err
}

func (r *managerRepository) UpdateManagerById(ctx context.Context, id int, email string, name string, userImageUri string, companyName string, companyImageUri string) (int, error) {

	_, err := r.GetManagerByEmail(ctx, email)
	if err == nil { // successfully found a manager with the same email
		return 0, domain.ErrUserEmailAlreadyExists
	}
	result, err := r.db.ExecContext(ctx, "UPDATE managers SET name = $1, user_image_uri = $2, company_name = $3, company_image_uri = $4 WHERE id = $5", name, userImageUri, companyName, companyImageUri, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), nil
}

func (r *managerRepository) UpdateManagerByIDSomeFields(ctx context.Context, id int, fields []string, args []interface{}) (int, error) {
	query := `UPDATE managers SET `
	for i, field := range fields {
		query += fmt.Sprintf("%s = $%d", field, i+1)
		if i != len(fields)-1 {
			query += ", "
		}
	}
	query += fmt.Sprintf(" WHERE id = $%d", len(fields)+1)

	args = append(args, id)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), nil
}
