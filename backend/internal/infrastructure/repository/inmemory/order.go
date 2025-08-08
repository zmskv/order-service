package inmemory

import (
	"context"
	"errors"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/zmskv/order-service/internal/domain/order/entity"
	"go.uber.org/zap"
)

type InMemoryRepository struct {
	cache  *lru.Cache[string, entity.Order]
	logger *zap.Logger
}

func NewInMemoryRepository(log *zap.Logger) (*InMemoryRepository, error) {
	cache, err := lru.New[string, entity.Order](1000)
	if err != nil {
		return nil, err
	}
	return &InMemoryRepository{
		cache:  cache,
		logger: log.Named("in_memory_repo"),
	}, nil
}

func (r *InMemoryRepository) Save(ctx context.Context, order entity.Order) error {
	r.cache.Add(order.OrderUID, order)
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, uid string) (entity.Order, error) {
	order, ok := r.cache.Get(uid)
	if !ok {
		r.logger.Info("Order not found", zap.String("order_uid", uid))
		return entity.Order{}, errors.New("order not found")
	}
	return order, nil
}

func (r *InMemoryRepository) GetAll(ctx context.Context) ([]entity.Order, error) {
	var orders []entity.Order
	for _, key := range r.cache.Keys() {
		order, ok := r.cache.Get(key)
		if !ok {
			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}
