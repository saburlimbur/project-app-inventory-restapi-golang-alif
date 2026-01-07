package handler

import (
	"alfdwirhmn/inventory/dto"
	"alfdwirhmn/inventory/middleware"
	"alfdwirhmn/inventory/model"
	"alfdwirhmn/inventory/service"
	"alfdwirhmn/inventory/utils"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ItemsHandler struct {
	ItemsService service.ItemsService
	Validator    *validator.Validate
	Logger       *zap.Logger

	Config utils.Configuration
}

func NewItemsHandler(service service.ItemsService, validator *validator.Validate, logger *zap.Logger, config utils.Configuration) *ItemsHandler {
	return &ItemsHandler{
		ItemsService: service,
		Validator:    validator,
		Logger:       logger,
		Config:       config,
	}
}

func (h *ItemsHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	items, err := h.ItemsService.Create(r.Context(), user, req)
	if err != nil {
		utils.JSONError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	h.Logger.Info("Items created successfully",
		zap.Int("id", items.ID),
		zap.String("name", items.Name),
	)

	utils.JSONSuccess(w, http.StatusCreated, "Category created successfully",
		dto.ItemResponseDTO{
			ID:           items.ID,
			CategoryID:   items.CategoryID,
			SKU:          items.SKU,
			Name:         items.Name,
			Description:  items.Description,
			Unit:         items.Unit,
			Price:        items.Price,
			Cost:         items.Cost,
			Stock:        items.Stock,
			MinimumStock: items.MinimumStock,
			Weight:       items.Weight,
			Dimensions:   items.Dimensions,
			IsActive:     items.IsActive,
			CreatedBy:    *items.CreatedBy,
			CreatedAt:    items.CreatedAt,
			UpdatedAt:    items.UpdatedAt,
		},
	)
}

func (h *ItemsHandler) Lists(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		utils.JSONError(w, http.StatusBadRequest, "invalid page", nil)
	}

	limit, err := strconv.Atoi(h.Config.Limit)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "invalid limit config", nil)
		return
	}

	items, pagination, err := h.ItemsService.FindAll(page, limit)
	if err != nil {
		h.Logger.Error("failed get items", zap.Error(err))
		utils.JSONError(w, http.StatusInternalServerError, "failed", nil)
		return
	}

	utils.JSONWithPagination(w, http.StatusOK, "succesfully get items data", items, *pagination)
}

func (h *ItemsHandler) DetailById(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "invalid item id", nil)
		return
	}

	res, err := h.ItemsService.FindByID(context.Background(), id, user)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if res == nil {
		utils.JSONError(w, http.StatusNotFound, "item not found", nil)
		return
	}

	utils.JSONSuccess(w, http.StatusOK, "succesfully get item detail", dto.ItemResponseDTO{
		CategoryID:   res.CategoryID,
		RackID:       res.RackID,
		SKU:          res.SKU,
		Name:         res.Name,
		Description:  res.Description,
		Unit:         res.Unit,
		Price:        res.Price,
		Cost:         res.Cost,
		Stock:        res.Stock,
		MinimumStock: res.MinimumStock,
		Weight:       res.Weight,
		Dimensions:   res.Dimensions,
		IsActive:     res.IsActive,
		CreatedBy:    *res.CreatedBy,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	})
}

func (h *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	// ambil user dari context
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// ambil id item dari param
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid item id", nil)
		return
	}

	var req dto.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// validasi
	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	item, err := h.ItemsService.Update(r.Context(), user, id, req)
	if err != nil {
		utils.JSONError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	h.Logger.Info("item updated successfully",
		zap.Int("item_id", item.ID),
		zap.String("updated_by", user.Role),
	)

	utils.JSONSuccess(
		w,
		http.StatusOK,
		"item updated successfully",
		dto.ItemResponseDTO{
			ID:           item.ID,
			CategoryID:   item.CategoryID,
			RackID:       item.RackID,
			SKU:          item.SKU,
			Name:         item.Name,
			Description:  item.Description,
			Unit:         item.Unit,
			Price:        item.Price,
			Cost:         item.Cost,
			Stock:        item.Stock,
			MinimumStock: item.MinimumStock,
			Weight:       item.Weight,
			Dimensions:   item.Dimensions,
			IsActive:     item.IsActive,
			CreatedBy:    *item.CreatedBy,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		},
	)
}

func (h *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid item id", nil)
		return
	}

	err = h.ItemsService.Delete(r.Context(), user, id)
	if err != nil {
		utils.JSONError(w, http.StatusForbidden, err.Error(), nil)
		return
	}

	h.Logger.Info("item deleted successfully",
		zap.Int("category_id", id),
		zap.String("deleted_by", user.Role),
	)

	utils.JSONSuccess(w, http.StatusOK, "item deleted successfully", id)
}
