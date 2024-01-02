package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/surajNirala/srj-crud/config"
	"github.com/surajNirala/srj-crud/models"
	"gorm.io/gorm"
)

var DB *gorm.DB

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

type APIResponse struct {
	StatusCode int         //`json:"code"`
	Status     string      //`json:"status"`
	Message    string      //`json:"message"`
	Data       interface{} `json:"Data,omitempty"`
}

func respondWithError(w http.ResponseWriter, code int, message string, payload interface{}) {
	respondWithJSON(w, code, APIResponse{
		StatusCode: code,
		Status:     "error",
		Message:    message,
	})
}

func respondWithSuccess(w http.ResponseWriter, code int, message string, payload interface{}) {
	respondWithJSON(w, code, APIResponse{
		StatusCode: code,
		Status:     "success",
		Message:    message,
		Data:       payload,
	})
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	config.LoadEnv()
	dbConfig := config.GetDatabaseConfig()
	db, err := config.ConnectDB(dbConfig)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	DB = db
	var users []models.User
	DB.Find(&users)
	// json.NewEncoder(w).Encode(users)
	respondWithSuccess(w, http.StatusOK, "Fetched User List.", &users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	config.LoadEnv()
	dbConfig := config.GetDatabaseConfig()
	db, err := config.ConnectDB(dbConfig)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	DB = db
	w.Header().Set("Content-type", "application/json")
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	result := DB.Create(&user)
	if result.Error != nil {
		log.Printf("Error creating record: %v", result.Error)
		respondWithError(w, http.StatusInternalServerError, "Internal Server error.", nil)
		return
		// Handle the error (return, log, etc.)
	}
	respondWithSuccess(w, http.StatusOK, "User created successfully.", &user)
}
