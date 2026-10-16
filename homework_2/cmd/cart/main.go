package main

import (
	"log"
	"net/http"

	"github.com/vyache31/go-ozon-homeworks/homework_2/internal/app/cart"
)

func main() {
	log.Println("Starting Cart service..")
	mux := http.NewServeMux()
	cartService, err := cart.NewCartService()
	if err != nil {
		log.Fatalf("failed to create CartService: %v", err)
	}
	cartHandler := cart.NewCartHandler(cartService)

	mux.HandleFunc("/user/{user_id}/cart", cartHandler.GetCart)
	mux.HandleFunc("/user/{user_id}/cart/{sku_id}", cartHandler.PostProduct)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
