package server

import (
	"paybridge-transaction-service/docs"
	"paybridge-transaction-service/internal/domain/health"
	"paybridge-transaction-service/internal/server/middleware"
	"paybridge-transaction-service/internal/domain/wallet"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

func NewRouter(db *pgxpool.Pool, log *zap.Logger) *echo.Echo {

	e := echo.New()

	e.Use(middleware.TracingMiddleware())

	apiVersion := "/api/v1"

	// Swagger Information
	docs.SwaggerInfo.Title = "Paybridge Transaction Service API"
	docs.SwaggerInfo.Description = "API documentation for Transaction Services"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = apiVersion

	// Swagger route
	e.GET("/swagger-ui/*", echoSwagger.WrapHandler)

	healthService := health.NewService(db)
	healthHandler := health.NewHandler(*healthService)
	healthHandler.RegisterRoutes(e.Group(apiVersion))
	walletRepo := wallet.NewRepository(db, log)
	walletService := wallet.NewService(walletRepo, log)
	walletHandler := wallet.NewHandler(walletService, log)
	walletHandler.RegisterRoutes(e.Group(apiVersion))

	return e
}
