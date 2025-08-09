package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zmskv/order-service/internal/domain/order/entity"
	mock_interfaces "github.com/zmskv/order-service/internal/domain/order/mocks"
)

func TestOrderService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderService := mock_interfaces.NewMockOrderService(ctrl)

	ctx := context.Background()
	orderID := "123"

	expectedOrder := entity.Order{
		OrderUID: orderID,
	}

	t.Run("success", func(t *testing.T) {
		mockOrderService.EXPECT().
			GetByID(ctx, orderID).
			Return(expectedOrder, nil).
			Times(1)

		order, err := mockOrderService.GetByID(ctx, orderID)
		assert.NoError(t, err)
		assert.Equal(t, expectedOrder, order)
	})

	t.Run("not found", func(t *testing.T) {
		mockOrderService.EXPECT().
			GetByID(ctx, orderID).
			Return(entity.Order{}, errors.New("order not found")).
			Times(1)

		order, err := mockOrderService.GetByID(ctx, orderID)
		assert.Error(t, err)
		assert.Equal(t, entity.Order{}, order)
		assert.EqualError(t, err, "order not found")
	})
}
