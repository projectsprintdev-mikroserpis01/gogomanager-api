package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type departmentRepository struct {
	DB *sqlx.DB
}

const (
	queryCreate                 = "INSERT INTO departments (name, manager_id, created_at) VALUES($1, $2, $3) RETURNING id"
	queryDelete                 = "DELETE FROM departments WHERE id = $1"
	queryFindAllWithLimitOffset = "SELECT * FROM departments LIMIT $1 OFFSET $2"
	queryFindAll                = "SELECT * FROM departments"
	queryFindByName             = "SELECT * FROM departments WHERE name ILIKE $1 ORDER BY name ASC LIMIT $2 OFFSET $3"
	queryUpdate                 = "UPDATE departments SET name = $1 WHERE id = $2"
)

func NewDepartmentRepository(db *sqlx.DB) contracts.DepartmentRepository {
	return &departmentRepository{DB: db}
}

func (repo *departmentRepository) Create(ctx context.Context, data entity.Department) (int, error) {
	var id int

	err := repo.DB.QueryRowContext(ctx, queryCreate, data.Name, data.ManagerID, data.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *departmentRepository) Delete(ctx context.Context, id int) error {
	_, err := repo.DB.ExecContext(ctx, queryDelete, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *departmentRepository) FindAll(ctx context.Context) ([]*entity.Department, error) {
	var listDepartment []*entity.Department

	err := repo.DB.SelectContext(ctx, &listDepartment, queryFindAll)
	if err != nil {
		return nil, err
	}

	return listDepartment, nil
}

func (repo *departmentRepository) FindByName(ctx context.Context, name string, limit, offset int) ([]*entity.Department, error) {

	var listDepartment []*entity.Department
	searchTerm := "%" + name + "%"
	err := repo.DB.SelectContext(ctx, &listDepartment, queryFindByName, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}

	return listDepartment, nil
}

func (repo *departmentRepository) Update(ctx context.Context, id int, newName string) (int, error) {
	result, err := repo.DB.ExecContext(ctx, queryUpdate, newName, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()

	return int(rowsAffected), nil
}

func (repo *departmentRepository) FindAllWithLimitOffset(ctx context.Context, limit int, offset int) ([]*entity.Department, error) {
	var listDepartment []*entity.Department

	err := repo.DB.SelectContext(ctx, &listDepartment, queryFindAllWithLimitOffset, limit, offset)
	if err != nil {
		return nil, err
	}

	return listDepartment, nil
}

func (repo *departmentRepository) FindByID(ctx context.Context, id int) (*entity.Department, error) {
	var department entity.Department

	err := repo.DB.GetContext(ctx, &department, "SELECT * FROM departments WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &department, nil
}
