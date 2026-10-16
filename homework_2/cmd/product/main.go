package main

import (
	"log"
	"net"

	"github.com/vyache31/go-ozon-homeworks/homework_2/internal/app/product"
	pb "github.com/vyache31/go-ozon-homeworks/homework_2/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting Product service..")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	reflection.Register(server)

	repo := product.NewJSONProductRepository("data/products.json")
	service := product.NewProductService(repo)
	pb.RegisterProductServer(server, service)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}
}
