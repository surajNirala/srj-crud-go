package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string `json:"name"`
	Code        string `json:"code"`
	Price       string `json:"price"`
	Description string `json:"description"`
}
