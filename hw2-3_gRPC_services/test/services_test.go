//go:build api

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	pb "github.com/vyache31/go-ozon-homeworks/hw2-3_gRPC_services/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestProductService(t *testing.T) {
	client := connectProduct(t)
	tests := []struct {
		name     string
		sku      uint32
		wantCode codes.Code
	}{
		{name: "product good", sku: 1076963, wantCode: codes.OK},
		{name: "zero sku", sku: 0, wantCode: codes.InvalidArgument},
		{name: "product not found", sku: 666, wantCode: codes.NotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			got, err := client.GetProduct(ctx, &pb.GetProductRequest{Sku: tt.sku})
			require.Equal(t, tt.wantCode, status.Code(err))
			if tt.wantCode == codes.OK {
				require.NoError(t, err)
				require.Equal(t, tt.sku, got.Sku)
				require.NotEmpty(t, got.Name)
			} else {
				require.Error(t, err)
				require.Nil(t, got)
			}
		})
	}
}

func TestCartService(t *testing.T) {
	client := connectProduct(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	product, err := client.GetProduct(ctx, &pb.GetProductRequest{Sku: 1076963})
	require.NoError(t, err)

	userID := strconv.FormatUint(uint64(newUserID(t)), 10)
	httpClient := newHTTPClient()
	cartURL := "http://localhost:8080/user/" + userID + "/cart"

	response, err := httpClient.Get(cartURL)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, http.StatusNotFound, response.StatusCode)

	postURL := cartURL + "/" + strconv.FormatUint(uint64(product.Sku), 10)
	response, err = httpClient.Post(postURL, "application/json", bytes.NewBufferString(`{"count":2}`))
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, http.StatusOK, response.StatusCode)

	response, err = httpClient.Get(cartURL)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusOK, response.StatusCode)

	var got struct {
		Items []struct {
			Sku   uint32
			Name  string
			Price uint32
			Count uint32
		}
		TotalPrice uint32
	}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&got))
	require.Len(t, got.Items, 1)
	require.Equal(t, product.Sku, got.Items[0].Sku)
	require.Equal(t, product.Name, got.Items[0].Name)
	require.Equal(t, uint32(product.Price), got.Items[0].Price)
	require.Equal(t, uint32(2), got.Items[0].Count)
	require.Equal(t, uint32(product.Price)*2, got.TotalPrice)

	response, err = httpClient.Post(
		"http://localhost:8080/user/"+userID+"/cart/666",
		"application/json",
		bytes.NewBufferString(`{"count":1}`),
	)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, http.StatusPreconditionFailed, response.StatusCode)

	response, err = httpClient.Get("http://localhost:8080/user/0/cart")
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
}
