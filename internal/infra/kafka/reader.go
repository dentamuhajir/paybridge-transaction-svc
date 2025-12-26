package kafka

import (
	"paybridge-transaction-service/internal/config"

	"github.com/segmentio/kafka-go"
)

func NewReader(cfg *config.Config, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{cfg.Broker.Host},
		Topic:          topic,
		GroupID:        cfg.Server.Name,
		CommitInterval: 0,
	})
}
