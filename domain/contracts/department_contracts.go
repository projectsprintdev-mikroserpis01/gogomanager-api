package contracts

import ("github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
	"context"
)

type DepartmentRepository interface {
	Create(ctx context.Context data entity.Department) (int, error)
	FindByName(ctx context.Context, limit, offset int, name string) ([]*entity.Department, error)
	FindAll(ctx context.Context, limit, offset int) ([]*entity.Department, error)
	Update(ctx context.Context, newName string) (float64, error)
	Delete(ctx context.Context, Id int) error
}

type DepartmentService interface {
	Create(ctx context.Context, name string) (entity.DataRes, error)
	Update(ctx context.Context, name string) (entity.DataRes, error)
	FindAll(ctx context.Context, limit, offset int) (entity.DataRes, error)
	FindByName(ctx context.Context, limit, offset int, name string) (entity.DataRes, error)
	Delete(ctx context.Context, id int) error
}
