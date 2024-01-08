package transforms

import "github.com/surajNirala/srj-crud/models"

type ProductResponse struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Code        string       `json:"code"`
	Price       string       `json:"price"`
	Description string       `json:"description"`
	UserID      uint         `json:"user_id"`
	User        UserResponse `gorm:"foreignKey:UserID" json:"user"`
	// CreatedAt
}

func TransformProduct(product models.Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Code:        product.Code,
		Price:       product.Price,
		Description: product.Description,
		UserID:      product.UserID,
		User:        TransformUser(product.User),
	}
}
