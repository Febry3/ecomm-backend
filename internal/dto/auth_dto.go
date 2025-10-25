package dto

type RegisterRequest struct {
	Username    string `json:"username" validate:"required"`
	FirstName   string `json:"first_name" validate:"min=3,max=30"`
	LastName    string `json:"last_name" validate:"min=3,max=30"`
	PhoneNumber string `json:"phone_number" validate:"required,number,min=8,max=12"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
}

type RegisterResponse struct {
	ID          int64  `json:"id" `
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	ID          int64  `json:"id" `
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}
