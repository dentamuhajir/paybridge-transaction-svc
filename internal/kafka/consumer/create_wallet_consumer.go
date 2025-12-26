package consumer

import (
	"context"
	"fmt"
	"log"
	"paybridge-transaction-service/internal/config"
	kafkaInfra "paybridge-transaction-service/internal/infra/kafka"

	"github.com/segmentio/kafka-go"
)

type WalletCreateEvent struct {
	UserID   string `json:"userId"`
	Currency string `json:"currency"`
}

type WalletCreateConsumer struct {
	reader *kafka.Reader
	//service *wallet.Service
}

func NewWalletCreateConsumer(cfg *config.Config) *WalletCreateConsumer {
	return &WalletCreateConsumer{
		reader: kafkaInfra.NewReader(cfg, "wallet"),
		// service: s,
	}
}

func (c *WalletCreateConsumer) Start(ctx context.Context) {
	defer c.reader.Close()

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Println("Fetch error:", err)
			continue
		}

		fmt.Printf(
			"[Kafka] topic=%s key=%s value=%s offset=%d\n",
			msg.Topic,
			string(msg.Key),
			string(msg.Value),
			msg.Offset,
		)

		_ = c.reader.CommitMessages(ctx, msg)
	}
}
