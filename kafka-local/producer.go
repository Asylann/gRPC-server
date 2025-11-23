package kafka_local

import (
	notepb "gRPC-server/proto"
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{writer: kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})}
}

func (producer *KafkaProducer) PublishCreateNote(ctx context.Context, userId string, text string) error {
	msg := &notepb.CreateNoteRequest{UserId: userId, NoteText: text}

	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return producer.writer.WriteMessages(ctx, kafka.Message{
		Value: bytes,
	})
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
