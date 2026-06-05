package adapterkafka

import (
	"context"
	"encoding/json"

	orderconfig "github.com/DencCPU/gRPCServices/OrderService/config"
	inboxdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/inbox"
	"github.com/segmentio/kafka-go"
)

type KafkaBroker struct {
	reader kafka.Reader
}

func NewKafkaBroker(cfg orderconfig.Kafka) *KafkaBroker {
	broker := KafkaBroker{}

	broker.reader = *kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:        cfg.Brokers,
			Topic:          cfg.Topic,
			GroupID:        cfg.GroupID,
			StartOffset:    kafka.LastOffset,
			MinBytes:       cfg.MinBytes,
			MaxBytes:       cfg.MaxBytes,
			MaxWait:        cfg.MaxWait,
			CommitInterval: 0,
		},
	)
	return &broker
}

func (k *KafkaBroker) ReadMessage(ctx context.Context) (inboxdomain.KafkaMessage, error) {

	msg, err := k.reader.FetchMessage(ctx)
	if err != nil {
		return inboxdomain.KafkaMessage{}, err
	}

	var kafkaMessage inboxdomain.KafkaMessage
	kafkaMessage.Key = string(msg.Key)
	kafkaMessage.Offset = msg.Offset
	kafkaMessage.Partition = msg.Partition

	err = json.Unmarshal(msg.Value, &kafkaMessage.Value)
	if err != nil {
		return inboxdomain.KafkaMessage{}, err
	}
	return kafkaMessage, nil
}

func (k *KafkaBroker) Commit(ctx context.Context, msg inboxdomain.KafkaMessage) error {
	return k.reader.CommitMessages(ctx, kafka.Message{
		Topic:     k.reader.Config().Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
	})
}

func (k *KafkaBroker) Close() {
	k.reader.Close()
}
