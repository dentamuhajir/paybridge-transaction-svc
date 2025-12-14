package app

import (
	"paybridge-transaction-service/internal/config"
	"paybridge-transaction-service/internal/db"
	"paybridge-transaction-service/internal/server"
)

type Bootstrap struct {
	cfg *config.Config
}

func NewBootstrap(cfg *config.Config) *Bootstrap {
	return &Bootstrap{cfg: cfg}
}

func (b *Bootstrap) Start() error {

	// Infrastructure
	pg, err := db.NewPostgres(b.cfg.Database.DSN)
	if err != nil {
		return err
	}

	deps := &server.Dependencies{
		DB: pg,
	}

	// HTTP
	// go func() {
	if err := server.Run(b.cfg, deps); err != nil {
		return err
	}
	// }()

	// Kafka consumers start here (later)
	// go kafka.StartWalletConsumer(...)

	select {}
}
