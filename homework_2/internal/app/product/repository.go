package product

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/vyache31/go-ozon-homeworks/homework_2/internal/domain"
)

type JSONProductRepository struct {
	Products map[domain.ProductSKU]domain.Product
	path     string
	mu       sync.RWMutex
}

func NewJSONProductRepository(path string) *JSONProductRepository {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &JSONProductRepository{
			Products: make(map[domain.ProductSKU]domain.Product),
			path:     path,
			mu:       sync.RWMutex{},
		}
	} else if err == nil {
		repo := &JSONProductRepository{
			mu: sync.RWMutex{},
		}
		if err := uploadJSON(path, repo); err != nil {
			log.Fatalf("Failed to upload json file: %v", err)
		}
		log.Print("Product data uploaded")
		repo.path = path
		return repo
	} else {
		log.Fatalf("Failed to init JSONProductRepository: %v", err)
		return nil
	}
}

func (r *JSONProductRepository) Get(sku domain.ProductSKU) (domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.Products[sku]
	if !exists {
		return domain.Product{}, domain.ErrorProductNotFound
	}

	return product, nil
}

func uploadJSON(path string, repo *JSONProductRepository) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read JSON: %w", err)
	}
	err = json.Unmarshal(file, &repo.Products)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}
