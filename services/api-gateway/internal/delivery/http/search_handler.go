package http

import (
"net/http"
"net/url"
"strconv"

searchv1 "realty/api/gen/search/v1"
)

type SearchHandler struct {
client searchv1.SearchServiceClient
}

func NewSearchHandler(client searchv1.SearchServiceClient) *SearchHandler {
return &SearchHandler{client: client}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
q := r.URL.Query()

priceMin, _ := strconv.ParseFloat(q.Get("price_min"), 64)
priceMax, _ := strconv.ParseFloat(q.Get("price_max"), 64)
areaMin, _  := strconv.ParseFloat(q.Get("area_min"), 64)
areaMax, _  := strconv.ParseFloat(q.Get("area_max"), 64)
rooms, _    := strconv.Atoi(q.Get("rooms"))
lat, _      := strconv.ParseFloat(q.Get("lat"), 64)
lng, _      := strconv.ParseFloat(q.Get("lng"), 64)
radius, _   := strconv.ParseFloat(q.Get("radius_km"), 64)
page, _     := strconv.Atoi(q.Get("page"))
limit, _    := strconv.Atoi(q.Get("limit"))

if page < 1  { page = 1 }
if limit < 1 { limit = 20 }

city, _     := url.QueryUnescape(q.Get("city"))
district, _ := url.QueryUnescape(q.Get("district"))
queryStr, _ := url.QueryUnescape(q.Get("q"))
typeStr, _  := url.QueryUnescape(q.Get("type"))

filter := &searchv1.SearchFilter{
Query:    queryStr,
Type:     typeStr,
City:     city,
District: district,
Rooms:    int32(rooms),
}

if priceMin > 0 || priceMax > 0 {
filter.Price = &searchv1.PriceRange{Min: priceMin, Max: priceMax}
}
if areaMin > 0 || areaMax > 0 {
filter.Area = &searchv1.AreaRange{Min: areaMin, Max: areaMax}
}
if lat != 0 && lng != 0 {
filter.Geo = &searchv1.GeoPoint{Lat: lat, Lng: lng, RadiusKm: radius}
}

resp, err := h.client.Search(r.Context(), &searchv1.SearchRequest{
Filter: filter,
Page:   int32(page),
Limit:  int32(limit),
})
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}

func (h *SearchHandler) Autocomplete(w http.ResponseWriter, r *http.Request) {
q := r.URL.Query()

query, _ := url.QueryUnescape(q.Get("q"))
field, _ := url.QueryUnescape(q.Get("field"))

resp, err := h.client.Autocomplete(r.Context(), &searchv1.AutocompleteRequest{
Query: query,
Field: field,
})
if err != nil {
writeGRPCError(w, err)
return
}

writeJSON(w, http.StatusOK, resp)
}
