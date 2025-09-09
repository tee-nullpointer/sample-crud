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
	"sample-crud/proto/pb/product"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	commonhandler "github.com/tee-nullpointer/go-common-kit/handler"
	commoninterceptor "github.com/tee-nullpointer/go-common-kit/interceptor"
	commonmiddleware "github.com/tee-nullpointer/go-common-kit/middleware"
	"github.com/tee-nullpointer/go-common-kit/pkg/logger"
	"github.com/tee-nullpointer/go-common-kit/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	productRepo := repo.NewGormProductRepository(gormDB)
	productService := service.NewProductService(productRepo, redisClient)
	setupRouter(ginRouter, productService)
	go ginServer.Start(cfg.Server.Host, cfg.Server.Port)

	grpcServer := server.NewGRPCServer(
		grpc.UnaryInterceptor(
			commoninterceptor.ChainUnaryInterceptors(
				commoninterceptor.RecoveryUnaryInterceptor,
				commoninterceptor.TraceUnaryInterceptor,
				commoninterceptor.LoggingUnaryInterceptor,
			),
		),
	)
	product.RegisterProductServiceServer(grpcServer.GetServer(), handler.NewProductGRPCHandler(productService))
	go grpcServer.Start("localhost", strconv.Itoa(cfg.Grpc.Port))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	zap.L().Info("Received shutdown signal", zap.String("signal", sig.String()))
	ginServer.GracefulShutdown()
}

func setupRouter(router *gin.Engine, productService service.ProductService) {
	router.Use(gin.Recovery())
	router.Use(commonmiddleware.TraceMiddleware())
	router.Use(commonmiddleware.LoggingMiddleware())
	router.Use(middleware.ErrorRecover())
	productHandler := handler.NewProductApiHandler(productService)
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
