package di

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"go.uber.org/zap"

	"github.com/zmskv/order-service/internal/application"
	"github.com/zmskv/order-service/internal/domain/order/interfaces"
	"github.com/zmskv/order-service/internal/infrastructure/messaging/kafka"
	"github.com/zmskv/order-service/internal/infrastructure/repository/inmemory"
	"github.com/zmskv/order-service/internal/infrastructure/repository/postgres"
	"github.com/zmskv/order-service/internal/presentation/http/ginapp"
	"github.com/zmskv/order-service/logger"
)

type Container struct {
	Config       Config
	Logger       *zap.Logger
	DB           *sqlx.DB
	OrderRepo    interfaces.OrderRepository
	InMemoryRepo interfaces.OrderRepository
	Consumer     *kafka.Consumer
	OrderService *application.OrderService
	HTTPServer   *http.Server
}

func NewContainer(ctx context.Context) *Container {
	config := ReadConfig()

	log := logger.New()
	defer log.Sync()

	db, err := InitDB(config.DB)

	if err != nil {
		log.Fatal("failed to initialize database", zap.Error(err))
	}

	orderRepo := postgres.NewOrderRepository(db, log)
	inMemoryRepo, err := inmemory.NewInMemoryRepository(log)
	if err != nil {
		log.Fatal("failed to create InMemoryRepository", zap.Error(err))
	}

	var consumer *kafka.Consumer
	consumer, err = kafka.NewConsumer(
		config.Kafka.Brokers,
		config.Kafka.Topic,
		config.Kafka.GroupID,
		func(data []byte) error {
			log.Info("Kafka message received but handler not set yet")
			return nil
		},
		log,
	)
	if err != nil {
		log.Fatal("failed to create kafka consumer", zap.Error(err))
	}

	orderService := application.NewOrderService(orderRepo, consumer, inMemoryRepo, log)

	httpServer := InitHTTPServer(config.HTTP, orderService, log)

	container := &Container{
		Config:       config,
		Logger:       log,
		DB:           db,
		OrderRepo:    orderRepo,
		InMemoryRepo: inMemoryRepo,
		Consumer:     consumer,
		OrderService: orderService,
		HTTPServer:   httpServer,
	}

	if err := loadCacheFromDB(ctx, container); err != nil {
		log.Fatal("failed to load cache from DB", zap.Error(err))
	}
	return container
}

func InitDB(cfg DBConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

func loadCacheFromDB(ctx context.Context, c *Container) error {
	c.Logger.Info("loading cache from database")

	orders, err := c.OrderRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if err := c.InMemoryRepo.Save(ctx, order); err != nil {
			c.Logger.Warn("failed to add order to in-memory repo", zap.Error(err), zap.Any("order", order))
		}
	}

	c.Logger.Info("cache successfully loaded", zap.Int("orders_count", len(orders)))
	return nil
}

func InitHTTPServer(cfg HTTPConfig, orderService interfaces.OrderService, logger *zap.Logger) *http.Server {
	router := gin.Default()
	ginapp.InitRoutes(router, orderService, logger)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}
