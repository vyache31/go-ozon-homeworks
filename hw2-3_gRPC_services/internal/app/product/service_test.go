package product

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vyache31/go-ozon-homeworks/hw2-3_gRPC_services/internal/domain"
	pb "github.com/vyache31/go-ozon-homeworks/hw2-3_gRPC_services/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type testRepository struct {
	product domain.Product
	err     error
}

func (r testRepository) Get(domain.ProductSKU) (domain.Product, error) {
	return r.product, r.err
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		sku        uint32
		repository testRepository
		want       *pb.GetProductResponse
		wantCode   codes.Code
	}{
		{
			name:     "wrong sku",
			sku:      0,
			wantCode: codes.InvalidArgument,
		},
		{
			name:       "not exist",
			sku:        666,
			repository: testRepository{err: domain.ErrorProductNotFound},
			wantCode:   codes.NotFound,
		},
		{
			name:       "repository error",
			sku:        123,
			repository: testRepository{err: errors.New("database error")},
			wantCode:   codes.Internal,
		},
		{
			name: "sucess",
			sku:  777,
			repository: testRepository{
				product: domain.Product{Sku: 777, Name: "Vodka", Price: 52},
			},
			want: &pb.GetProductResponse{Sku: 777, Name: "Vodka", Price: 52},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := NewProductService(tt.repository)
			got, err := service.GetProduct(context.Background(), &pb.GetProductRequest{Sku: tt.sku})

			if tt.wantCode != codes.OK {
				require.Error(t, err)
				require.Nil(t, got)
				require.Equal(t, tt.wantCode, status.Code(err))
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
