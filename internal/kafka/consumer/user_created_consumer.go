package consumer

import (
	"context"
	"encoding/json"
	"log"
	"paybridge-transaction-service/internal/account"
	"paybridge-transaction-service/internal/config"
	kafkaInfra "paybridge-transaction-service/internal/infra/kafka"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type CreateUserEvent struct {
	UserID uuid.UUID `json:"user_id"`
}

type UserCreateConsumer struct {
	reader  *kafka.Reader
	Service account.Service
}

func NewUserCreateConsumer(cfg *config.Config, svc account.Service) *UserCreateConsumer {
	return &UserCreateConsumer{
		reader:  kafkaInfra.NewReader(cfg, "account"),
		Service: svc,
	}
}

func (c *UserCreateConsumer) Start(ctx context.Context) {
	defer c.reader.Close()

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Println("Fetch error:", err)
			continue
		}

		var event CreateUserEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("Invalid message:", err)
			_ = c.reader.CommitMessages(ctx, msg) // skip bad message
			continue
		}

		if err := c.Service.CreateAccountWithInitialBalances(ctx, event.UserID); err != nil {
			log.Println("Create wallet failed:", err)
			continue // DO NOT COMMIT → Kafka will retry
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Println("Commit failed:", err)
		}

		log.Printf("[Kafka] account initialized user=%s", event.UserID)
	}
}
