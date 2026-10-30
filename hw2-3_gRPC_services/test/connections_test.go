//go:build api

package integration_test

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	pb "github.com/vyache31/go-ozon-homeworks/hw2-3_gRPC_services/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connectProduct(t *testing.T) pb.ProductClient {
	t.Helper()

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, conn.Close())
	})

	return pb.NewProductClient(conn)
}

func newUserID(t *testing.T) uint32 {
	t.Helper()

	maxID := new(big.Int).SetUint64(uint64(^uint32(0)))
	id, err := rand.Int(rand.Reader, maxID)
	require.NoError(t, err)

	return uint32(id.Uint64() + 1)
}

func newHTTPClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}
