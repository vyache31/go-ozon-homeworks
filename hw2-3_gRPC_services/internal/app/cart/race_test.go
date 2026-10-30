//go:build race

package cart

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vyache31/go-ozon-homeworks/hw2-3_gRPC_services/internal/domain"
	pb "github.com/vyache31/go-ozon-homeworks/hw2-3_gRPC_services/pkg/api"
)

func TestConcurrentPostProduct(t *testing.T) {
	t.Parallel()

	const requests = 50
	service := newTestService(testProductGetter{products: map[domain.ProductSKU]*pb.GetProductResponse{
		123: {Sku: 123, Name: "Vodka", Price: 10},
	}})

	var wg sync.WaitGroup
	wg.Add(requests)
	for range requests {
		go func() {
			defer wg.Done()
			if err := service.PostProduct(context.Background(), domain.PostProductRequest{
				UserID: 1, ProductSKU: 123, Count: 1,
			}); err != nil {
				t.Errorf("PostProduct() error = %v", err)
			}
		}()
	}
	wg.Wait()

	got, err := service.GetProduct(context.Background(), domain.GetProductRequest{UserID: 1})
	require.NoError(t, err)
	require.Equal(t, uint32(requests), got.Items[0].Count)
	require.Equal(t, uint32(requests*10), got.TotalPrice)
}
