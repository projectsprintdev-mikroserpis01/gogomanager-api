package contracts

import (
	"context"

	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type ManagerRepository interface {
	GetManagerById(ctx context.Context, id int) (*entity.Manager, error)
	UpdateManagerById(ctx context.Context, id int, req dto.UpdateManagerRequest) (int, error)
}

type ManagerService interface {
	GetManagerById(ctx context.Context, id int) (*dto.GetCurrentManagerResponse, error)
	UpdateManagerById(ctx context.Context, id int, req dto.UpdateManagerRequest) (*dto.UpdateManagerResponse, error)
}
