package usecase

import (
"context"
"fmt"

"realty/services/search-service/internal/domain/entity"
"realty/services/search-service/internal/domain/repository"
redisInfra "realty/services/search-service/internal/infrastructure/redis"
)

type searchUseCase struct {
searchRepo  repository.SearchRepository
searchCache *redisInfra.SearchCache
}

func NewSearchUseCase(
searchRepo repository.SearchRepository,
searchCache *redisInfra.SearchCache,
) SearchUseCase {
return &searchUseCase{
searchRepo:  searchRepo,
searchCache: searchCache,
}
}

func (uc *searchUseCase) Search(ctx context.Context, input SearchInput) (*entity.SearchResult, error) {
if input.Page < 1 {
input.Page = 1
}
if input.Limit < 1 || input.Limit > 100 {
input.Limit = 20
}

cacheKey := uc.searchCache.BuildKey(input.Filter, input.Page, input.Limit, input.SortBy)

cached, err := uc.searchCache.Get(ctx, cacheKey)
if err == nil && cached != nil {
return cached, nil
}

result, err := uc.searchRepo.Search(ctx, input.Filter, input.Page, input.Limit, input.SortBy)
if err != nil {
return nil, fmt.Errorf("поиск объявлений: %w", err)
}

_ = uc.searchCache.Set(ctx, cacheKey, result)

return result, nil
}

func (uc *searchUseCase) Autocomplete(ctx context.Context, query, field string) ([]string, error) {
if query == "" {
return nil, nil
}

suggestions, err := uc.searchRepo.Autocomplete(ctx, query, field)
if err != nil {
return nil, fmt.Errorf("автодополнение: %w", err)
}

return suggestions, nil
}

func (uc *searchUseCase) IndexListing(ctx context.Context, listing *entity.ListingIndex) error {
if err := uc.searchRepo.Index(ctx, listing); err != nil {
return fmt.Errorf("индексация объявления: %w", err)
}
return nil
}

func (uc *searchUseCase) DeleteIndex(ctx context.Context, listingID string) error {
if err := uc.searchRepo.Delete(ctx, listingID); err != nil {
return fmt.Errorf("удаление из индекса: %w", err)
}
return nil
}
