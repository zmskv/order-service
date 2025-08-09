package ginapp

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zmskv/order-service/internal/domain/order/interfaces"
	"github.com/zmskv/order-service/internal/presentation/http/ginapp/handlers"
	"go.uber.org/zap"
)

func InitRoutes(r *gin.Engine, orderService interfaces.OrderService, logger *zap.Logger) {
	orderHandler := handlers.NewOrderHandler(orderService, logger)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	order := r.Group("/order")
	{
		order.GET("/:id", orderHandler.GetByID)
	}
}
