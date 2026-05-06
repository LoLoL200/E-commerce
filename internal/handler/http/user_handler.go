package http

import (
	"ecommers/internal/service/auth"
	userService "ecommers/internal/service/auth"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService userService.UserService
	authService auth.AuthService
}

func NewUserHandler(srvc userService.UserService, authSrv auth.AuthService) *UserHandler {
	return &UserHandler{
		authService: authSrv,
		userService: srvc,
	}
}

// Router for register
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/profile", func(r chi.Router) {
		r.Get("/", h.GetProfile)
		r.Put("/", h.UpdateProfile)
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current authenticated user profile data
// @Tags profile
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.User
// @Failure 401 {object} ErrorResponse "User not authenticated"
// @Failure 404 {object} ErrorResponse "User not found"
// @Router /api/v1/profile [get]
// GET
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not authentificated")
		return
	}
	user, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}
	respondJSON(w, http.StatusOK, user)
}

// UdateProfile godoc
// @Summary Update user profile
// @Description Update current authenticated user profile details
// @Tags profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body userService.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 401 {object} ErrorResponse "User not authenticated"
// @Failure 404 {object} ErrorResponse "Failed to update profile"
// @Router /api/v1/profile [put]
// Update
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "user not authentificated")
		return
	}
	var req userService.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userService.UpdateProfile(r.Context(), userID, req, false)
	if err != nil {
		respondError(w, http.StatusNotFound, "Failed to update profile")
		return
	}
	respondJSON(w, http.StatusOK, user)
}
