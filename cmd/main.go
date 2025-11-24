package main

import (
	"context"
	"github.com/Asylann/gRPC-server/internal/repository"
	kafka_local "github.com/Asylann/gRPC-server/kafka-local"
	notepb "github.com/Asylann/gRPC-server/proto"
	"github.com/Asylann/gRPC-server/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	repo, _ := repository.NewRepository()
	serv := service.NewService(repo)

	// ----------------------
	// gRPC server
	// ----------------------
	grpcServer := grpc.NewServer()
	notepb.RegisterNoteServiceServer(grpcServer, serv)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("gRPC NoteService running on port 50051")

	// ----------------------
	// Kafka consumer
	// ----------------------
	brokers := []string{"kafka:9092"}

	consumer := kafka_local.NewConsumer(
		brokers,
		"notes.create",
		"notes_service_group",
		*serv,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go consumer.Run(ctx)

	// ----------------------
	// gRPC server run
	// ----------------------
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// ----------------------
	// Graceful shutdown
	// ----------------------
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	log.Println("Shutting down gRPC and Kafka...")

	cancel()
	grpcServer.GracefulStop()
	log.Println("Stopped.")
}
