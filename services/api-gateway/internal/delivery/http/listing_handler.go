package http

import (
"encoding/json"
"net/http"

"github.com/go-chi/chi/v5"
listingv1 "realty/api/gen/listing/v1"
"realty/services/api-gateway/internal/middleware"
)

type ListingHandler struct {
client listingv1.ListingServiceClient
}

func NewListingHandler(client listingv1.ListingServiceClient) *ListingHandler {
return &ListingHandler{client: client}
}

func (h *ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
userID := middleware.GetUserID(r)

var listing listingv1.Listing
if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
writeError(w, http.StatusBadRequest, "неверный формат запроса")
return
}
listing.SellerId = userID

resp, err := h.client.CreateListing(r.Context(), &listingv1.CreateListingRequest{Listing: &listing})
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusCreated, resp)
}

func (h *ListingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
id := chi.URLParam(r, "id")

resp, err := h.client.GetListing(r.Context(), &listingv1.GetListingRequest{Id: id})
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}

func (h *ListingHandler) Update(w http.ResponseWriter, r *http.Request) {
id := chi.URLParam(r, "id")
userID := middleware.GetUserID(r)

var listing listingv1.Listing
if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
writeError(w, http.StatusBadRequest, "неверный формат запроса")
return
}
listing.Id = id
listing.SellerId = userID

resp, err := h.client.UpdateListing(r.Context(), &listingv1.UpdateListingRequest{Listing: &listing})
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}

func (h *ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
id := chi.URLParam(r, "id")
userID := middleware.GetUserID(r)

resp, err := h.client.DeleteListing(r.Context(), &listingv1.DeleteListingRequest{
Id:       id,
SellerId: userID,
})
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}

func (h *ListingHandler) Publish(w http.ResponseWriter, r *http.Request) {
id := chi.URLParam(r, "id")
userID := middleware.GetUserID(r)

resp, err := h.client.PublishListing(r.Context(), &listingv1.PublishListingRequest{
Id:       id,
SellerId: userID,
})
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}

func (h *ListingHandler) GetUploadURL(w http.ResponseWriter, r *http.Request) {
id := chi.URLParam(r, "id")

var req listingv1.GetUploadURLRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
writeError(w, http.StatusBadRequest, "неверный формат запроса")
return
}
req.ListingId = id

resp, err := h.client.GetUploadURL(r.Context(), &req)
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}
