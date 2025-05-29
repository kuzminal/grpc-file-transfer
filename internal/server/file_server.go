package server

import (
	"fmt"
	"io"
	"log"
	"os"

	pb "GrpcLargeFileTransfer/api/v1"
)

type FileServer struct {
	pb.UnimplementedFileServiceServer
}

func (s *FileServer) Upload(stream pb.FileService_UploadServer) error {
	defer log.Println("Upload complete")
	var fileSize uint32
	var fileName string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.UploadResponse{
				Id:   "generated-file-id",
				Size: fileSize,
			})
		}
		if err != nil {
			return err
		}

		if fileName == "" {
			fileName = req.Filename
			log.Printf("Starting upload for file: %s", fileName)
		}

		chunk := req.Chunk
		fileSize += uint32(len(chunk))
		// Здесь можно сохранять чанки в файл или хранилище
		file, err := os.OpenFile("uploaded-"+fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		file.Write(chunk)
		file.Close()
		log.Printf("Writen chunk with length %d to file %s", len(chunk), fileName)
	}
}

func (s *FileServer) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	// Пример: возвращаем тестовые данные
	data := [][]byte{[]byte("chunk1"), []byte("chunk2")}

	for _, chunk := range data {
		if err := stream.Send(&pb.DownloadResponse{
			Chunk: chunk,
		}); err != nil {
			return err
		}
	}
	return nil
}
