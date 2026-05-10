package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"realty/services/search-service/internal/domain/entity"
	"realty/services/search-service/internal/domain/repository"

	"github.com/elastic/go-elasticsearch/v8"
)

const indexName = "listings"

type searchRepository struct {
	client *elasticsearch.Client
}

func NewSearchRepository(client *elasticsearch.Client) repository.SearchRepository {
	return &searchRepository{client: client}
}

func (r *searchRepository) Index(ctx context.Context, listing *entity.ListingIndex) error {
	data, err := json.Marshal(listing)
	if err != nil {
		return fmt.Errorf("сериализация объявления: %w", err)
	}

	res, err := r.client.Index(
		indexName,
		bytes.NewReader(data),
		r.client.Index.WithDocumentID(listing.ListingID),
		r.client.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("индексация объявления: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("ошибка индексации: %s", string(body))
	}

	return nil
}

func (r *searchRepository) Delete(ctx context.Context, listingID string) error {
	res, err := r.client.Delete(
		indexName,
		listingID,
		r.client.Delete.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("удаление из индекса: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("ошибка удаления: %s", string(body))
	}

	return nil
}

func (r *searchRepository) Search(ctx context.Context, filter entity.SearchFilter, page, limit int, sortBy string) (*entity.SearchResult, error) {
	query := buildSearchQuery(filter)
	sort := buildSortQuery(sortBy)
	from := (page - 1) * limit

	body := map[string]any{
		"query": query,
		"sort":  sort,
		"from":  from,
		"size":  limit,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("сериализация запроса: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithIndex(indexName),
		r.client.Search.WithBody(bytes.NewReader(data)),
		r.client.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("поиск: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("ошибка поиска: %s", string(body))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("декодирование ответа: %w", err)
	}

	return parseSearchResult(result, page), nil
}

func (r *searchRepository) Autocomplete(ctx context.Context, query, field string) ([]string, error) {
	if field == "" {
		field = "city"
	}

	body := map[string]any{
		"query": map[string]any{
			"prefix": map[string]any{
				field: query,
			},
		},
		"_source": []string{field},
		"size":    10,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("сериализация запроса автодополнения: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithIndex(indexName),
		r.client.Search.WithBody(bytes.NewReader(data)),
		r.client.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("автодополнение: %w", err)
	}
	defer res.Body.Close()

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("декодирование ответа: %w", err)
	}

	var suggestions []string
	seen := map[string]bool{}

	hits, _ := result["hits"].(map[string]any)
	hitsArr, _ := hits["hits"].([]any)
	for _, h := range hitsArr {
		hit, _ := h.(map[string]any)
		source, _ := hit["_source"].(map[string]any)
		val, _ := source[field].(string)
		if val != "" && !seen[val] {
			suggestions = append(suggestions, val)
			seen[val] = true
		}
	}

	return suggestions, nil
}

func buildSearchQuery(f entity.SearchFilter) map[string]any {
	must := []map[string]any{}
	filter := []map[string]any{
		{"term": map[string]any{"status": "published"}},
	}

	if f.Query != "" {
		must = append(must, map[string]any{
			"multi_match": map[string]any{
				"query":  f.Query,
				"fields": []string{"title^3", "description", "city", "district"},
			},
		})
	}

	if f.Type != "" {
		filter = append(filter, map[string]any{
			"term": map[string]any{"type": f.Type},
		})
	}

	if f.City != "" {
		filter = append(filter, map[string]any{
			"term": map[string]any{"city": f.City},
		})
	}

	if f.District != "" {
		filter = append(filter, map[string]any{
			"term": map[string]any{"district": f.District},
		})
	}

	if f.PriceMin > 0 || f.PriceMax > 0 {
		priceRange := map[string]any{}
		if f.PriceMin > 0 {
			priceRange["gte"] = f.PriceMin
		}
		if f.PriceMax > 0 {
			priceRange["lte"] = f.PriceMax
		}
		filter = append(filter, map[string]any{
			"range": map[string]any{"price": priceRange},
		})
	}

	if f.AreaMin > 0 || f.AreaMax > 0 {
		areaRange := map[string]any{}
		if f.AreaMin > 0 {
			areaRange["gte"] = f.AreaMin
		}
		if f.AreaMax > 0 {
			areaRange["lte"] = f.AreaMax
		}
		filter = append(filter, map[string]any{
			"range": map[string]any{"area_sqm": areaRange},
		})
	}

	if f.Rooms > 0 {
		filter = append(filter, map[string]any{
			"term": map[string]any{"rooms": f.Rooms},
		})
	}

	if f.Lat != 0 && f.Lng != 0 && f.RadiusKm > 0 {
		filter = append(filter, map[string]any{
			"geo_distance": map[string]any{
				"distance": fmt.Sprintf("%.1fkm", f.RadiusKm),
				"location": map[string]any{
					"lat": f.Lat,
					"lon": f.Lng,
				},
			},
		})
	}

	q := map[string]any{
		"bool": map[string]any{
			"filter": filter,
		},
	}

	if len(must) > 0 {
		q["bool"].(map[string]any)["must"] = must
	}

	return q
}

func buildSortQuery(sortBy string) []map[string]any {
	switch sortBy {
	case "price_asc":
		return []map[string]any{{"price": map[string]any{"order": "asc"}}}
	case "price_desc":
		return []map[string]any{{"price": map[string]any{"order": "desc"}}}
	case "area_desc":
		return []map[string]any{{"area_sqm": map[string]any{"order": "desc"}}}
	default:
		return []map[string]any{{"_score": map[string]any{"order": "desc"}}}
	}
}

func parseSearchResult(result map[string]any, page int) *entity.SearchResult {
	sr := &entity.SearchResult{Page: page}

	hits, _ := result["hits"].(map[string]any)
	total, _ := hits["total"].(map[string]any)
	totalVal, _ := total["value"].(float64)
	sr.Total = int(totalVal)

	hitsArr, _ := hits["hits"].([]any)
	for _, h := range hitsArr {
		hit, _ := h.(map[string]any)
		source, _ := hit["_source"].(map[string]any)
		score, _ := hit["_score"].(float64)

		sh := &entity.SearchHit{
			Score: float32(score),
		}

		if v, ok := source["listing_id"].(string); ok {
			sh.ListingID = v
		}
		if v, ok := source["title"].(string); ok {
			sh.Title = v
		}
		if v, ok := source["price"].(float64); ok {
			sh.Price = v
		}
		if v, ok := source["city"].(string); ok {
			sh.City = v
		}
		if v, ok := source["district"].(string); ok {
			sh.District = v
		}
		if v, ok := source["area_sqm"].(float64); ok {
			sh.AreaSqm = v
		}
		if v, ok := source["rooms"].(float64); ok {
			sh.Rooms = int32(v)
		}
		if v, ok := source["preview_url"].(string); ok {
			sh.PreviewURL = v
		}
		if v, ok := source["lat"].(float64); ok {
			sh.Lat = v
		}
		if v, ok := source["lng"].(float64); ok {
			sh.Lng = v
		}

		sr.Hits = append(sr.Hits, sh)
	}

	return sr
}

func EnsureIndex(ctx context.Context, client *elasticsearch.Client) error {
	res, err := client.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("проверка индекса: %w", err)
	}
	res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	mapping := `{
"mappings": {
"properties": {
"listing_id":  { "type": "keyword" },
"type":        { "type": "keyword" },
"status":      { "type": "keyword" },
"title":       { "type": "text", "analyzer": "russian" },
"description": { "type": "text", "analyzer": "russian" },
"price":       { "type": "double" },
"area_sqm":    { "type": "double" },
"rooms":       { "type": "integer" },
"city":        { "type": "keyword" },
"district":    { "type": "keyword" },
"preview_url": { "type": "keyword" },
"location": {
"type": "geo_point"
}
}
}
}`

	createRes, err := client.Indices.Create(
		indexName,
		client.Indices.Create.WithBody(strings.NewReader(mapping)),
		client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("создание индекса: %w", err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		body, _ := io.ReadAll(createRes.Body)
		return fmt.Errorf("ошибка создания индекса: %s", string(body))
	}

	return nil
}
