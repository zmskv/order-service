package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/zmskv/order-service/internal/di"
	"go.uber.org/zap"
)

// @title Order Service API
// @version 1.0
// @description This is the API documentation for the Order Service.
// @host localhost:8000
// @BasePath /
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	container := di.NewContainer(ctx)

	go func() {
		if err := container.OrderService.Start(ctx); err != nil {
			container.Logger.Fatal("failed to start order service", zap.Error(err))
		}
	}()

	go func() {
		container.Logger.Info("starting HTTP server", zap.String("addr", container.HTTPServer.Addr))
		if err := container.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			container.Logger.Fatal("failed to start HTTP server", zap.Error(err))
		}
	}()

	waitForShutdown(cancel, container)
}

func waitForShutdown(cancel context.CancelFunc, container *di.Container) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Logger.Info("Shutting down...")

	cancel()

	if err := container.OrderService.Stop(); err != nil {
		container.Logger.Error("order service stop failed", zap.Error(err))
	}
}
