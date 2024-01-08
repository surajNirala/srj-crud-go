package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/surajNirala/srj-crud/config"
	"github.com/surajNirala/srj-crud/models"
	"github.com/surajNirala/srj-crud/responses"
	"github.com/surajNirala/srj-crud/transforms"
)

func GetAllProduct(w http.ResponseWriter, r *http.Request) {
	DB := config.DB
	var products []models.Product
	DB.Preload("User").Find(&products)

	var transformProducts []transforms.ProductResponse
	for _, product := range products {
		transformProducts = append(transformProducts, transforms.TransformProduct(product))
	}
	responses.ResponseSuccess(w, http.StatusOK, "Fetched Product List.", transformProducts)
	// json.NewEncoder(w).Encode(products)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	DB := config.DB
	w.Header().Set("Content-type", "application/json")
	var product models.Product
	var params = mux.Vars(r)
	var productId = params["product_id"]
	// result := DB.Select("name", "code", "price").First(&product, productId)
	result := DB.Preload("User").First(&product, productId)
	if result.Error != nil || result.RowsAffected == 0 {
		responses.ResponseError(w, http.StatusNotFound, "Product not found", nil)
		return
	}

	transforms := transforms.TransformProduct(product)
	responses.ResponseSuccess(w, http.StatusOK, "Feteched Product Detail.", transforms)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	DB := config.DB
	w.Header().Set("Content-type", "application/json")
	var product models.Product
	var params = mux.Vars(r)
	var productId = params["product_id"]
	result := DB.Delete(&product, productId)
	if result.Error != nil || result.RowsAffected == 0 {
		responses.ResponseError(w, http.StatusNotFound, "Product not found", nil)
		return
	}

	responses.ResponseSuccess(w, http.StatusOK, "Product deleted successfully.", nil)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	DB := config.DB
	w.Header().Set("Content-type", "application/json")
	var product models.Product
	var params = mux.Vars(r)
	var productId = params["product_id"]
	result := DB.First(&product, productId)
	if result.Error != nil || result.RowsAffected == 0 {
		responses.ResponseError(w, http.StatusNotFound, "Product not found", nil)
		return
	}
	product.Name = "Srj-product"
	product.Code = "SRJ-COde-007"
	product.Price = "100"
	DB.Save(&product)
	responses.ResponseSuccess(w, http.StatusOK, "Product updated successfully.", product)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	DB := config.DB
	w.Header().Set("Content-type", "application/json")
	var product models.Product
	json.NewDecoder(r.Body).Decode(&product)
	// DB.Create(&product)
	result := DB.Create(&product)
	if result.Error != nil {
		log.Printf("Error creating record: %v", result.Error)
	}
	responses.ResponseSuccess(w, http.StatusOK, "Product created successfully.", product)
}
