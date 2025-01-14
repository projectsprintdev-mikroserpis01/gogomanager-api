package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/log"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/validator"
)

type employeeService struct {
	repo      contracts.EmployeeRepository
	validator validator.ValidatorInterface
}

func NewEmployeeService(
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

	employeeImageUri, err := url.ParseRequestURI(data.EmployeeImageURI)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid employee image uri")
	}

	if employeeImageUri.Scheme == "" || employeeImageUri.Host == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid employee image uri")
	}

	// Additional validation: Check if the host contains a domain or is not empty
	if !strings.Contains(employeeImageUri.Host, ".") {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid employee image uri")
	}

	strDepartmentID, _ := strconv.Atoi(data.DepartmentID)
	employee := entity.Employee{
		IdentityNumber:   data.IdentityNumber,
		Name:             data.Name,
		Gender:           data.Gender,
		DepartmentID:     strDepartmentID,
		EmployeeImageURI: data.EmployeeImageURI,
	}

	err = e.repo.Create(ctx, employee)
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

	log.Info(log.LogInfo{
		"employeeDataRes": employeeDataRes,
	}, "[EmployeeService.Create]")

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

	log.Info(log.LogInfo{
		"listResponseData": listResponseData,
	}, "[EmployeeService.Find]")

	return listResponseData, nil
}

func (e employeeService) Update(
	ctx context.Context,
	data dto.EmployeeUpdateReq,
	identityNumber string,
) (*dto.EmployeeDataRes, error) {
	valErr := e.validator.Validate(&data)
	if valErr != nil {
		return nil, valErr
	}

	employeeImageUri, err := url.ParseRequestURI(data.EmployeeImageURI)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid employee image uri")
	}

	if employeeImageUri.Scheme == "" || employeeImageUri.Host == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid employee image uri")
	}

	// Additional validation: Check if the host contains a domain or is not empty
	if !strings.Contains(employeeImageUri.Host, ".") {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid employee image uri")
	}

	oldData, err := e.repo.FindByIdentityNumber(ctx, identityNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("employee with id %s not found", identityNumber))
		}
		return nil, fmt.Errorf("failed to update employee: %w", err)
	}

	updatedData := generateUpdateData(data, *oldData)

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
) entity.Employee {
	updatedData := oldData

	log.Error(log.LogInfo{
		"newData": newData,
		"oldData": oldData,
	}, "[generateUpdateData]")

	if newData.IdentityNumber != "" {
		updatedData.IdentityNumber = newData.IdentityNumber
	}

	if newData.Name != "" {
		updatedData.Name = newData.Name
	}

	if newData.EmployeeImageURI != "" {
		updatedData.EmployeeImageURI = newData.EmployeeImageURI
	}

	if newData.DepartmentID != "" {
		log.Info(log.LogInfo{
			"newData.DepartmentID": newData.DepartmentID,
			"is not empty":         newData.DepartmentID != "",
			"is not nil":           &newData.DepartmentID != nil,
			"oldData.DepartmentID": oldData.DepartmentID,
		}, "newData.DepartmentID")
		updatedData.DepartmentID, _ = strconv.Atoi(newData.DepartmentID)
	}

	if newData.Gender != "" {
		updatedData.Gender = newData.Gender
	}

	log.Info(log.LogInfo{
		"updatedData": updatedData,
	}, "updatedData")

	return updatedData
}
