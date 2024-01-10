package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

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
	router.HandleFunc("/", handlers.Index).Methods("GET")

	router.HandleFunc("/products", handlers.GetAllProduct).Methods("GET")
	router.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{product_id}", handlers.GetProduct).Methods("GET")
	router.HandleFunc("/products/{product_id}", handlers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{product_id}", handlers.DeleteProduct).Methods("DELETE")

	/*****************Start User Handler****************/
	router.HandleFunc("/users", handlers.GetAllUser).Methods("GET")
	router.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/users/{user_id}", handlers.GetUser).Methods("GET")
	router.HandleFunc("/users/{user_id}", handlers.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{user_id}", handlers.DeleteUser).Methods("DELETE")
	/*****************End User Handler****************/

	/*****************Start Login Handler****************/
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	/*****************End Login Handler****************/

	/*****************Start user-detail Handler****************/
	router.HandleFunc("/user-info", handlers.UserInfo).Methods("GET")
	router.HandleFunc("/profie-image-update", handlers.ProfileImageUpdate).Methods("POST")
	// router.HandleFunc("/profie-image-update", handleFileUpload).Methods("POST")
	/*****************End user-detail Handler****************/

	log.Fatal(http.ListenAndServe(":4002", router))
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	// Parse the form data, including the uploaded file
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get a reference to the uploaded file
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error getting file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the full path for storing the uploaded file
	uploadDir := "uploaded/files/"
	fullPath := filepath.Join(uploadDir, handler.Filename)

	// Save the file to the full path
	err = saveFile(file, fullPath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "File uploaded successfully!")
}

func saveFile(file multipart.File, fullPath string) error {
	// Open a new file on the server for writing
	outFile, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Copy the uploaded file to the destination file on the server
	_, err = io.Copy(outFile, file)
	return err
}
