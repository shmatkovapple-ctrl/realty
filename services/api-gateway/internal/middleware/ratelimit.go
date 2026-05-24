package middleware

import (
"fmt"
"net/http"
"time"

"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
client  *redis.Client
limit   int
window  time.Duration
}

func NewRateLimiter(client *redis.Client, limit int, window time.Duration) *RateLimiter {
return &RateLimiter{client: client, limit: limit, window: window}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
ip := r.RemoteAddr
key := fmt.Sprintf("rate:%s", ip)

ctx := r.Context()
count, err := rl.client.Incr(ctx, key).Result()
if err != nil {
next.ServeHTTP(w, r)
return
}

if count == 1 {
rl.client.Expire(ctx, key, rl.window)
}

if count > int64(rl.limit) {
writeError(w, http.StatusTooManyRequests, "слишком много запросов")
return
}

next.ServeHTTP(w, r)
})
}
