package entity

type ListingIndex struct {
ListingID  string  `json:"listing_id"`
Type       string  `json:"type"`
Status     string  `json:"status"`
Title      string  `json:"title"`
Description string `json:"description"`
Price      float64 `json:"price"`
AreaSqm    float64 `json:"area_sqm"`
Rooms      int32   `json:"rooms"`
City       string  `json:"city"`
District   string  `json:"district"`
Lat        float64 `json:"lat"`
Lng        float64 `json:"lng"`
PreviewURL string  `json:"preview_url"`
}

type SearchFilter struct {
Query    string
Type     string
City     string
District string
PriceMin float64
PriceMax float64
AreaMin  float64
AreaMax  float64
Rooms    int32
Lat      float64
Lng      float64
RadiusKm float64
}

type SearchResult struct {
Hits  []*SearchHit
Total int
Page  int
}

type SearchHit struct {
ListingID  string  `json:"listing_id"`
Score      float32 `json:"score"`
Title      string  `json:"title"`
Price      float64 `json:"price"`
City       string  `json:"city"`
District   string  `json:"district"`
AreaSqm    float64 `json:"area_sqm"`
Rooms      int32   `json:"rooms"`
PreviewURL string  `json:"preview_url"`
Lat        float64 `json:"lat"`
Lng        float64 `json:"lng"`
}
