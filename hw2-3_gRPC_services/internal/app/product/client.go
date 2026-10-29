package product

import (
	"context"
	"fmt"

	"github.com/vyache31/go-ozon-homeworks/homework_2/internal/domain"
	pb "github.com/vyache31/go-ozon-homeworks/homework_2/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type ProductServiceClient struct {
	client pb.ProductClient
}

func NewProductClient(address string) (*ProductServiceClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc conn: %w", err)
	}

	return &ProductServiceClient{
		client: pb.NewProductClient(conn),
	}, nil
}

func (c *ProductServiceClient) GetProduct(ctx context.Context, sku domain.ProductSKU) (*pb.GetProductResponse, error) {
	resp, err := c.client.GetProduct(ctx, &pb.GetProductRequest{
		Sku: uint32(sku),
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, domain.ErrorProductNotFound
		}

		return nil, fmt.Errorf("client failed to get product: %w", err)
	}

	return resp, nil
}
