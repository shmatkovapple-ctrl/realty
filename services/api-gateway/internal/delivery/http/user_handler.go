package http

import (
	"encoding/json"
	"net/http"

	userv1 "realty/api/gen/user/v1"
	"realty/services/api-gateway/internal/middleware"
)

type UserHandler struct {
	client userv1.UserServiceClient
}

func NewUserHandler(client userv1.UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req userv1.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}

	resp, err := h.client.Register(r.Context(), &req)
	if err != nil {
		writeGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userv1.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}

	resp, err := h.client.Login(r.Context(), &req)
	if err != nil {
		writeGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req userv1.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}

	resp, err := h.client.Logout(r.Context(), &req)
	if err != nil {
		writeGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req userv1.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}

	resp, err := h.client.RefreshToken(r.Context(), &req)
	if err != nil {
		writeGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	resp, err := h.client.GetProfile(r.Context(), &userv1.GetProfileRequest{
		UserId: userID,
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req userv1.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный формат запроса")
		return
	}
	req.UserId = userID

	resp, err := h.client.UpdateProfile(r.Context(), &req)
	if err != nil {
		writeGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
