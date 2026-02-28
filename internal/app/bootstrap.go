package app

import (
	"context"
	"paybridge-transaction-service/internal/config"
	"paybridge-transaction-service/internal/infra/otel"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := otel.InitTracer(ctx, b.container.Cfg.Server.Name, b.container.Cfg.Otel.OtelExporterEndpoint)
	if err != nil {
		return err
	}
	defer func() {
		_ = shutdown(context.Background())
	}()

	// start kafka consumers
	userCreatedConsumer := consumer.NewUserCreateConsumer(
		b.container.Cfg,
		b.container.Service.OpenAccountUsecase,
	)

	go userCreatedConsumer.Start(ctx)

	// start HTTP server (blocking)
	return server.Run(
		b.container.Cfg,
		b.container.DB,
		b.container.Logger,
	)
}
