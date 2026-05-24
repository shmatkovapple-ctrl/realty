package http

import (
"net/http"
"time"

"github.com/go-chi/chi/v5"
chimiddleware "github.com/go-chi/chi/v5/middleware"
"github.com/go-chi/cors"
"realty/services/api-gateway/internal/middleware"
)

func NewRouter(
authMW      *middleware.AuthMiddleware,
rateLimiter *middleware.RateLimiter,
user         *UserHandler,
listing      *ListingHandler,
search       *SearchHandler,
deal         *DealHandler,
notification *NotificationHandler,
) http.Handler {
r := chi.NewRouter()

r.Use(chimiddleware.Logger)
r.Use(chimiddleware.Recoverer)
r.Use(chimiddleware.Timeout(30 * time.Second))
r.Use(chimiddleware.RequestID)
r.Use(rateLimiter.Limit)
r.Use(cors.Handler(cors.Options{
AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
AllowedHeaders:   []string{"Authorization", "Content-Type"},
AllowCredentials: true,
}))

r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
})

r.Route("/api/v1", func(r chi.Router) {
r.Route("/auth", func(r chi.Router) {
r.Post("/register", user.Register)
r.Post("/login", user.Login)
r.Post("/logout", user.Logout)
r.Post("/refresh", user.RefreshToken)
})

r.Route("/listings", func(r chi.Router) {
r.Get("/", search.Search)
r.Get("/autocomplete", search.Autocomplete)
r.With(authMW.Authenticate, authMW.RequireAnyRole("seller", "agent")).Get("/mine", listing.ListMine)
r.Get("/{id}", listing.GetByID)

r.Group(func(r chi.Router) {
r.Use(authMW.Authenticate)
r.Use(authMW.RequireAnyRole("seller", "agent"))
r.Post("/", listing.Create)
r.Put("/{id}", listing.Update)
r.Delete("/{id}", listing.Delete)
r.Post("/{id}/publish", listing.Publish)
r.Post("/{id}/upload-url", listing.GetUploadURL)
r.Post("/{id}/upload", listing.UploadPhoto)
})
})

r.Group(func(r chi.Router) {
r.Use(authMW.Authenticate)

r.Route("/profile", func(r chi.Router) {
r.Get("/", user.GetProfile)
r.Put("/", user.UpdateProfile)
})

r.Route("/viewings", func(r chi.Router) {
r.Post("/", deal.CreateViewing)
r.Put("/{id}", deal.UpdateViewing)
})

r.Route("/deals", func(r chi.Router) {
r.Post("/", deal.CreateDeal)
r.Put("/{id}", deal.UpdateDeal)
})

r.Route("/favorites", func(r chi.Router) {
r.Get("/", deal.ListFavorites)
r.Post("/", deal.AddToFavorites)
r.Delete("/{listing_id}", deal.RemoveFromFavorites)
})

r.Route("/notifications", func(r chi.Router) {
r.Get("/", notification.List)
r.Put("/{id}/read", notification.MarkAsRead)
r.Put("/read-all", notification.MarkAllAsRead)
})
})
})

return r
}