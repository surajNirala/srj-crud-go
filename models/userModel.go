package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string  `gorm:"size:255" json:"first_name"`
	LastName  string  `gorm:"size:255" json:"last_name"`
	Email     string  `gorm:"size:255" json:"email"`
	Phone     string  `gorm:"size:255" json:"phone"`
	Image     *string `gorm:"size:255" json:"image"`
	// FullName  string  `gorm:"size:255->;type:GENERATED ALWAYS AS (concat(first_name,' ',last_name));default:(-)"`
	Age    uint `gorm:"default:18" json:"age"`
	Status bool `gorm:"default:true" json:"status"`
}
