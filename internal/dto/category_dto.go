package dto

import "github.com/febry3/gamingin/internal/entity"

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func ToCategoryResponse(category []entity.Category) []CategoryResponse {
	var response []CategoryResponse
	for _, v := range category {
		response = append(response, CategoryResponse{
			ID:   v.ID,
			Name: v.Name,
			Slug: v.Slug,
		})
	}
	return response
}
