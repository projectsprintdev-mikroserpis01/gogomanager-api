package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/validator"
)

type departmentService struct {
	repo      contracts.DepartmentRepository
	validator validator.ValidatorInterface
}

func NewDepartmentService(repository contracts.DepartmentRepository, validator validator.ValidatorInterface) contracts.DepartmentService {
	return departmentService{repo: repository, validator: validator}
}

func (d departmentService) Create(ctx context.Context, managerId int, name string) (*dto.DepartmentRes, error) {
	type Request struct {
		Name string `json:"name" validate:"required,min=4,max=33"`
	}

	req := Request{Name: name}

	valErr := d.validator.Validate(&req)
	if valErr != nil {
		return nil, valErr
	}

	department := entity.Department{
		Name:      name,
		ManagerID: managerId,
		CreatedAt: time.Now(),
	}

	id, err := d.repo.Create(ctx, department)
	if err != nil {
		return nil, err
	}

	return &dto.DepartmentRes{
		ID:   strconv.Itoa(id),
		Name: department.Name,
	}, nil
}

func (d departmentService) Delete(ctx context.Context, id int) error {
	_, err := d.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("department with id %d not found", id))
		}

		return err
	}

	err = d.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (d departmentService) FindAll(ctx context.Context, limit int, offset int) ([]*dto.DepartmentRes, error) {
	var departments []*entity.Department
	var err error

	departments, err = d.repo.FindAllWithLimitOffset(ctx, limit, offset)

	if err != nil {
		return nil, err
	}

	var departmentRes []*dto.DepartmentRes
	for _, dept := range departments {
		departmentRes = append(departmentRes, &dto.DepartmentRes{
			ID:   strconv.Itoa(dept.ID),
			Name: dept.Name,
		})
	}

	return departmentRes, nil
}

func (d departmentService) FindByName(ctx context.Context, limit int, offset int, name string) ([]*dto.DepartmentRes, error) {

	departments, err := d.repo.FindByName(ctx, name, limit, offset)
	if err != nil {
		return nil, err
	}

	var departmentRes []*dto.DepartmentRes
	for _, dept := range departments {
		departmentRes = append(departmentRes, &dto.DepartmentRes{
			ID:   strconv.Itoa(dept.ID),
			Name: dept.Name,
		})
	}

	return departmentRes, nil
}

func (d departmentService) Update(ctx context.Context, id int, name string) (*dto.DepartmentRes, error) {
	type Request struct {
		Name string `json:"name" validate:"required,min=4,max=33"`
	}

	req := Request{Name: name}

	valErr := d.validator.Validate(&req)
	if valErr != nil {
		return nil, valErr
	}

	_, err := d.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("department with id %d not found", id))
		}

		return nil, err
	}

	_, err = d.repo.Update(ctx, id, name)
	if err != nil {
		return nil, err
	}

	return &dto.DepartmentRes{
		ID:   strconv.Itoa(id),
		Name: name,
	}, nil
}
