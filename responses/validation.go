package responses

import (
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// ValidateEmail checks if the email is valid.
func ValidateEmail(email string) bool {
	// You can use a more comprehensive email validation library if needed.
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword checks if the password meets certain criteria.
func ValidatePassword(password string) bool {
	// You can customize the password requirements as needed.
	return len(password) >= 8
}

// HashPassword hashes the given password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func WrongPassword(dbPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func ExtractToken(r *http.Request) string {
	// Extract token from Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return ""
	}
	return tokenString[7:]
}
