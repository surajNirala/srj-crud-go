package models

import (
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint    `json:"id"`
	FirstName string  `gorm:"size:255" json:"first_name"`
	LastName  string  `gorm:"size:255" json:"last_name"`
	Email     string  `gorm:"size:255" json:"email"`
	Password  string  `gorm:"size:255" json:"password"`
	Phone     string  `gorm:"size:255" json:"phone"`
	Image     *string `gorm:"size:255" json:"image"`
	Age       uint    `gorm:"default:18" json:"age"`
	Status    bool    `gorm:"default:true" json:"status"`
}

type Claims struct {
	ID        uint    `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Phone     string  `json:"phone"`
	Image     *string `json:"image"`
	Age       uint    `json:"age"`
	Status    bool    `json:"status"`
	jwt.StandardClaims
}
