package main

import (
	"fmt"
	"os"
	"os/signal"
	"sample-crud/infra/cache"
	"sample-crud/infra/db"
	"sample-crud/internal/config"
	"sample-crud/internal/handler"
	"sample-crud/internal/middleware"
	"sample-crud/internal/repo"
	"sample-crud/internal/service"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	commonhandler "github.com/tee-nullpointer/go-common-kit/handler"
	commonmiddleware "github.com/tee-nullpointer/go-common-kit/middleware"
	"github.com/tee-nullpointer/go-common-kit/pkg/logger"
	"github.com/tee-nullpointer/go-common-kit/server"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	if err := logger.InitLogger(cfg.Logger.Level, cfg.Logger.Format); err != nil {
		panic(fmt.Sprintf("Fail to initialize logger: %v", err))
	}
	defer logger.Sync()

	gormDB := db.Init(cfg.Database)
	defer db.ShutDown()

	redisClient := cache.NewRedisClient(cfg.Redis)
	defer cache.Close()

	ginServer := server.NewGinServer(cfg.Server.Mode)
	ginRouter := ginServer.GetRouter()
	setupRouter(ginRouter, gormDB, redisClient)
	go ginServer.Start(cfg.Server.Host, cfg.Server.Port)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	ginServer.GracefulShutdown()
}

func setupRouter(router *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	router.Use(gin.Recovery())
	router.Use(commonmiddleware.LoggingMiddleware())
	router.Use(middleware.ErrorRecover())

	productRepo := repo.NewGormProductRepository(db)
	productService := service.NewProductService(productRepo, redisClient)
	productHandler := handler.NewProductHandler(productService)
	monitor := router.Group("/")
	{
		monitor.GET("/health", commonhandler.HealthCheck)
	}
	v1 := router.Group("/api/v1")
	{
		v1.POST("/products", productHandler.Create)
		v1.GET("/products/:id", productHandler.FindByID)
		v1.PUT("/products/:id", productHandler.Update)
		v1.DELETE("/products/:id", productHandler.Delete)
	}
}
