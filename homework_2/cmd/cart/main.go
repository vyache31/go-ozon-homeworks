package main

import (
	"log"
	"net/http"

	"github.com/vyache31/go-ozon-homeworks/homework_2/internal/app/cart"
)

func main() {
	log.SetPrefix("[CART] ")
	log.Println("Starting HTTP server..")
	mux := http.NewServeMux()
	cartService, err := cart.NewCartService()
	if err != nil {
		log.Fatalf("failed to create CartService: %v", err)
	}
	cartHandler := cart.NewCartHandler(cartService)

	mux.Handle("GET /docs/{$}", serveBytes("text/html; charset=utf-8", docsHTML))
	mux.Handle("GET /swagger.json", serveBytes("application/json", openAPIyaml))

	mux.HandleFunc("/user/{user_id}/cart", cartHandler.GetCart)
	mux.HandleFunc("/user/{user_id}/cart/{sku_id}", cartHandler.PostProduct)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func serveBytes(contentType string, body []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", contentType)
		_, _ = w.Write(body)
	})
}
