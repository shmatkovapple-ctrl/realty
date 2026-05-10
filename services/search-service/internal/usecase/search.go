package usecase

import (
"context"

"realty/services/search-service/internal/domain/entity"
)

type SearchUseCase interface {
Search(ctx context.Context, input SearchInput) (*entity.SearchResult, error)
Autocomplete(ctx context.Context, query, field string) ([]string, error)
IndexListing(ctx context.Context, listing *entity.ListingIndex) error
DeleteIndex(ctx context.Context, listingID string) error
}

type SearchInput struct {
Filter entity.SearchFilter
Page   int
Limit  int
SortBy string
}
