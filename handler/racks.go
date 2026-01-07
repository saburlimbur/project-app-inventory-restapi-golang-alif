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

type RacksHandler struct {
	RacksService service.RacksService
	Validator    *validator.Validate
	Logger       *zap.Logger

	Config utils.Configuration
}

func NewRacksHandler(service service.RacksService, validator *validator.Validate, logger *zap.Logger, config utils.Configuration) *RacksHandler {
	return &RacksHandler{
		RacksService: service,
		Validator:    validator,
		Logger:       logger,
		Config:       config,
	}
}

func (h *RacksHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.CreateRackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "validation error", err)
		return
	}

	rack, err := h.RacksService.Create(r.Context(), user, req)
	if err != nil {
		utils.JSONError(w, http.StatusForbidden, err.Error(), nil)
		return
	}

	h.Logger.Info("rack created",
		zap.Int("rack_id", rack.ID),
		zap.Int("created_by", *rack.CreatedBy),
	)

	utils.JSONSuccess(w, http.StatusCreated, "Rack created successfully",
		dto.RackResponseDTO{
			ID:          rack.ID,
			WarehouseID: rack.WarehouseID,
			Code:        rack.Code,
			Name:        rack.Name,
			Location:    rack.Location,
			Capacity:    rack.Capacity,
			Description: rack.Description,
			IsActive:    rack.IsActive,
			CreatedBy:   *rack.CreatedBy,
			CreatedAt:   rack.CreatedAt,
			UpdatedAt:   rack.UpdatedAt,
		},
	)
}

func (h *RacksHandler) Lists(w http.ResponseWriter, r *http.Request) {
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

	racks, pagination, err := h.RacksService.FindAll(page, limit)
	if err != nil {
		h.Logger.Error("failed get racks", zap.Error(err))
		utils.JSONError(w, http.StatusInternalServerError, "failed", nil)
		return
	}

	utils.JSONWithPagination(w, http.StatusOK, "successfully get rack data", racks, *pagination)
}

func (h *RacksHandler) DetailById(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid racks id", nil)
		return
	}

	res, err := h.RacksService.FindById(id, user)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if res == nil {
		utils.JSONError(w, http.StatusNotFound, "racks not found", nil)
		return
	}

	utils.JSONSuccess(w, http.StatusOK, "succesfully get racks detail", dto.RackResponseDTO{
		ID:          res.ID,
		WarehouseID: res.WarehouseID,
		Code:        res.Code,
		Name:        res.Name,
		Location:    res.Location,
		Capacity:    res.Capacity,
		Description: res.Description,
		IsActive:    res.IsActive,
		CreatedBy:   *res.CreatedBy,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	})
}

func (h *RacksHandler) Update(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid rack id", nil)
		return
	}

	var req dto.UpdateRackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	rack, err := h.RacksService.Update(r.Context(), user, id, req)
	if err != nil {
		utils.JSONError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	h.Logger.Info("rack updated successfully",
		zap.Int("rack_id", rack.ID),
		zap.Int("updated_by", user.ID),
	)

	utils.JSONSuccess(w, http.StatusOK, "rack updated successfully",
		dto.RackResponseDTO{
			ID:          rack.ID,
			WarehouseID: rack.WarehouseID,
			Code:        rack.Code,
			Name:        rack.Name,
			Location:    rack.Location,
			Capacity:    rack.Capacity,
			Description: rack.Description,
			IsActive:    rack.IsActive,
			CreatedBy:   *rack.CreatedBy,
			CreatedAt:   rack.CreatedAt,
			UpdatedAt:   rack.UpdatedAt,
		},
	)
}

func (h *RacksHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid rack id", nil)
		return
	}

	if err := h.RacksService.Delete(r.Context(), user, id); err != nil {
		utils.JSONError(w, http.StatusForbidden, err.Error(), nil)
		return
	}

	h.Logger.Info("rack deleted successfully",
		zap.Int("rack_id", id),
		zap.Int("deleted_by", user.ID),
	)

	utils.JSONSuccess(w, http.StatusOK, "rack deleted successfully", id)
}
