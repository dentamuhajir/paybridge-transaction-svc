package server

import (
	"paybridge-transaction-service/docs"
	"paybridge-transaction-service/internal/health"
	"paybridge-transaction-service/internal/logger"
	"paybridge-transaction-service/internal/server/middleware"
	"paybridge-transaction-service/internal/wallet"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Dependencies struct {
	DB *pgxpool.Pool
}

func NewRouter(deps *Dependencies) *echo.Echo {

	logger.Init()

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

	healthService := health.NewService(deps.DB)
	healthHandler := health.NewHandler(*healthService)
	healthHandler.RegisterRoutes(e.Group(apiVersion))
	walletRepo := wallet.NewRepository(deps.DB)
	walletService := wallet.NewService(walletRepo)
	walletHandler := wallet.NewHandler(walletService)
	walletHandler.RegisterRoutes(e.Group(apiVersion))

	return e
}
