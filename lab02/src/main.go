package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	ProductList        = make(map[uint64]Product)
	CurrentID   uint64 = 1
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome\n"))
	})

	r.Route("/product", func(r chi.Router) {
		r.Get("/{product_id}", GetProductByID)
		r.Put("/{product_id}", UpdateProductByID)
		r.Delete("/{product_id}", DeleteProductByID)
	})
	r.Post("/product", AddProduct)
	r.Get("/products", GetAllProducts)

	http.ListenAndServe(":8080", r)
}
