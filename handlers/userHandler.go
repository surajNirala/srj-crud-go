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

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	DB := config.DB
	var users []models.User
	DB.Find(&users)
	// json.NewEncoder(w).Encode(users)
	var transformedUsers []transforms.UserResponse
	// var simpleTranforms []map[string]interface{}
	for _, user := range users {
		// simpleTranform := map[string]interface{}{
		// 	"userName": user.FirstName + user.LastName,
		// }
		// simpleTranforms = append(simpleTranforms, simpleTranform)
		transformedUsers = append(transformedUsers, transforms.TransformUser(user))
	}
	// responses.ResponseSuccess(w, http.StatusOK, "Fetched User List.", simpleTranforms)
	responses.ResponseSuccess(w, http.StatusOK, "Fetched User List.", transformedUsers)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	DB := config.DB
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		responses.ResponseError(w, http.StatusInternalServerError, "Invalid Data Request.", nil)
	}
	email := user.Email
	password := user.Password
	if !responses.ValidateEmail(email) {
		responses.ResponseError(w, http.StatusBadRequest, "Invalid Email Address", nil)
		return
	}
	if !responses.ValidatePassword(password) {
		responses.ResponseError(w, http.StatusBadRequest, "Password must be least 8 characters.", nil)
		return
	}
	hashedPassword, err := responses.HashPassword(user.Password)
	if err != nil {
		responses.ResponseError(w, http.StatusInternalServerError, "Error hasing password", nil)
		return
	}
	user.Password = hashedPassword
	result := DB.Create(&user)
	if result.Error != nil {
		log.Printf("Error creating record: %v", result.Error)
		responses.ResponseError(w, http.StatusInternalServerError, "Internal Server error.", nil)
		return
		// Handle the error (return, log, etc.)
	}
	responses.ResponseSuccess(w, http.StatusOK, "User created successfully.", &user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	DB := config.DB
	var user models.User
	params := mux.Vars(r)
	user_id := params["user_id"]
	userDetail := DB.Find(&user, user_id)
	if userDetail.Error != nil || userDetail.RowsAffected == 0 {
		responses.ResponseError(w, http.StatusNotFound, "User not found", nil)
		return
	}
	data := transforms.TransformUser(user)
	responses.ResponseSuccess(w, http.StatusOK, "Fected User Details", data)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	DB := config.DB
	var user models.User
	var updatedUser models.User
	params := mux.Vars(r)
	userID := params["user_id"]
	userDetail := DB.Find(&user, userID)
	if userDetail.Error != nil || userDetail.RowsAffected == 0 {
		responses.ResponseError(w, http.StatusNotFound, "User not found", nil)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		responses.ResponseError(w, http.StatusBadRequest, "Invalid request data", nil)
		return
	}
	var existingUser models.User
	if DB.Where("email = ? AND id != ?", updatedUser.Email, userID).First(&existingUser).RowsAffected > 0 {
		responses.ResponseError(w, http.StatusBadRequest, "Email already exists", nil)
		return
	}
	// Update user details
	user.FirstName = updatedUser.FirstName
	user.LastName = updatedUser.LastName
	user.Email = updatedUser.Email
	// Save the updated user to the database
	DB.Save(&user)
	data := transforms.TransformUser(user)
	responses.ResponseSuccess(w, http.StatusOK, "User Detail Updated Successfully.", data)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	DB := config.DB
	var user models.User
	params := mux.Vars(r)
	user_id := params["user_id"]
	userDetail := DB.Delete(&user, user_id)
	if userDetail.Error != nil || userDetail.RowsAffected == 0 {
		responses.ResponseError(w, http.StatusNotFound, "User not found", nil)
		return
	}
	responses.ResponseSuccess(w, http.StatusOK, "User Deleted Successfully.", nil)
}
