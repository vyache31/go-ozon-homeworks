package cart

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/vyache31/go-ozon-homeworks/homework_2/internal/app/product"
	"github.com/vyache31/go-ozon-homeworks/homework_2/internal/domain"
	pb "github.com/vyache31/go-ozon-homeworks/homework_2/pkg/api"
)

type ProductGetter interface {
	GetProduct(context.Context, domain.ProductSKU) (*pb.GetProductResponse, error)
}

type CartService struct {
	usersCarts map[domain.UserID]domain.Cart
	products   ProductGetter

	mu sync.RWMutex
}

func NewCartService() (*CartService, error) {
	productClient, err := product.NewProductClient(":50051")
	if err != nil {
		return nil, fmt.Errorf("failed to create product client: %w", err)
	}

	return &CartService{
		usersCarts: make(map[domain.UserID]domain.Cart),
		products:   productClient,
		mu:         sync.RWMutex{},
	}, nil
}

func (s *CartService) PostProduct(ctx context.Context, in domain.PostProductRequest) error {
	if in.UserID <= 0 {
		return domain.ErrorUserInvalidID
	}
	if in.ProductSKU <= 0 {
		return domain.ErrorProductInvalidSKU
	}
	if in.Count <= 0 {
		return domain.ErrorPositionInvalidCount
	}

	product, err := s.products.GetProduct(ctx, domain.ProductSKU(in.ProductSKU))
	if err != nil {
		return fmt.Errorf("failed to check product: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	userCart := s.usersCarts[domain.UserID(in.UserID)]
	if userCart.Positions == nil {
		userCart.Positions = make(map[domain.ProductSKU]domain.Position)
	}

	position, exists := userCart.Positions[domain.ProductSKU(product.Sku)]
	if exists {
		position.Count += in.Count
	} else {
		position = domain.Position{
			Product: domain.Product{
				Sku:   product.Sku,
				Name:  product.Name,
				Price: uint32(product.Price),
			},
			Count: in.Count,
		}
	}

	userCart.Positions[domain.ProductSKU(product.Sku)] = position
	userCart.TotalPrice += uint32(product.Price) * in.Count
	s.usersCarts[domain.UserID(in.UserID)] = userCart

	return nil
}

func (s *CartService) GetProduct(ctx context.Context, in domain.GetProductRequest) (domain.GetProductResponse, error) {
	if in.UserID <= 0 {
		return domain.GetProductResponse{}, domain.ErrorUserInvalidID
	}

	s.mu.RLock()
	cart, exists := s.usersCarts[domain.UserID(in.UserID)]
	s.mu.RUnlock()

	if !exists {
		return domain.GetProductResponse{}, domain.ErrorCartNoExistsOrEmpty
	}

	items := make([]domain.Item, 0, len(cart.Positions))
	for _, position := range cart.Positions {
		items = append(items, domain.Item{
			Sku:   int64(position.Product.Sku),
			Name:  position.Product.Name,
			Price: position.Product.Price,
			Count: position.Count,
		})
	}
	slices.SortFunc(items, func(a, b domain.Item) int {
		return cmp.Compare(a.Sku, b.Sku)
	})

	return domain.GetProductResponse{
		Items:      items,
		TotalPrice: cart.TotalPrice,
	}, nil
}
