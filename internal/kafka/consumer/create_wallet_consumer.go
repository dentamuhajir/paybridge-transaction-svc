package consumer

import (
	"context"
	"encoding/json"
	"log"
	"paybridge-transaction-service/internal/config"
	"paybridge-transaction-service/internal/domain/wallet"
	kafkaInfra "paybridge-transaction-service/internal/infra/kafka"

	"github.com/segmentio/kafka-go"
)

type WalletCreateEvent struct {
	UserID   string `json:"userId"`
	Currency string `json:"currency"`
}

type WalletCreateConsumer struct {
	reader  *kafka.Reader
	Service wallet.Service
}

func NewWalletCreateConsumer(cfg *config.Config, svc wallet.Service) *WalletCreateConsumer {
	return &WalletCreateConsumer{
		reader:  kafkaInfra.NewReader(cfg, "wallet"),
		Service: svc,
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
 
		var event WalletCreateEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("Invalid message:", err)
			_ = c.reader.CommitMessages(ctx, msg) // skip bad message
			continue
		}

		req := wallet.CreateWalletReq{
			UserID:   event.UserID,
			Currency: event.Currency,
		}

		if _, err := c.Service.CreateWallet(ctx, req); err != nil {
			log.Println("Create wallet failed:", err)
			continue // DO NOT COMMIT → Kafka will retry
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Println("Commit failed:", err)
		}

		log.Printf("[Kafka] wallet created user=%s currency=%s", event.UserID, event.Currency)
	}
}
