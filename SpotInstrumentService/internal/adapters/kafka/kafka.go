package adapterkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	spotconfig "github.com/DencCPU/gRPCServices/SpotInstrumentService/config"
	outboxdomain "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/domain/outbox"
	"github.com/segmentio/kafka-go"
)

type KafkaBroker struct {
	writer kafka.Writer
}

func NewKafkaBroker(cfg spotconfig.Kafka) *KafkaBroker {
	addr := cfg.Host + ":" + cfg.Port
	fmt.Println(addr)
	broker := KafkaBroker{}
	broker.writer = kafka.Writer{
		Addr:            kafka.TCP(addr),
		Topic:           cfg.Topic,
		Balancer:        &kafka.LeastBytes{},
		MaxAttempts:     cfg.MaxAttempts,
		WriteBackoffMin: cfg.WriteBackoffMin,
		WriteBackoffMax: cfg.WriteBackoffMax,
		BatchSize:       cfg.BatchSize,
		BatchBytes:      cfg.BatchBytes,
	}
	err := broker.CreateTopic(addr)
	if err != nil {
		fmt.Println(err)
	}
	return &broker
}

func (k *KafkaBroker) Send(ctx context.Context, event *outboxdomain.MarketEvent) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		key := []byte(event.BasedEvent.EventId)
		value, err := json.Marshal(event)
		if err != nil {
			return err
		}

		message := kafka.Message{
			Key:   key,
			Value: value,
		}
		err = k.writer.WriteMessages(ctx, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *KafkaBroker) CreateTopic(addr string) error {
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer conn.Close()

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             k.writer.Topic,
		NumPartitions:     3,
		ReplicationFactor: 1,
	})

	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf("failed to create topic: %w", err)
	}
	return nil
}
