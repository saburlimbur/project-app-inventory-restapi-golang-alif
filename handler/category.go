package handler

import (
	"alfdwirhmn/inventory/dto"
	"alfdwirhmn/inventory/middleware"
	"alfdwirhmn/inventory/model"
	"alfdwirhmn/inventory/service"
	"alfdwirhmn/inventory/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	CategoryService service.CategoryService
	Validator       *validator.Validate
	Logger          *zap.Logger

	Config utils.Configuration
}

func NewCategoryHandler(service service.CategoryService, validator *validator.Validate, logger *zap.Logger, config utils.Configuration) *CategoryHandler {
	return &CategoryHandler{
		CategoryService: service,
		Validator:       validator,
		Logger:          logger,
		Config:          config,
	}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	category, err := h.CategoryService.Create(r.Context(), user, req)
	if err != nil {
		utils.JSONError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	h.Logger.Info("Category created successfully",
		zap.String("name", category.Name),
		zap.String("code", category.Code),
	)

	utils.JSONSuccess(w, http.StatusCreated, "Category created successfully",
		dto.CategoryResponseDTO{
			ID:          category.ID,
			Code:        category.Code,
			Name:        category.Name,
			Description: category.Description,
			IsActive:    category.IsActive,
			CreatedBy:   *category.CreatedBy,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		},
	)
}

func (h *CategoryHandler) Lists(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		utils.JSONError(w, http.StatusBadRequest, "invalid page", nil)
		return
	}

	limit, err := strconv.Atoi(h.Config.Limit)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "invalid limit config", nil)
		return
	}

	category, pagination, err := h.CategoryService.FindAll(page, limit)
	if err != nil {
		h.Logger.Error("failed get category", zap.Error(err))
		utils.JSONError(w, http.StatusInternalServerError, "failed", nil)
		return
	}

	utils.JSONWithPagination(w, http.StatusOK, "succesfully get category data", category, *pagination)
}

func (h *CategoryHandler) DetailById(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid category id", nil)
		return
	}

	res, err := h.CategoryService.FindById(id, user)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if res == nil {
		utils.JSONError(w, http.StatusNotFound, "category not found", nil)
		return
	}

	utils.JSONSuccess(w, http.StatusOK, "successfully get category detail", dto.CategoryResponseDTO{
		ID:          res.ID,
		Code:        res.Code,
		Name:        res.Name,
		Description: res.Description,
		IsActive:    res.IsActive,
		CreatedBy:   *res.CreatedBy,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	})
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	// user dari context
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid category id", nil)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	category, err := h.CategoryService.Update(r.Context(), user, id, req)
	if err != nil {
		utils.JSONError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	h.Logger.Info("Category updated successfully",
		zap.Int("category_id", category.ID),
		zap.String("updated_by", user.Role),
	)

	utils.JSONSuccess(w, http.StatusOK, "category updated successfully",
		dto.CategoryResponseDTO{
			ID:          category.ID,
			Code:        category.Code,
			Name:        category.Name,
			Description: category.Description,
			IsActive:    category.IsActive,
			CreatedBy:   *category.CreatedBy,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		},
	)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid category id", nil)
		return
	}

	err = h.CategoryService.Delete(r.Context(), user, id)
	if err != nil {
		utils.JSONError(w, http.StatusForbidden, err.Error(), nil)
		return
	}

	h.Logger.Info("Category deleted successfully",
		zap.Int("category_id", id),
		zap.String("deleted_by", user.Role),
	)

	utils.JSONSuccess(w, http.StatusOK, "category deleted successfully", id)
}
