package cart

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vyache31/go-ozon-homeworks/hw2-3_gRPC_services/internal/domain"
	pb "github.com/vyache31/go-ozon-homeworks/hw2-3_gRPC_services/pkg/api"
)

type testProductGetter struct {
	products map[domain.ProductSKU]*pb.GetProductResponse
	err      error
}

func (g testProductGetter) GetProduct(_ context.Context, sku domain.ProductSKU) (*pb.GetProductResponse, error) {
	return g.products[sku], g.err
}

func newTestService(getter ProductGetter) *CartService {
	return &CartService{
		usersCarts: make(map[domain.UserID]domain.Cart),
		products:   getter,
	}
}

func TestPostProduct(t *testing.T) {
	t.Parallel()

	serviceError := errors.New("service error")
	tests := []struct {
		name     string
		requests []domain.PostProductRequest
		getter   testProductGetter
		wantErr  error
		wantCart domain.GetProductResponse
	}{
		{
			name: "bad user",
			requests: []domain.PostProductRequest{
				{UserID: 0, ProductSKU: 123, Count: 1},
			},
			wantErr: domain.ErrorUserInvalidID,
		},
		{
			name: "bad sku",
			requests: []domain.PostProductRequest{
				{UserID: 1, ProductSKU: 0, Count: 1},
			},
			wantErr: domain.ErrorProductInvalidSKU,
		},
		{
			name: "bad count",
			requests: []domain.PostProductRequest{
				{UserID: 1, ProductSKU: 123, Count: 0},
			},
			wantErr: domain.ErrorPositionInvalidCount,
		},
		{
			name: "product service error",
			requests: []domain.PostProductRequest{
				{UserID: 1, ProductSKU: 123, Count: 1},
			},
			getter:  testProductGetter{err: serviceError},
			wantErr: serviceError,
		},
		{
			name: "add product and update total",
			requests: []domain.PostProductRequest{
				{UserID: 1, ProductSKU: 123, Count: 2},
			},
			getter: testProductGetter{products: map[domain.ProductSKU]*pb.GetProductResponse{
				123: {Sku: 123, Name: "Book", Price: 500},
			}},
			wantCart: domain.GetProductResponse{
				Items:      []domain.Item{{Sku: 123, Name: "Book", Price: 500, Count: 2}},
				TotalPrice: 1000,
			},
		},
		{
			name: "add same product increas count",
			requests: []domain.PostProductRequest{
				{UserID: 1, ProductSKU: 123, Count: 2},
				{UserID: 1, ProductSKU: 123, Count: 3},
			},
			getter: testProductGetter{products: map[domain.ProductSKU]*pb.GetProductResponse{
				123: {Sku: 123, Name: "Vodka", Price: 500},
			}},
			wantCart: domain.GetProductResponse{
				Items:      []domain.Item{{Sku: 123, Name: "Vodka", Price: 500, Count: 5}},
				TotalPrice: 2500,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := newTestService(tt.getter)
			for _, request := range tt.requests {
				err := service.PostProduct(context.Background(), request)
				if tt.wantErr != nil {
					require.ErrorIs(t, err, tt.wantErr)
					break
				}
				require.NoError(t, err)
			}

			if tt.wantErr != nil {
				return
			}

			got, err := service.GetProduct(context.Background(), domain.GetProductRequest{UserID: 1})
			require.NoError(t, err)
			require.Equal(t, tt.wantCart, got)
		})
	}
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  int64
		setup   []domain.PostProductRequest
		getter  testProductGetter
		getErr  error
		want    domain.GetProductResponse
		wantErr error
	}{
		{
			name:    "bad user",
			userID:  0,
			wantErr: domain.ErrorUserInvalidID,
		},
		{
			name:    "empty cart",
			userID:  1,
			wantErr: domain.ErrorCartNoExistsOrEmpty,
		},
		{
			name:   "items sorted by sku",
			userID: 1,
			setup: []domain.PostProductRequest{
				{UserID: 1, ProductSKU: 42, Count: 2},
				{UserID: 1, ProductSKU: 8, Count: 3},
			},
			getter: testProductGetter{products: map[domain.ProductSKU]*pb.GetProductResponse{
				42: {Sku: 42, Name: "Display", Price: 700},
				8:  {Sku: 8, Name: "Keyboard", Price: 150},
			}},
			want: domain.GetProductResponse{
				Items: []domain.Item{
					{Sku: 8, Name: "Keyboard", Price: 150, Count: 3},
					{Sku: 42, Name: "Display", Price: 700, Count: 2},
				},
				TotalPrice: 1850,
			},
		},
		{
			name:   "product disappeared",
			userID: 1,
			setup: []domain.PostProductRequest{
				{UserID: 1, ProductSKU: 42, Count: 1},
			},
			getter: testProductGetter{products: map[domain.ProductSKU]*pb.GetProductResponse{
				42: {Sku: 42, Name: "Display", Price: 700},
			}},
			getErr:  domain.ErrorProductNotFound,
			wantErr: domain.ErrorProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := newTestService(tt.getter)
			for _, req := range tt.setup {
				require.NoError(t, service.PostProduct(context.Background(), req))
			}
			if tt.getErr != nil {
				service.products = testProductGetter{err: tt.getErr}
			}

			got, err := service.GetProduct(context.Background(), domain.GetProductRequest{UserID: tt.userID})
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, domain.GetProductResponse{}, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetCartProductError(t *testing.T) {
	t.Parallel()

	service := newTestService(testProductGetter{products: map[domain.ProductSKU]*pb.GetProductResponse{
		123: {Sku: 123, Name: "Book", Price: 500},
	}})
	require.NoError(t, service.PostProduct(context.Background(), domain.PostProductRequest{
		UserID: 1, ProductSKU: 123, Count: 1,
	}))
	service.products = testProductGetter{err: domain.ErrorProductNotFound}

	request := httptest.NewRequest(http.MethodGet, "/user/1/cart", nil)
	request.SetPathValue("user_id", "1")
	response := httptest.NewRecorder()

	NewCartHandler(service).GetCart(response, request)

	require.Equal(t, http.StatusPreconditionFailed, response.Code)
}
