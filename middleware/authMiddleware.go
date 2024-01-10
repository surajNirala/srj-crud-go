package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/surajNirala/srj-crud/models"
)

var jwtKey = []byte("your-secret-key")

func Islogin(w http.ResponseWriter, r *http.Request) (int, string) {
	tokenString := extractToken(r)
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return http.StatusUnauthorized, "Unauthorized."
	}
	if !token.Valid {
		return http.StatusUnauthorized, "Invalid Token."
	}
	return http.StatusOK, "Valid User"
}

func extractToken(r *http.Request) string {
	// Extract token from Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return ""
	}
	return tokenString[7:]
}
