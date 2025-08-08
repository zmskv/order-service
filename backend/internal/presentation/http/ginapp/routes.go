package ginapp

import (
	"github.com/gin-gonic/gin"
	"github.com/zmskv/order-service/internal/domain/order/interfaces"
	"github.com/zmskv/order-service/internal/presentation/http/ginapp/handlers"
	"go.uber.org/zap"
)

func InitRoutes(r *gin.Engine, orderService interfaces.OrderService, logger zap.Logger) {
	orderHandler := handlers.NewOrderHandler(orderService, logger)

	order := r.Group("/order")
	{
		order.GET("/:id", orderHandler.GetByID)
	}
}
