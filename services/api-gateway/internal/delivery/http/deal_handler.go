package http

import (
"encoding/json"
"net/http"

"github.com/go-chi/chi/v5"
dealv1 "realty/api/gen/deal/v1"
"realty/services/api-gateway/internal/middleware"
)

type DealHandler struct {
client dealv1.DealServiceClient
}

func NewDealHandler(client dealv1.DealServiceClient) *DealHandler {
return &DealHandler{client: client}
}

func (h *DealHandler) CreateViewing(w http.ResponseWriter, r *http.Request) {
userID := middleware.GetUserID(r)

var req dealv1.CreateViewingRequestReq
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
writeError(w, http.StatusBadRequest, "неверный формат запроса")
return
}
req.BuyerId = userID

resp, err := h.client.CreateViewingRequest(r.Context(), &req)
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusCreated, resp)
}

func (h *DealHandler) UpdateViewing(w http.ResponseWriter, r *http.Request) {
id := chi.URLParam(r, "id")

var req dealv1.UpdateViewingRequestReq
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
writeError(w, http.StatusBadRequest, "неверный формат запроса")
return
}
req.Id = id

resp, err := h.client.UpdateViewingRequest(r.Context(), &req)
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}

func (h *DealHandler) CreateDeal(w http.ResponseWriter, r *http.Request) {
userID := middleware.GetUserID(r)

var req dealv1.CreateDealReq
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
writeError(w, http.StatusBadRequest, "неверный формат запроса")
return
}
req.BuyerId = userID

resp, err := h.client.CreateDeal(r.Context(), &req)
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusCreated, resp)
}

func (h *DealHandler) UpdateDeal(w http.ResponseWriter, r *http.Request) {
id := chi.URLParam(r, "id")

var req dealv1.UpdateDealReq
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
writeError(w, http.StatusBadRequest, "неверный формат запроса")
return
}
req.Id = id

resp, err := h.client.UpdateDeal(r.Context(), &req)
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}

func (h *DealHandler) AddToFavorites(w http.ResponseWriter, r *http.Request) {
userID := middleware.GetUserID(r)

var req dealv1.AddToFavoritesReq
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
writeError(w, http.StatusBadRequest, "неверный формат запроса")
return
}
req.UserId = userID

resp, err := h.client.AddToFavorites(r.Context(), &req)
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusCreated, resp)
}

func (h *DealHandler) RemoveFromFavorites(w http.ResponseWriter, r *http.Request) {
userID := middleware.GetUserID(r)
listingID := chi.URLParam(r, "listing_id")

resp, err := h.client.RemoveFromFavorites(r.Context(), &dealv1.RemoveFromFavoritesReq{
UserId:    userID,
ListingId: listingID,
})
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}

func (h *DealHandler) ListFavorites(w http.ResponseWriter, r *http.Request) {
userID := middleware.GetUserID(r)

resp, err := h.client.ListFavorites(r.Context(), &dealv1.ListFavoritesReq{
UserId: userID,
Page:   1,
Limit:  20,
})
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}
