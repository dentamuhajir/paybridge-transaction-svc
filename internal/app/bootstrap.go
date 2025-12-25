package app

import (
	"context"
	"paybridge-transaction-service/internal/config"
	"paybridge-transaction-service/internal/kafka/consumer"
	"paybridge-transaction-service/internal/server"
)

type Bootstrap struct {
	container *Container
}

func NewBootstrap(cfg *config.Config) *Bootstrap {

	ctr, err := NewContainer(cfg)
	if err != nil {
		panic(err)
	}

	return &Bootstrap{container: ctr}
}

func (b *Bootstrap) Start() error {

	ctx := context.Background()

	// start kafka consumers
	walletConsumer := consumer.NewWalletCreateConsumer(b.container.Cfg)
	go walletConsumer.Start(ctx)

	// start HTTP server (blocking)
	return server.Run(
		b.container.Cfg,
		b.container.DB,
		b.container.Logger,
	)
}
