package interfaces

import (
	"context"

	"github.com/zmskv/order-service/internal/domain/order/entity"
)

type OrderRepository interface {
	Save(ctx context.Context, order entity.Order) error
	Get(ctx context.Context, orderUID string) (entity.Order, error)
	GetAll(ctx context.Context) ([]entity.Order, error)
}
