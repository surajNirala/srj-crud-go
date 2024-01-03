package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/surajNirala/srj-crud/config"
	"github.com/surajNirala/srj-crud/handlers"
	"github.com/surajNirala/srj-crud/models"
)

func main() {
	config.LoadEnv()
	dbConfig := config.GetDatabaseConfig()
	config.InitDB(dbConfig)

	InitilizedMigration()
	routerInitialized()
}

func InitilizedMigration() {
	DB := config.DB
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Product{})
}

func routerInitialized() {
	router := mux.NewRouter()
	router.HandleFunc("/products", handlers.GetAllProduct).Methods("GET")
	router.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{product_id}", handlers.GetProduct).Methods("GET")
	router.HandleFunc("/products/{product_id}", handlers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{product_id}", handlers.DeleteProduct).Methods("DELETE")

	/*****************Start User Handler****************/
	router.HandleFunc("/users", handlers.GetAllUser).Methods("GET")
	router.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	// router.HandleFunc("/users/{user_id}", GetUser).Methods("GET")
	// router.HandleFunc("/users/{user_id}", UpdateUser).Methods("PUT")
	// router.HandleFunc("/users/{user_id}", DeleteUser).Methods("DELETE")
	/*****************End User Handler****************/
	log.Fatal(http.ListenAndServe(":4002", router))
}
