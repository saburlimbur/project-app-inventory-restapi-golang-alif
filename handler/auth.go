package handler

import (
	"alfdwirhmn/inventory/dto"
	"alfdwirhmn/inventory/service"
	"alfdwirhmn/inventory/utils"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthHandler struct {
	Service service.AuthService
	Logger  *zap.Logger
}

func NewAuthHandler(service service.AuthService, log *zap.Logger) *AuthHandler {
	return &AuthHandler{
		Service: service,
		Logger:  log,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	// decode body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	// validasi request
	if validationErrors, err := utils.ValidateErrors(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	user, err := h.Service.Register(r.Context(), req)
	if err != nil {
		h.Logger.Error("failed to register user", zap.Error(err))
		utils.JSONError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.JSONSuccess(
		w,
		http.StatusCreated,
		"register success",
		dto.ToUserResponseDTO(user),
	)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	json.NewDecoder(r.Body).Decode(&req)

	ip := r.RemoteAddr
	ua := r.UserAgent()

	resp, err := h.Service.Login(r.Context(), req, ip, ua)
	if err != nil {
		utils.JSONError(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.JSONSuccess(w, http.StatusOK, "login success", resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Context().Value("token").(string)
	token, _ := uuid.Parse(tokenStr)

	h.Service.Logout(r.Context(), token)
	utils.JSONSuccess(w, http.StatusOK, "logout success", nil)
}
