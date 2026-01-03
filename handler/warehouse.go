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

type WarehouseHandler struct {
	WarehouseService service.WarehouseService
	Validator        *validator.Validate
	Logger           *zap.Logger

	Config utils.Configuration
}

func NewWarehouseHandler(service service.WarehouseService, validator *validator.Validate, logger *zap.Logger, config utils.Configuration) *WarehouseHandler {
	return &WarehouseHandler{
		WarehouseService: service,
		Validator:        validator,
		Logger:           logger,
		Config:           config,
	}
}

func (h *WarehouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.CreateWarehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "validation error", err)
		return
	}

	warehouse, err := h.WarehouseService.Create(r.Context(), user, req)
	if err != nil {
		utils.JSONError(w, http.StatusForbidden, err.Error(), nil)
		return
	}

	h.Logger.Info("warehouse created",
		zap.Int("warehouse_id", warehouse.ID),
		zap.Int("created_by", *warehouse.CreatedBy),
	)

	utils.JSONSuccess(w, http.StatusCreated, "Warehouse created successfully",
		dto.WarehouseResponseDTO{
			ID:         warehouse.ID,
			Code:       warehouse.Code,
			Name:       warehouse.Name,
			Address:    warehouse.Address,
			City:       warehouse.City,
			Province:   warehouse.Province,
			PostalCode: warehouse.PostalCode,
			Phone:      warehouse.Phone,
			IsActive:   warehouse.IsActive,
			CreatedBy:  *warehouse.CreatedBy,
			CreatedAt:  warehouse.CreatedAt,
			UpdatedAt:  warehouse.UpdatedAt,
		},
	)
}

func (h *WarehouseHandler) Lists(w http.ResponseWriter, r *http.Request) {
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

	warehouse, pagination, err := h.WarehouseService.FindAll(page, limit)
	if err != nil {
		h.Logger.Error("failed get warehouse", zap.Error(err))
		utils.JSONError(w, http.StatusInternalServerError, "failed", nil)
		return
	}

	utils.JSONWithPagination(
		w,
		http.StatusOK,
		"succesfully get warehouse data",
		warehouse,
		*pagination,
	)
}

func (h *WarehouseHandler) Update(w http.ResponseWriter, r *http.Request) {
	// user dari context
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid warehouse id", nil)
		return
	}

	var req dto.UpdateWarehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	warehouse, err := h.WarehouseService.Update(r.Context(), user, id, req)
	if err != nil {
		utils.JSONError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	h.Logger.Info("warehouse updated successfully",
		zap.Int("warehouse_id", warehouse.ID),
		zap.String("updated_by", user.Role),
	)

	utils.JSONSuccess(
		w,
		http.StatusOK,
		"warehouse updated successfully",
		dto.WarehouseResponseDTO{
			ID:         warehouse.ID,
			Code:       warehouse.Code,
			Name:       warehouse.Name,
			Address:    warehouse.Address,
			City:       warehouse.City,
			Province:   warehouse.Province,
			PostalCode: warehouse.PostalCode,
			Phone:      warehouse.Phone,
			IsActive:   warehouse.IsActive,
			CreatedBy:  *warehouse.CreatedBy,
			CreatedAt:  warehouse.CreatedAt,
			UpdatedAt:  warehouse.UpdatedAt,
		},
	)
}

func (h *WarehouseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid warehouse id", nil)
		return
	}

	err = h.WarehouseService.Delete(r.Context(), user, id)
	if err != nil {
		utils.JSONError(w, http.StatusForbidden, err.Error(), nil)
		return
	}

	h.Logger.Info("Warehouses deleted successfully",
		zap.Int("warehouse_id", id),
		zap.String("deleted_by", user.Role),
	)

	utils.JSONSuccess(w, http.StatusOK, "warehouses deleted successfully", id)
}
