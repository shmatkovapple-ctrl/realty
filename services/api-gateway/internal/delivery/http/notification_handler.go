package http

import (
"net/http"
"strconv"

"github.com/go-chi/chi/v5"
notificationv1 "realty/api/gen/notification/v1"
"realty/services/api-gateway/internal/middleware"
)

type NotificationHandler struct {
client notificationv1.NotificationServiceClient
}

func NewNotificationHandler(client notificationv1.NotificationServiceClient) *NotificationHandler {
return &NotificationHandler{client: client}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
userID := middleware.GetUserID(r)
q := r.URL.Query()

unreadOnly := q.Get("unread_only") == "true"
page, _  := strconv.Atoi(q.Get("page"))
limit, _ := strconv.Atoi(q.Get("limit"))
if page < 1  { page = 1 }
if limit < 1 { limit = 20 }

resp, err := h.client.ListNotifications(r.Context(), &notificationv1.ListNotificationsReq{
UserId:     userID,
UnreadOnly: unreadOnly,
Page:       int32(page),
Limit:      int32(limit),
})
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
userID := middleware.GetUserID(r)
notifID := chi.URLParam(r, "id")

resp, err := h.client.MarkAsRead(r.Context(), &notificationv1.MarkAsReadReq{
NotificationId: notifID,
UserId:         userID,
})
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}

func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
userID := middleware.GetUserID(r)

resp, err := h.client.MarkAllAsRead(r.Context(), &notificationv1.MarkAllAsReadReq{
UserId: userID,
})
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}
