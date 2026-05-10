package grpc

import (
"context"

"google.golang.org/grpc/codes"
"google.golang.org/grpc/status"
searchv1 "realty/api/gen/search/v1"
"realty/services/search-service/internal/domain/entity"
"realty/services/search-service/internal/usecase"
)

type SearchHandler struct {
searchv1.UnimplementedSearchServiceServer
uc usecase.SearchUseCase
}

func NewSearchHandler(uc usecase.SearchUseCase) *SearchHandler {
return &SearchHandler{uc: uc}
}

func (h *SearchHandler) Search(ctx context.Context, req *searchv1.SearchRequest) (*searchv1.SearchResponse, error) {
filter := entity.SearchFilter{}

if req.Filter != nil {
filter.Query    = req.Filter.Query
filter.Type     = req.Filter.Type
filter.City     = req.Filter.City
filter.District = req.Filter.District
filter.Rooms    = req.Filter.Rooms

if req.Filter.Price != nil {
filter.PriceMin = req.Filter.Price.Min
filter.PriceMax = req.Filter.Price.Max
}
if req.Filter.Area != nil {
filter.AreaMin = req.Filter.Area.Min
filter.AreaMax = req.Filter.Area.Max
}
if req.Filter.Geo != nil {
filter.Lat      = req.Filter.Geo.Lat
filter.Lng      = req.Filter.Geo.Lng
filter.RadiusKm = req.Filter.Geo.RadiusKm
}
}

result, err := h.uc.Search(ctx, usecase.SearchInput{
Filter: filter,
Page:   int(req.Page),
Limit:  int(req.Limit),
SortBy: req.Sort.String(),
})
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

var hits []*searchv1.SearchHit
for _, h := range result.Hits {
hits = append(hits, &searchv1.SearchHit{
ListingId:  h.ListingID,
Score:      h.Score,
Title:      h.Title,
Price:      h.Price,
City:       h.City,
District:   h.District,
AreaSqm:    h.AreaSqm,
Rooms:      h.Rooms,
PreviewUrl: h.PreviewURL,
Lat:        h.Lat,
Lng:        h.Lng,
})
}

return &searchv1.SearchResponse{
Hits:  hits,
Total: int32(result.Total),
Page:  int32(result.Page),
}, nil
}

func (h *SearchHandler) Autocomplete(ctx context.Context, req *searchv1.AutocompleteRequest) (*searchv1.AutocompleteResponse, error) {
suggestions, err := h.uc.Autocomplete(ctx, req.Query, req.Field)
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &searchv1.AutocompleteResponse{Suggestions: suggestions}, nil
}

func (h *SearchHandler) GetSuggestions(ctx context.Context, req *searchv1.GetSuggestionsRequest) (*searchv1.GetSuggestionsResponse, error) {
result, err := h.uc.Search(ctx, usecase.SearchInput{
Filter: entity.SearchFilter{},
Page:   1,
Limit:  10,
SortBy: "created_desc",
})
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

var hits []*searchv1.SearchHit
for _, h := range result.Hits {
hits = append(hits, &searchv1.SearchHit{
ListingId:  h.ListingID,
Score:      h.Score,
Title:      h.Title,
Price:      h.Price,
City:       h.City,
District:   h.District,
AreaSqm:    h.AreaSqm,
Rooms:      h.Rooms,
PreviewUrl: h.PreviewURL,
})
}

return &searchv1.GetSuggestionsResponse{Hits: hits}, nil
}

func (h *SearchHandler) IndexListing(ctx context.Context, req *searchv1.IndexListingRequest) (*searchv1.IndexListingResponse, error) {
listing := &entity.ListingIndex{
ListingID:   req.ListingId,
Type:        req.Type,
Status:      req.Status,
Title:       req.Title,
Description: req.Description,
Price:       req.Price,
AreaSqm:     req.AreaSqm,
Rooms:       req.Rooms,
City:        req.City,
District:    req.District,
Lat:         req.Lat,
Lng:         req.Lng,
PreviewURL:  req.PreviewUrl,
}

if err := h.uc.IndexListing(ctx, listing); err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &searchv1.IndexListingResponse{Success: true}, nil
}

func (h *SearchHandler) DeleteIndex(ctx context.Context, req *searchv1.DeleteIndexRequest) (*searchv1.DeleteIndexResponse, error) {
if err := h.uc.DeleteIndex(ctx, req.ListingId); err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &searchv1.DeleteIndexResponse{Success: true}, nil
}
