package consumer

import (
	"context"
	"encoding/json"
	"log"
	"paybridge-transaction-service/internal/config"
	"paybridge-transaction-service/internal/event"
	kafkaInfra "paybridge-transaction-service/internal/infra/kafka"
	"paybridge-transaction-service/internal/usecase"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type CreateUserEvent struct {
	UserID uuid.UUID `json:"userId"`
}

type UserCreateConsumer struct {
	reader  *kafka.Reader
	usecase *usecase.OpenAccountUsecase
}

func NewUserCreateConsumer(
	cfg *config.Config,
	uc *usecase.OpenAccountUsecase,
) *UserCreateConsumer {
	return &UserCreateConsumer{
		reader:  kafkaInfra.NewReader(cfg, event.UserCreatedTopic),
		usecase: uc,
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

		if err := c.usecase.ExecuteUserAccount(ctx, event.UserID); err != nil {
			log.Println("Create account failed:", err)
			continue // DO NOT COMMIT → Kafka retry
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Println("Commit failed:", err)
		}

		log.Printf("[Kafka] account initialized user=%s", event.UserID)
	}
}
