package transforms

import "github.com/surajNirala/srj-crud/models"

type UserResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
}

func TransformUser(user models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		FullName: user.FirstName + "" + user.LastName,
	}
}
