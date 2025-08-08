package interfaces

import (
	"context"

	"github.com/zmskv/order-service/internal/domain/order/entity"
)

type OrderService interface {
	GetByID(ctx context.Context, id string) (entity.Order, error)
}
