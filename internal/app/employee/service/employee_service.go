package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/validator"
)

type employeeService struct {
	repo      contracts.EmployeeRepository
	validator validator.ValidatorInterface
}

func NewemployeeService(
	repository contracts.EmployeeRepository,
	validator validator.ValidatorInterface,
) contracts.EmployeeService {
	return employeeService{repo: repository, validator: validator}
}

func (e employeeService) Create(
	ctx context.Context,
	data dto.EmployeeCreateReq,
) (*dto.EmployeeDataRes, error) {
	valErr := e.validator.Validate(&data)
	if valErr != nil {
		return nil, valErr
	}

	strDepartmentID, _ := strconv.Atoi(data.DepartmentID)
	employee := entity.Employee{
		IdentityNumber:   data.IdentityNumber,
		Name:             data.Name,
		Gender:           data.Gender,
		DepartmentID:     strDepartmentID,
		EmployeeImageURI: data.EmployeeImageURI,
	}
	err := e.repo.Create(ctx, employee)
	if err != nil {
		return nil, err
	}

	employeeDataRes := dto.EmployeeDataRes{
		IdentityNumber:   data.IdentityNumber,
		Name:             data.Name,
		Gender:           data.Gender,
		DepartmentID:     data.DepartmentID,
		EmployeeImageURI: data.EmployeeImageURI,
	}

	return &employeeDataRes, nil
}

func (e employeeService) Delete(
	ctx context.Context,
	identityNumber string,
) error {

	identityNumberInt, err := strconv.Atoi(identityNumber)
	if err != nil {
		return fiber.NewError(400, fiber.ErrBadRequest.Message)
	}

	err = e.repo.Delete(ctx, identityNumberInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("employee with id %s not found", identityNumber))
		}
		return fmt.Errorf("failed to delete employee: %w", err)
	}
	return nil
}

func (e employeeService) Find(
	ctx context.Context,
	identityNumber string,
	name string,
	gender string,
	departmentID int,
	limit int,
	offset int,
) ([]*dto.EmployeeDataRes, error) {

	listData, err := e.repo.Find(
		ctx,
		identityNumber,
		name,
		gender,
		departmentID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	listResponseData := []*dto.EmployeeDataRes{}
	for _, data := range listData {
		responseData := dto.EmployeeDataRes{
			IdentityNumber:   data.IdentityNumber,
			Name:             data.Name,
			EmployeeImageURI: data.EmployeeImageURI,
			Gender:           data.Gender,
			DepartmentID:     fmt.Sprintf("%d", data.DepartmentID),
		}
		listResponseData = append(listResponseData, &responseData)
	}

	return listResponseData, nil
}

func (e employeeService) Update(
	ctx context.Context,
	data dto.EmployeeUpdateReq,
	identityNumber string,
) (*dto.EmployeeDataRes, error) {

	oldData, err := e.repo.FindByIdentityNumber(ctx, data.IdentityNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("employee with id %s not found", data.IdentityNumber))
		}
		return nil, fmt.Errorf("failed to delete employee: %w", err)
	}

	updatedData := generateUpdateData(data, *oldData, identityNumber)

	err = e.repo.Update(ctx, updatedData)
	if err != nil {
		return nil, err
	}
	res := &dto.EmployeeDataRes{
		IdentityNumber:   updatedData.IdentityNumber,
		Name:             updatedData.Name,
		EmployeeImageURI: updatedData.EmployeeImageURI,
		Gender:           updatedData.Gender,
		DepartmentID:     fmt.Sprintf("%d", updatedData.DepartmentID),
	}
	return res, nil
}

func generateUpdateData(
	newData dto.EmployeeUpdateReq,
	oldData entity.Employee,
	oldIdentityNumber string,
) entity.Employee {
	updatedData := entity.Employee{}

	return updatedData
}
