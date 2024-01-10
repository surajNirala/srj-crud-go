package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/surajNirala/srj-crud/config"
	"github.com/surajNirala/srj-crud/middleware"
	"github.com/surajNirala/srj-crud/models"
	"github.com/surajNirala/srj-crud/responses"
)

var jwtKey = []byte("your-secret-key")

func Login(w http.ResponseWriter, r *http.Request) {
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
	result := DB.Where("email =?", email).First(&user)
	if result.Error != nil {
		responses.ResponseError(w, http.StatusNotFound, "User is not found.", nil)
		return
		// Handle the error (return, log, etc.)
	}
	if !responses.WrongPassword(user.Password, password) {
		responses.ResponseError(w, http.StatusBadRequest, "Invalid Email/Password.", nil)
		return
	}
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours
	claims := &models.Claims{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Phone:     user.Phone,
		Image:     user.Image,
		Age:       user.Age,
		Status:    user.Status,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		responses.ResponseError(w, http.StatusInternalServerError, "Error Creating token.", nil)
		return
	}

	payloadToken := map[string]string{
		"token": signedToken,
	}
	responses.ResponseSuccess(w, http.StatusOK, "Token Generated successfully.", payloadToken)
}

func UserInfo(w http.ResponseWriter, r *http.Request) {
	tokenString := responses.ExtractToken(r)
	// claims := &Claims{}
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		responses.ResponseError(w, http.StatusUnauthorized, "Unauthorized.", nil)
		return
	}
	if !token.Valid {
		responses.ResponseError(w, http.StatusUnauthorized, "Token is not valid.", nil)
		return
	}
	responses.ResponseSuccess(w, http.StatusOK, "User Information from token.", claims)
	return
}

func ProfileImageUpdate(w http.ResponseWriter, r *http.Request) {
	var user models.User
	status, message := middleware.Islogin(w, r)
	if status == 401 {
		responses.ResponseError(w, status, message, nil)
		return
	}

	// DB := config.DB
	err := r.ParseMultipartForm(10 << 20) // Set a limit for the uploaded file size
	if err != nil {
		responses.ResponseError(w, http.StatusBadRequest, "Error parsing form", nil)
		return
	}
	// Get a reference to the uploaded file
	file, handler, err := r.FormFile("file")
	if err != nil {
		responses.ResponseError(w, http.StatusBadRequest, "Error getting file", nil)
		return
	}
	defer file.Close()
	user_id := r.FormValue("user_id")
	if user_id == "" {
		responses.ResponseError(w, http.StatusBadRequest, "User ID not found", nil)
		return
	}
	DB := config.DB
	result := DB.Where("id = ?", user_id).First(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		responses.ResponseError(w, http.StatusNotFound, "User not found", nil)
		return
	}
	uploadDir := "uploaded/files/"
	filename := generateUniqueFilename(handler.Filename)
	fullPath := filepath.Join(uploadDir, filename)
	// Save the file to the full path
	err = saveFile(file, fullPath)
	if err != nil {
		responses.ResponseError(w, http.StatusBadRequest, "Error saving file", nil)
		return
	}
	user.Image = &filename
	DB.Save(&user)
	// DB.Exec("INSERT INTO files (image) VALUES (?)", filename)
	responses.ResponseSuccess(w, http.StatusOK, "File uploaded successfully!", nil)
	return
}

func generateUniqueFilename(originalFilename string) string {
	// Get the current timestamp
	currentTime := time.Now()

	// Generate a random identifier (e.g., random number or hash)
	randomIdentifier := rand.Intn(1000) // Adjust the range as needed

	// Create a unique filename using a combination of timestamp and random identifier
	uniqueFilename := fmt.Sprintf("%d_%d_%s", currentTime.Unix(), randomIdentifier, originalFilename)

	// Optionally, sanitize the filename (replace spaces, special characters, etc.)
	// uniqueFilename = sanitizeFilename(uniqueFilename)

	return uniqueFilename
}

func saveFile(file multipart.File, filepath string) error {
	// Open a new file on the server for writing
	outFile, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer outFile.Close()
	// Copy the uploaded file to the destination file on the server
	_, err = io.Copy(outFile, file)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// fmt.Println(err.Error())
	return err
}

// func storeFileInfo(filename, filepath string) error {
// 	// Insert file information into the database
// 	_, err := db.Exec("INSERT INTO files (filename, filepath) VALUES (?, ?)", filename, filepath)
// 	return err
// }
