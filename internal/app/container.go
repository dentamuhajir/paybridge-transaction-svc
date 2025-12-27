package app

import (
	"paybridge-transaction-service/internal/config"
	"paybridge-transaction-service/internal/domain/wallet"
	"paybridge-transaction-service/internal/infra/logger"
	"paybridge-transaction-service/internal/infra/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Service struct {
	WalletService wallet.Service
}

type Container struct {
	Cfg     *config.Config
	DB      *pgxpool.Pool
	Logger  *zap.Logger
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

	walletRepo := wallet.NewRepository(db, log)
	walletSvc := wallet.NewService(walletRepo, log)

	return &Container{
		Cfg:    cfg,
		DB:     db,
		Logger: log,
		Service: &Service{
			WalletService: walletSvc,
		},
	}, nil
}
