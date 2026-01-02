package handler

import (
	"alfdwirhmn/inventory/dto"
	appMiddleware "alfdwirhmn/inventory/middleware"
	"alfdwirhmn/inventory/model"
	"alfdwirhmn/inventory/service"
	"alfdwirhmn/inventory/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type UserHandler struct {
	UserService service.UserService
	Validator   *validator.Validate
	Logger      *zap.Logger

	Config utils.Configuration
}

func NewUserHandler(service service.UserService, validator *validator.Validate, logger *zap.Logger, config utils.Configuration) *UserHandler {
	return &UserHandler{
		UserService: service,
		Validator:   validator,
		Logger:      logger,
		Config:      config,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	// ambil current user dari context yang di set di auth midlleware
	currentUser, ok := r.Context().
		Value(appMiddleware.UserContextKey).(*model.User)

	if !ok {
		utils.JSONError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// decode req body
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// validasi
	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	user, err := h.UserService.Create(r.Context(), currentUser, req)
	if err != nil {
		utils.JSONError(w, http.StatusForbidden, err.Error(), nil)
		return
	}

	h.Logger.Info("User created successfully",
		zap.Int("user_id", user.ID),
		zap.String("username", user.Username),
		zap.String("role", user.Role),
		zap.Int("created_by", user.CreatedAt.Day()),
	)

	// jadi dto tdk perlu memasukan semua field, karna di dlm ToUserResponseDTO sudah di define user response nya
	utils.JSONSuccess(w, http.StatusCreated, "User created successfully", dto.ToUserResponseDTO(user))
}

func (h *UserHandler) Lists(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		utils.JSONError(w, http.StatusBadRequest, "Invalid page", nil)
		return
	}

	limit, err := strconv.Atoi(h.Config.Limit)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Invalid limit config", nil)
		return
	}

	users, pagination, err := h.UserService.FindAll(page, limit)
	if err != nil {
		h.Logger.Error("failed get users", zap.Error(err))
		utils.JSONError(w, http.StatusInternalServerError, "failed", nil)
		return
	}

	utils.JSONWithPagination(w, http.StatusOK, "success get data", users, *pagination)
}
