package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	pb "GrpcLargeFileTransfer/api/v1"
	"GrpcLargeFileTransfer/internal/server"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(100*1024*1024), // 100MB
		grpc.MaxSendMsgSize(100*1024*1024),
	)
	pb.RegisterFileServiceServer(s, &server.FileServer{})

	go func() {
		log.Printf("Server started at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Обработка сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down server...")
	s.GracefulStop()
}
