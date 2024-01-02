package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/surajNirala/srj-crud/config"
	"github.com/surajNirala/srj-crud/handlers"
	"github.com/surajNirala/srj-crud/models"
	"gorm.io/gorm"
)

var DB *gorm.DB

// var err error

type Product struct {
	gorm.Model
	Code        string
	Price       string
	Name        string
	Description string
}

func main() {
	config.LoadEnv()
	dbConfig := config.GetDatabaseConfig()
	db, err := config.ConnectDB(dbConfig)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	DB = db
	InitilizedMigration()
	routerInitialized()
}

func InitilizedMigration() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&Product{})
}

func routerInitialized() {
	router := mux.NewRouter()
	router.HandleFunc("/products", GetAllProduct).Methods("GET")
	router.HandleFunc("/products", CreateProduct).Methods("POST")
	router.HandleFunc("/products/{product_id}", GetProduct).Methods("GET")
	router.HandleFunc("/products/{product_id}", UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{product_id}", DeleteProduct).Methods("DELETE")

	/*****************Start User Handler****************/
	router.HandleFunc("/users", handlers.GetAllUser).Methods("GET")
	router.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	// router.HandleFunc("/users/{user_id}", GetUser).Methods("GET")
	// router.HandleFunc("/users/{user_id}", UpdateUser).Methods("PUT")
	// router.HandleFunc("/users/{user_id}", DeleteUser).Methods("DELETE")
	/*****************End User Handler****************/

	log.Fatal(http.ListenAndServe(":4002", router))
}

func GetAllProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var products []Product
	DB.Find(&products)
	json.NewEncoder(w).Encode(products)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var product Product
	var params = mux.Vars(r)
	var productId = params["product_id"]
	// result := DB.Select("name", "code", "price").First(&product, productId)
	result := DB.First(&product, productId)
	if result.Error != nil || result.RowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}
	// fmt.Fprintf(w, "ID from URL: %s", productId)
	json.NewEncoder(w).Encode(product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var product Product
	var params = mux.Vars(r)
	var productId = params["product_id"]
	result := DB.Delete(&product, productId)
	if result.Error != nil || result.RowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: fmt.Sprintf("Product with ID %d deleted successfully", productId),
	})
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var product Product
	var params = mux.Vars(r)
	var productId = params["product_id"]
	result := DB.First(&product, productId)
	if result.Error != nil || result.RowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}
	product.Name = "Srj-product"
	product.Code = "SRJ-COde-007"
	product.Price = "100"
	DB.Save(&product)
	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  "success",
		Message: fmt.Sprintf("Product with ID %d Updated successfully", productId),
	})
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var product Product
	json.NewDecoder(r.Body).Decode(&product)
	// DB.Create(&product)
	result := DB.Create(&product)
	if result.Error != nil {
		log.Printf("Error creating record: %v", result.Error)
		// Handle the error (return, log, etc.)
	}
	json.NewEncoder(w).Encode(product)
}

type APIResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, APIResponse{
		StatusCode: code,
		Status:     "error",
		Message:    message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
