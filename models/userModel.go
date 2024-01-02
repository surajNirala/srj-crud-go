package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string  `gorm:"size:255"`
	LastName  string  `gorm:"size:255"`
	Email     string  `gorm:"size:255"`
	Phone     string  `gorm:"size:255"`
	Image     *string `gorm:"size:255"`
	// FullName  string  `gorm:"size:255->;type:GENERATED ALWAYS AS (concat(first_name,' ',last_name));default:(-)"`
	Age    uint `gorm:"default:18"`
	Status bool `gorm:"default:true"`
}
