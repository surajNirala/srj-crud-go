package transforms

import "github.com/surajNirala/srj-crud/models"

type UserResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Image    string `json:"image"`
	Age      int    `json:"age"`
	Status   bool   `json:"status"`
	// CreatedAt
}

func TransformUser(user models.User) UserResponse {
	var Image string
	if user.Image != nil {
		Image = *user.Image
	}
	return UserResponse{
		ID:       user.ID,
		FullName: user.FirstName + "" + user.LastName,
		Email:    user.Email,
		Phone:    user.Phone,
		Image:    Image,
		Age:      int(user.Age),
		Status:   user.Status,
	}
}
