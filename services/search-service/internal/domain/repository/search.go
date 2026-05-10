package repository

import (
"context"

"realty/services/search-service/internal/domain/entity"
)

type SearchRepository interface {
Index(ctx context.Context, listing *entity.ListingIndex) error
Delete(ctx context.Context, listingID string) error
Search(ctx context.Context, filter entity.SearchFilter, page, limit int, sortBy string) (*entity.SearchResult, error)
Autocomplete(ctx context.Context, query, field string) ([]string, error)
}
