package product

import (
	"context"
	"errors"

	"github.com/vyache31/go-ozon-homeworks/homework_2/internal/domain"
	pb "github.com/vyache31/go-ozon-homeworks/homework_2/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductRepository interface {
	Get(sku domain.ProductSKU) (domain.Product, error)
}

type ProductService struct {
	pb.UnimplementedProductServer

	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, in *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	if in.Sku <= 0 {
		return nil, status.Error(codes.InvalidArgument, domain.ErrorProductInvalidSKU.Error())
	}

	product, err := s.repo.Get(domain.ProductSKU(in.Sku))
	if errors.Is(err, domain.ErrorProductNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, domain.ErrorProductInternal.Error())
	}

	reponse := pb.GetProductResponse{
		Sku:   product.Sku,
		Name:  product.Name,
		Price: int32(product.Price),
	}

	return &reponse, nil
}
