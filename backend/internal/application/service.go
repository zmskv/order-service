package application

import (
	"context"
	"encoding/json"

	"github.com/zmskv/order-service/internal/domain/order/entity"
	"github.com/zmskv/order-service/internal/domain/order/interfaces"
	"github.com/zmskv/order-service/internal/infrastructure/messaging/kafka"
	"github.com/zmskv/order-service/internal/infrastructure/messaging/kafka/dto"
	"go.uber.org/zap"
)

type OrderService struct {
	repo         interfaces.OrderRepository
	consumer     *kafka.Consumer
	inMemoryRepo interfaces.OrderRepository
	logger       *zap.Logger
}

func NewOrderService(repo interfaces.OrderRepository, consumer *kafka.Consumer, inMemoryRepo interfaces.OrderRepository, logger *zap.Logger) *OrderService {
	svc := &OrderService{
		repo:         repo,
		consumer:     consumer,
		inMemoryRepo: inMemoryRepo,
		logger:       logger.Named("service"),
	}
	svc.consumer.Handler = svc.handleMessage

	return svc
}

func (s *OrderService) Start(ctx context.Context) error {
	return s.consumer.Start()
}

func (s *OrderService) Stop() error {
	return s.consumer.Stop()
}

func (s *OrderService) handleMessage(msg []byte) error {
	var orderDTO dto.Order
	if err := json.Unmarshal(msg, &orderDTO); err != nil {
		s.logger.Error("invalid message format", zap.Error(err))
		return err
	}

	orderEntity, err := dto.ToOrderEntity(orderDTO)
	if err != nil {
		s.logger.Error("failed to map DTO to entity", zap.Error(err))
		return err
	}

	if err := s.inMemoryRepo.Save(context.Background(), orderEntity); err != nil {
		s.logger.Error("failed to save order in inMemoryRepo", zap.Error(err))
	}

	if err := s.repo.Save(context.Background(), orderEntity); err != nil {
		s.logger.Error("failed to save order in persistent repo", zap.Error(err))
		return err
	}

	s.logger.Info("order processed and saved", zap.String("order_uid", orderEntity.OrderUID))
	return nil
}

func (s *OrderService) GetByID(ctx context.Context, id string) (entity.Order, error) {
	order, err := s.inMemoryRepo.Get(ctx, id)
	if err == nil {
		return order, nil
	}

	order, err = s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get order in OrderRepo", zap.Error(err))
		return entity.Order{}, err
	}

	_ = s.inMemoryRepo.Save(ctx, order)

	return order, nil
}
