package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	listingv1 "realty/api/gen/listing/v1"
	"realty/services/api-gateway/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type ListingHandler struct {
	client     listingv1.ListingServiceClient
	httpClient *http.Client
}

func NewListingHandler(client listingv1.ListingServiceClient) *ListingHandler {
	return &ListingHandler{
		client:     client,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
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

func (h *ListingHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("ошибка парсинга формы: %v", err)
		writeError(w, http.StatusBadRequest, "ошибка парсинга формы")
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		log.Printf("файл не найден: %v", err)
		writeError(w, http.StatusBadRequest, "файл не найден")
		return
	}
	defer file.Close()

	urlResp, err := h.client.GetUploadURL(r.Context(), &listingv1.GetUploadURLRequest{
		ListingId:   id,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		log.Printf("ошибка получения upload URL: %v", err)
		writeGRPCError(w, err)
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("ошибка чтения файла: %v", err)
		writeError(w, http.StatusInternalServerError, "ошибка чтения файла")
		return
	}

	internalURL := strings.Replace(urlResp.UploadUrl, "localhost:9000", "minio:9000", 1)
	log.Printf("internal upload URL: %s", internalURL)

	req, err := http.NewRequestWithContext(r.Context(), "PUT", internalURL, bytes.NewReader(data))
	if err != nil {
		log.Printf("ошибка создания запроса: %v", err)
		writeError(w, http.StatusInternalServerError, "ошибка создания запроса")
		return
	}
	req.Header.Set("Content-Type", header.Header.Get("Content-Type"))

	resp, err := h.httpClient.Do(req)
	if err != nil {
		log.Printf("ошибка PUT в MinIO: %v", err)
		writeError(w, http.StatusInternalServerError, "ошибка загрузки в хранилище")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("MinIO ответ: %d %s", resp.StatusCode, string(body))

	if resp.StatusCode >= 400 {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("хранилище вернуло ошибку: %d", resp.StatusCode))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"file_url": urlResp.FileUrl})
}

func (h *ListingHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	resp, err := h.client.ListListingsBySeller(r.Context(), &listingv1.ListListingsBySellerRequest{
		SellerId: userID,
		Page:     int32(page),
		Limit:    int32(limit),
	})
	if err != nil {
		writeGRPCError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
