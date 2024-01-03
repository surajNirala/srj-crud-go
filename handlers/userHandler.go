package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/surajNirala/srj-crud/config"
	"github.com/surajNirala/srj-crud/models"
	"github.com/surajNirala/srj-crud/responses"
	"github.com/surajNirala/srj-crud/transforms"
)

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	DB := config.DB
	var users []models.User
	DB.Find(&users)
	// json.NewEncoder(w).Encode(users)
	var transformedUsers []transforms.UserResponse
	for _, user := range users {
		transformedUsers = append(transformedUsers, transforms.TransformUser(user))
	}
	responses.ResponseSuccess(w, http.StatusOK, "Fetched User List.", transformedUsers)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	DB := config.DB
	w.Header().Set("Content-type", "application/json")
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	result := DB.Create(&user)
	if result.Error != nil {
		log.Printf("Error creating record: %v", result.Error)
		responses.ResponseError(w, http.StatusInternalServerError, "Internal Server error.", nil)
		return
		// Handle the error (return, log, etc.)
	}
	responses.ResponseSuccess(w, http.StatusOK, "User created successfully.", &user)
}
