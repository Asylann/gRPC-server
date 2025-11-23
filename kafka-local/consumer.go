package kafka_local

import (
	notepb "gRPC-server/proto"
	"gRPC-server/service"
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"log"
)

type NoteConsumer struct {
	reader      *kafka.Reader
	NoteService *service.NoteService
}

func NewConsumer(brokers []string, topic, group string, noteService service.NoteService) *NoteConsumer {
	return &NoteConsumer{kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: group,
	}), &noteService}
}

func (consumer *NoteConsumer) Run(ctx context.Context) {
	for {
		m, err := consumer.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}
		var req notepb.CreateNoteRequest
		if err := proto.Unmarshal(m.Value, &req); err != nil {
			log.Println("Proto unmarshal error:", err)
			continue
		}

		if _, err := consumer.NoteService.CreateNote(ctx, &req); err != nil {
			log.Println("CreateNote error:", err)
		}

		log.Println("LOG: a New note was may created!")
	}
}
