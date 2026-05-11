package middleware

import (
"context"
"encoding/json"
"net/http"
"strings"

"google.golang.org/grpc"
"google.golang.org/grpc/credentials/insecure"
userv1 "realty/api/gen/user/v1"
)

type contextKey string

const (
ContextKeyUserID contextKey = "user_id"
ContextKeyRole   contextKey = "role"
)

type AuthMiddleware struct {
userClient userv1.UserServiceClient
}

func NewAuthMiddleware(userServiceAddr string) (*AuthMiddleware, error) {
conn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
return nil, err
}
return &AuthMiddleware{userClient: userv1.NewUserServiceClient(conn)}, nil
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
authHeader := r.Header.Get("Authorization")
if authHeader == "" {
writeError(w, http.StatusUnauthorized, "требуется авторизация")
return
}

parts := strings.SplitN(authHeader, " ", 2)
if len(parts) != 2 || parts[0] != "Bearer" {
writeError(w, http.StatusUnauthorized, "неверный формат токена")
return
}

resp, err := m.userClient.ValidateToken(r.Context(), &userv1.ValidateTokenRequest{
AccessToken: parts[1],
})
if err != nil || !resp.Valid {
writeError(w, http.StatusUnauthorized, "недействительный токен")
return
}

ctx := context.WithValue(r.Context(), ContextKeyUserID, resp.UserId)
ctx = context.WithValue(ctx, ContextKeyRole, resp.Role)
next.ServeHTTP(w, r.WithContext(ctx))
})
}

func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
return func(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
userRole, ok := r.Context().Value(ContextKeyRole).(string)
if !ok || userRole != role {
writeError(w, http.StatusForbidden, "недостаточно прав")
return
}
next.ServeHTTP(w, r)
})
}
}

func GetUserID(r *http.Request) string {
id, _ := r.Context().Value(ContextKeyUserID).(string)
return id
}

func GetRole(r *http.Request) string {
role, _ := r.Context().Value(ContextKeyRole).(string)
return role
}

func writeError(w http.ResponseWriter, code int, msg string) {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(code)
json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
