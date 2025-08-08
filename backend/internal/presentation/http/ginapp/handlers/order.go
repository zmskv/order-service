package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zmskv/order-service/internal/domain/order/interfaces"
	"github.com/zmskv/order-service/internal/infrastructure/repository/postgres"
	"github.com/zmskv/order-service/internal/infrastructure/response"
	"go.uber.org/zap"
)

type OrderHandler struct {
	service interfaces.OrderService
	logger  *zap.Logger
}

func NewOrderHandler(service interfaces.OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{service: service, logger: logger.Named("handler")}
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	order, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrOrderNotFound) {
			response.NewErrorResponse(c, http.StatusNotFound, "order not found")
			return
		}
		h.logger.Error("failed to get order", zap.Error(err))
		response.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	response.NewSuccessResponse(c, http.StatusOK, "success", order)
}
