package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zmskv/order-service/internal/domain/order/interfaces"
	"go.uber.org/zap"
)

type OrderHandler struct {
	service interfaces.OrderService
	logger  zap.Logger
}

func NewOrderHandler(service interfaces.OrderService, logger zap.Logger) *OrderHandler {
	return &OrderHandler{service: service, logger: logger}
}

func (h *OrderHandler) GetByID(c *gin.Context) {

}
