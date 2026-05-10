package redis

import (
"context"
"encoding/json"
"fmt"
"time"

"github.com/redis/go-redis/v9"
"realty/services/search-service/internal/domain/entity"
)

const cacheTTL = 5 * time.Minute

type SearchCache struct {
client *redis.Client
}

func NewSearchCache(client *redis.Client) *SearchCache {
return &SearchCache{client: client}
}

func (c *SearchCache) Get(ctx context.Context, key string) (*entity.SearchResult, error) {
val, err := c.client.Get(ctx, key).Result()
if err == redis.Nil {
return nil, nil
}
if err != nil {
return nil, fmt.Errorf("получение из кэша: %w", err)
}

var result entity.SearchResult
if err := json.Unmarshal([]byte(val), &result); err != nil {
return nil, fmt.Errorf("десериализация кэша: %w", err)
}

return &result, nil
}

func (c *SearchCache) Set(ctx context.Context, key string, result *entity.SearchResult) error {
data, err := json.Marshal(result)
if err != nil {
return fmt.Errorf("сериализация кэша: %w", err)
}

return c.client.Set(ctx, key, data, cacheTTL).Err()
}

func (c *SearchCache) BuildKey(filter entity.SearchFilter, page, limit int, sortBy string) string {
return fmt.Sprintf("search:%s:%s:%s:%s:%.0f:%.0f:%.0f:%.0f:%d:%.4f:%.4f:%.1f:%d:%d:%s",
filter.Query, filter.Type, filter.City, filter.District,
filter.PriceMin, filter.PriceMax, filter.AreaMin, filter.AreaMax,
filter.Rooms, filter.Lat, filter.Lng, filter.RadiusKm,
page, limit, sortBy,
)
}
