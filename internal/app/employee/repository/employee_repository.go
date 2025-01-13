package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type employeeRepository struct {
	DB *sqlx.DB
}

const (
	queryCreate = `
	INSERT INTO employees (identity_number, name, employee_image_uri, gender, department_id) 
	VALUES (:identity_number, :name, :employee_image_uri, :gender, :department_id) RETURNING id`
	queryFindBase             = "SELECT * FROM employees WHERE 1=1"
	queryFindByIdentityNumber = "SELECT * FROM employees WHERE identity_number=?"
	queryDelete               = "DELETE FROM employees WHERE id = ?"
	queryUpdate               = `
		UPDATE employees 
			SET name = ?, 
 			identity_number = ?, 
    		gender = ?, 
    		department_id = ?, 
    		employee_image_uri = ? 
		WHERE id = ?`
)

func NewEmployeeRepository(db *sqlx.DB) contracts.EmployeeRepository {
	return &employeeRepository{DB: db}
}

func (e *employeeRepository) Create(ctx context.Context, data entity.Employee) error {
	_, err := e.DB.ExecContext(
		ctx,
		queryCreate,
		data.IdentityNumber,
		data.Name,
		data.EmployeeImageURI,
		data.Gender,
		data.DepartmentID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (e *employeeRepository) FindByIdentityNumber(
	ctx context.Context,
	identityNumber string,
) (*entity.Employee, error) {
	employee := &entity.Employee{}

	err := e.DB.SelectContext(ctx, &employee, queryFindByIdentityNumber, identityNumber)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (e *employeeRepository) Find(
	ctx context.Context,
	identityNumber string,
	name string,
	gender string,
	departmentID int,
	limit int,
	offset int,
) ([]*entity.Employee, error) {
	employees := []*entity.Employee{}

	query := queryFindBase

	args := []interface{}{}
	if identityNumber != "" {
		query += " AND identity_number = ?"
		args = append(args, identityNumber)
	}
	if name != "" {
		query += " AND name ILIKE ?"
		args = append(args, "%"+name+"%")
	}
	if gender != "" {
		query += " AND gender = ?"
		args = append(args, gender)
	}
	if departmentID != 0 {
		query += " AND department_id = ?"
		args = append(args, departmentID)
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err := e.DB.SelectContext(ctx, &employees, query, args...)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (e *employeeRepository) Update(ctx context.Context, data entity.Employee) error {

	_, err := e.DB.ExecContext(ctx, queryUpdate,
		data.Name,
		data.IdentityNumber,
		data.Gender,
		data.DepartmentID,
		data.EmployeeImageURI,
		data.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

func (e *employeeRepository) Delete(ctx context.Context, id int) error {
	_, err := e.DB.ExecContext(ctx, queryDelete, id)
	return err
}
