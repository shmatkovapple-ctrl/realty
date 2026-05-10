package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenCache struct {
	client *redis.Client
}

func NewTokenCache(client *redis.Client) *TokenCache {
	return &TokenCache{client: client}
}

func (c *TokenCache) SaveRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	key := fmt.Sprintf("refresh:%s:%s", userID, token)
	return c.client.Set(ctx, key, userID, ttl).Err()
}

func (c *TokenCache) ValidateRefreshToken(ctx context.Context, userID, token string) (bool, error) {
	key := fmt.Sprintf("refresh:%s:%s", userID, token)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("проверка refresh токена: %w", err)
	}
	return val == userID, nil
}

func (c *TokenCache) DeleteRefreshToken(ctx context.Context, userID, token string) error {
	key := fmt.Sprintf("refresh:%s:%s", userID, token)
	return c.client.Del(ctx, key).Err()
}

func (c *TokenCache) DeleteAllUserTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("refresh:%s:*", userID)
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("получение ключей: %w", err)
	}
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}
