package dto

import (
	"time"

	"github.com/febry3/gamingin/internal/entity"
)

type AddressRequest struct {
	StreetAddress string `validate:"required" json:"street_address"`
	RT            string `validate:"required" json:"rt"`
	RW            string `validate:"required" json:"rw"`
	Village       string `validate:"required" json:"village"`
	District      string `validate:"required" json:"district"`
	City          string `validate:"required" json:"city"`
	Province      string `validate:"required" json:"province"`
	PostalCode    string `validate:"required" json:"postal_code"`
	Notes         string `validate:"required" json:"notes"`
}

type AddressResponse struct {
	AddressID     string    `json:"address_id"`
	UserID        int64     `json:"user_id"`
	StreetAddress string    `json:"street_address"`
	RT            string    `json:"rt"`
	RW            string    `json:"rw"`
	Village       string    `json:"village"`
	District      string    `json:"district"`
	City          string    `json:"city"`
	Province      string    `json:"province"`
	PostalCode    string    `json:"postal_code"`
	Notes         string    `json:"notes"`
	IsDefault     bool      `json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (req *AddressRequest) UpdateEntity(a *entity.Address) {
	if req.StreetAddress != "" {
		a.StreetAddress = req.StreetAddress
	}
	if req.RT != "" {
		a.RT = req.RT
	}
	if req.RW != "" {
		a.RW = req.RW
	}
	if req.Village != "" {
		a.Village = req.Village
	}
	if req.District != "" {
		a.District = req.District
	}
	if req.City != "" {
		a.City = req.City
	}
	if req.Province != "" {
		a.Province = req.Province
	}
	if req.PostalCode != "" {
		a.PostalCode = req.PostalCode
	}
	if req.Notes != "" {
		a.Notes = req.Notes
	}
}
