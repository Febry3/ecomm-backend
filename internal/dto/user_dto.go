package dto

import "github.com/febry3/gamingin/internal/entity"

type UserRequest struct {
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	ProfileUrl  string `json:"profile_url"`
	UserID      int64  `json:"user_id"`
}

type UserResponse struct {
	ID          int64  `json:"id" `
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	AccessToken string `json:"access_token"`
	ProfileUrl  string `json:"profile_url"`
}

func (req *UserRequest) UpdateEntity(u *entity.User) {
	if req.Username != "" {
		u.Username = req.Username
	}
	if req.FirstName != "" {
		u.FirstName = req.FirstName
	}
	if req.LastName != "" {
		u.LastName = req.LastName
	}
	if req.PhoneNumber != "" {
		u.PhoneNumber = req.PhoneNumber
	}
	if req.ProfileUrl != "" {
		u.ProfileUrl = req.ProfileUrl
	}
}
