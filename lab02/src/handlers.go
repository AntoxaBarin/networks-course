package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func AddProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	product.ID = CurrentID
	ProductList[CurrentID] = product
	addedProductJson, err := json.Marshal(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(addedProductJson)
	CurrentID += 1
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if _, ok := ProductList[uint64(id)]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resultJson, err := json.Marshal(ProductList[uint64(id)])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(resultJson)
}

func UpdateProductByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	product, ok := ProductList[uint64(id)]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var update Product
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if update.Name != product.Name && update.Name != "" {
		product.Name = update.Name
	}
	if update.Description != product.Description && update.Description != "" {
		product.Description = update.Description
	}
	ProductList[uint64(id)] = product
	updatedJson, err := json.Marshal(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(updatedJson)
}

func DeleteProductByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	product, ok := ProductList[uint64(id)]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	delete(ProductList, uint64(id))
	deletedProductJson, err := json.Marshal(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(deletedProductJson)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	result := make([]Product, 0, len(ProductList))
	for _, product := range ProductList {
		result = append(result, product)
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(resultJson)
}

func GetImageHandler(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	product, ok := ProductList[uint64(id)]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	imagePath := product.Icon
	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, imagePath)
}

func PostImageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	product, ok := ProductList[uint64(id)]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("icon")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	imagePath := "../assets/" + string(productID) + ".png"
	err = os.WriteFile(imagePath, fileBytes, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	product.Icon = imagePath
	ProductList[uint64(id)] = product
	w.WriteHeader(http.StatusOK)
}
