package server

import (
	"paybridge-transaction-service/internal/health"
	"paybridge-transaction-service/internal/wallet"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Dependencies struct {
	DB *pgxpool.Pool
}

func NewRouter(deps *Dependencies) *echo.Echo {
	e := echo.New()
	apiVersion := "/api/v1"

	healthService := health.NewService(deps.DB)
	healthHandler := health.NewHandler(*healthService)
	healthHandler.RegisterRoutes(e.Group(apiVersion))
	walletRepo := wallet.NewRepository(deps.DB)
	walletService := wallet.NewService(walletRepo)
	walletHandler := wallet.NewHandler(walletService)
	walletHandler.RegisterRoutes(*e.Group(apiVersion))

	return e
}
