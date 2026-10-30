package cart

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/vyache31/go-ozon-homeworks/hw2-3_gRPC_services/internal/app/product"
	"github.com/vyache31/go-ozon-homeworks/hw2-3_gRPC_services/internal/domain"
	pb "github.com/vyache31/go-ozon-homeworks/hw2-3_gRPC_services/pkg/api"
)

type ProductGetter interface {
	GetProduct(context.Context, domain.ProductSKU) (*pb.GetProductResponse, error)
}

const productRequestTimeout = 5 * time.Second

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

	productCtx, cancel := context.WithTimeout(ctx, productRequestTimeout)
	defer cancel()

	product, err := s.products.GetProduct(productCtx, domain.ProductSKU(in.ProductSKU))
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
	if !exists {
		s.mu.RUnlock()
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
	s.mu.RUnlock()

	if len(items) == 0 {
		return domain.GetProductResponse{}, domain.ErrorCartNoExistsOrEmpty
	}

	slices.SortFunc(items, func(a, b domain.Item) int {
		return cmp.Compare(a.Sku, b.Sku)
	})

	ctx, cancel := context.WithTimeout(ctx, productRequestTimeout)
	defer cancel()

	products := make([]*pb.GetProductResponse, len(items))
	errs := make([]error, len(items))
	var wg sync.WaitGroup
	for i := range items {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			products[i], errs[i] = s.products.GetProduct(ctx, domain.ProductSKU(items[i].Sku))
		}(i)
	}
	wg.Wait()

	var totalPrice uint32
	for i, product := range products {
		if errs[i] != nil {
			return domain.GetProductResponse{}, fmt.Errorf("failed to get product: %w", errs[i])
		}
		if product == nil {
			return domain.GetProductResponse{}, fmt.Errorf("failed to get product: empty response")
		}

		items[i].Name = product.Name
		items[i].Price = uint32(product.Price)
		totalPrice += items[i].Price * items[i].Count
	}

	return domain.GetProductResponse{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}
