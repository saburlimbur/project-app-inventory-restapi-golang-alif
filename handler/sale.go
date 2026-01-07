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

type SaleHandler struct {
	SaleService service.SaleService
	Validator   *validator.Validate
	Logger      *zap.Logger

	Config utils.Configuration
}

func NewSaleHandler(service service.SaleService, validator *validator.Validate, logger *zap.Logger, config utils.Configuration) *SaleHandler {
	return &SaleHandler{
		SaleService: service,
		Validator:   validator,
		Logger:      logger,
		Config:      config,
	}
}

func (h *SaleHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*model.User)
	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.CreateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	sale, err := h.SaleService.Create(r.Context(), user, req)
	if err != nil {
		utils.JSONError(w, http.StatusConflict, err.Error(), nil)
		return
	}

	h.Logger.Info("Sale created successfully",
		zap.Int("id", sale.ID),
		zap.String("name", *sale.CustomerName),
	)

	utils.JSONSuccess(w, http.StatusCreated, "Sale created succesfully",
		dto.SaleResponseDTO{
			ID:            sale.ID,
			InvoiceNumber: sale.InvoiceNumber,
			CustomerName:  sale.CustomerName,
			CustomerPhone: sale.CustomerPhone,
			CustomerEmail: sale.CustomerEmail,
			SaleDate:      sale.SaleDate,
			TotalAmount:   sale.TotalAmount,
			Discount:      sale.Discount,
			Tax:           sale.Tax,
			GrandTotal:    sale.GrandTotal,
			PaymentMethod: sale.PaymentMethod,
			PaymentStatus: sale.PaymentStatus,
			Notes:         sale.Notes,
			CreatedBy:     *sale.CreatedBy,
			CreatedAt:     sale.CreatedAt,
			UpdatedAt:     sale.UpdatedAt,
		})
}

func (h *SaleHandler) Lists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	sale, pagination, err := h.SaleService.FindAll(ctx, page, limit)
	if err != nil {
		h.Logger.Error("failed get sale", zap.Error(err))
		utils.JSONError(w, http.StatusInternalServerError, "failed", nil)
		return
	}

	utils.JSONWithPagination(w, http.StatusOK, "succesfully get sale data", sale, *pagination)
}

func (h *SaleHandler) FindById(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserContextKey).(*model.User)

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	sale, err := h.SaleService.Detail(r.Context(), user, id)
	if err != nil {
		utils.JSONError(w, http.StatusForbidden, err.Error(), nil)
		return
	}

	utils.JSONSuccess(w, http.StatusOK, "success", sale)
}

func (h *SaleHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserContextKey).(*model.User)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	var req dto.UpdateSaleRequest
	json.NewDecoder(r.Body).Decode(&req)

	sale := &model.Sale{
		ID:            id,
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		CustomerEmail: req.CustomerEmail,
		Notes:         req.Notes,
	}

	if err := h.SaleService.Update(r.Context(), user, sale); err != nil {
		utils.JSONError(w, http.StatusForbidden, err.Error(), nil)
		return
	}

	utils.JSONSuccess(w, http.StatusOK, "sale updated", nil)
}

func (h *SaleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserContextKey).(*model.User)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	if err := h.SaleService.Delete(r.Context(), user, id); err != nil {
		utils.JSONError(w, http.StatusForbidden, err.Error(), nil)
		return
	}

	utils.JSONSuccess(w, http.StatusOK, "sale deleted", nil)
}
