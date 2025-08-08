package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zmskv/order-service/internal/domain/order/interfaces"
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
	id := c.Query("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "error invalid request")
		return
	}

	order, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusNotFound, "error not found")
		return
	}

	response.NewSuccessResponse(c, http.StatusOK, "success", order)
}
