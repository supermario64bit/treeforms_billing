package dtos

type SignupDTO struct {
	Name            string `json:"name"`
	Email           string `json"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RotateRefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
	UserID       uint   `json:"user_id"`
}
