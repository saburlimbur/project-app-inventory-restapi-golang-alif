package dto

import "time"

type CategoryResponseDTO struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedBy   int    `json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Code        string `json:"code" validate:"required,max=20"`
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`

	// CreatedBy *int  `json:"created_by,omitempty" db:"created_by"`
}

// func ToCategoryResponseDTO(ctg *model.Category) *CategoryResponseDTO {
// 	return &CategoryResponseDTO{
// 		ID:          ctg.ID,
// 		Code:        ctg.Code,
// 		Name:        ctg.Name,
// 		Description: ctg.Description,
// 		IsActive:    ctg.IsActive,

// 		CreatedAt: ctg.CreatedAt.Format("2006-01-02 15:04:05"),
// 		UpdatedAt: ctg.UpdatedAt.Format("2006-01-02 15:04:05"),
// 	}
// }
