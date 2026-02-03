package server

import (
	"paybridge-transaction-service/docs"
	"paybridge-transaction-service/internal/account"
	"paybridge-transaction-service/internal/health"
	"paybridge-transaction-service/internal/loan"
	"paybridge-transaction-service/internal/server/middleware"
	"paybridge-transaction-service/internal/wallet"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

func NewRouter(db *pgxpool.Pool, log *zap.Logger) *echo.Echo {

	e := echo.New()

	e.Use(middleware.TracingMiddleware())

	apiVersion := "/api/v1"

	internal := "internal"

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

	loanAppsRepo := loan.NewRepository(db, log)
	loanAppsSvc := loan.NewService(loanAppsRepo, log)
	loanAppsHandler := loan.NewHandler(loanAppsSvc, log)
	loanAppsHandler.RegisterRoutes(e.Group((apiVersion)))

	accountRepo := account.NewRepository(db, log)
	accountSvc := account.NewService(accountRepo, log)
	accountHandler := account.NewHandler(accountSvc, log)
	accountHandler.RegisterInternalRoutes(e.Group(internal))

	return e
}
