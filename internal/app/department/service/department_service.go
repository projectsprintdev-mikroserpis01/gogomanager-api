package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type departmentService struct {
	repo contracts.DepartmentRepository
}

func NewDepartmentService(repository contracts.DepartmentRepository) contracts.DepartmentService {
	return departmentService{repo: repository}
}

func (d departmentService) Create(ctx context.Context, managerId int, name string) (*dto.DepartmentRes, error) {
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
	err := d.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (d departmentService) FindAll(ctx context.Context, limit int, offset int) ([]*dto.DepartmentRes, error) {
	var departments []*entity.Department
	var err error

	if limit == 0 || offset == 0 {
		departments, err = d.repo.FindAll(ctx)
	} else {
		departments, err = d.repo.FindAllWithLimitOffset(ctx, limit, offset)
	}

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

	departments, err := d.repo.FindByName(ctx, name)
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

	rowsAffected, err := d.repo.Update(ctx, id, name)
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("department with ID %d not found", id)
	}

	return &dto.DepartmentRes{
		ID:   strconv.Itoa(id),
		Name: name,
	}, nil
}
