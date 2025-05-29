package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "GrpcLargeFileTransfer/api/v1"
)

func uploadFile(client pb.FileServiceClient, filename string) {
	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	buf := make([]byte, 64*1024) // 64KB chunks
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		if err := stream.Send(&pb.UploadRequest{
			Filename: filename,
			Chunk:    buf[:n],
		}); err != nil {
			log.Fatalf("Failed to send chunk: %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive response: %v", err)
	}

	log.Printf("Uploaded file ID: %s, Size: %d bytes", res.Id, res.Size)
}

func main() {

	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(100*1024*1024),
			grpc.MaxCallSendMsgSize(100*1024*1024),
		),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	fileName := flag.String("fileName", "", "Filename to upload")
	flag.Parse()
	client := pb.NewFileServiceClient(conn)
	uploadFile(client, *fileName)
}
