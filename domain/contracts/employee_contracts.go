package contracts

import (
	"context"

	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type EmployeeRepository interface {
	Create(ctx context.Context, data entity.Employee) error
	Find(ctx context.Context, identityNumber, name, gender string, departmentID, limit, offset int) ([]*entity.Employee, error)
	FindByIdentityNumber(ctx context.Context, identityNumber string) (*entity.Employee, error)
	Update(ctx context.Context, data entity.Employee) error
	Delete(ctx context.Context, id int) error
}

type EmployeeService interface {
	Create(ctx context.Context, data dto.EmployeeCreateReq) (*dto.EmployeeDataRes, error)
	Update(ctx context.Context, data dto.EmployeeUpdateReq, identityNumber string) (*dto.EmployeeDataRes, error)
	Find(ctx context.Context, identityNumber, name, gender string, departmentID, limit, offset int) ([]*dto.EmployeeDataRes, error)
	Delete(ctx context.Context, identityNumber string) error
}
