package contracts

import (
	"context"

	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type DepartmentRepository interface {
	Create(ctx context.Context, data entity.Department) (int, error)
	FindByName(ctx context.Context, name string) ([]*entity.Department, error)
	FindAll(ctx context.Context) ([]*entity.Department, error)
	FindAllWithLimitOffset(ctx context.Context, limit, offset int) ([]*entity.Department, error)
	Update(ctx context.Context, id int, newName string) (int, error)
	Delete(ctx context.Context, id int) error
}

type DepartmentService interface {
	Create(ctx context.Context, managerId int, name string) (*dto.DepartmentRes, error)
	Update(ctx context.Context, id int, name string) (*dto.DepartmentRes, error)
	FindAll(ctx context.Context, limit, offset int) ([]*dto.DepartmentRes, error)
	FindByName(ctx context.Context, limit, offset int, name string) ([]*dto.DepartmentRes, error)
	Delete(ctx context.Context, id int) error
}
