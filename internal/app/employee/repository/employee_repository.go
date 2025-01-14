package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/log"
)

type employeeRepository struct {
	DB *sqlx.DB
}

const (
	queryCreate = `
	INSERT INTO employees (identity_number, name, employee_image_uri, gender, department_id)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`
	queryFindBase             = "SELECT * FROM employees WHERE 1=1"
	queryFindByIdentityNumber = "SELECT * FROM employees WHERE identity_number = $1"
	queryDelete               = "DELETE FROM employees WHERE id = $1"
	queryUpdate               = `
		UPDATE employees
			SET name = $1,
 			identity_number = $2,
    		gender = $3,
    		department_id = $4,
    		employee_image_uri = $5
		WHERE id = $6`
)

func NewEmployeeRepository(db *sqlx.DB) contracts.EmployeeRepository {
	return &employeeRepository{DB: db}
}

func (e *employeeRepository) Create(ctx context.Context, data entity.Employee) error {
	log.Info(log.LogInfo{
		"identityNumber": data.IdentityNumber,
	}, "[EmployeeRepository.Create]")

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
	var employee entity.Employee

	err := e.DB.GetContext(ctx, &employee, queryFindByIdentityNumber, identityNumber)
	if err != nil {
		return nil, err
	}

	return &employee, nil
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

	args := map[string]interface{}{}
	if identityNumber != "" {
		query += " AND identity_number = :identity_number"
		args["identity_number"] = identityNumber
	}
	if name != "" {
		query += " AND name ILIKE :name"
		args["name"] = "%" + name + "%"
	}
	if gender != "" {
		query += " AND gender = :gender"
		args["gender"] = gender
	}
	if departmentID != 0 {
		query += " AND department_id = :department_id"
		args["department_id"] = departmentID
	}

	query += " LIMIT :limit OFFSET :offset"
	args["limit"] = limit
	args["offset"] = offset

	finalQuery, finalArgs, err := sqlx.Named(query, args)
	if err != nil {
		return nil, err
	}

	finalQuery = e.DB.Rebind(finalQuery)

	err = e.DB.SelectContext(ctx, &employees, finalQuery, finalArgs...)
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
