package app

import (
	"paybridge-transaction-service/internal/account"
	"paybridge-transaction-service/internal/config"
	"paybridge-transaction-service/internal/infra/logger"
	"paybridge-transaction-service/internal/infra/postgres"
	"paybridge-transaction-service/internal/ledger"
	"paybridge-transaction-service/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	OpenAccountUsecase *usecase.OpenAccountUsecase
}

type Container struct {
	Cfg     *config.Config
	DB      *pgxpool.Pool
	Logger  *logger.Logger
	Service *Service
}

func NewContainer(cfg *config.Config) (*Container, error) {
	db, err := postgres.NewPostgres(cfg.Database.DSN)
	if err != nil {
		return nil, err
	}

	log, err := logger.New()
	if err != nil {
		return nil, err
	}

	accountRepo := account.NewRepository(db, log)

	ledgerRepo := ledger.NewRepository(db, log)
	openAccountUC := usecase.NewOpenAccountUsecase(
		db,
		accountRepo,
		ledgerRepo,
	)

	return &Container{
		Cfg:    cfg,
		DB:     db,
		Logger: log,
		Service: &Service{
			OpenAccountUsecase: openAccountUC,
		},
	}, nil
}
