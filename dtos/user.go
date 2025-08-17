package dtos

type UserDTO struct {
	Name   string `json:"name"`
	Email  string `json"email"`
	Phone  string `json:"phone"`
	Role   string `json:"role"`
	Status string `json:"status"`
}
