package grpc

import (
"context"

"github.com/google/uuid"
"google.golang.org/grpc/codes"
"google.golang.org/grpc/status"
listingv1 "realty/api/gen/listing/v1"
"realty/services/listing-service/internal/domain/entity"
"realty/services/listing-service/internal/usecase"
)

type ListingHandler struct {
listingv1.UnimplementedListingServiceServer
uc usecase.ListingUseCase
}

func NewListingHandler(uc usecase.ListingUseCase) *ListingHandler {
return &ListingHandler{uc: uc}
}

func (h *ListingHandler) CreateListing(ctx context.Context, req *listingv1.CreateListingRequest) (*listingv1.CreateListingResponse, error) {
sellerID, err := uuid.Parse(req.Listing.SellerId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный seller_id")
}

listing, err := h.uc.Create(ctx, usecase.CreateInput{
SellerID:    sellerID,
Type:        protoTypeToEntity(req.Listing.Type),
Title:       req.Listing.Title,
Description: req.Listing.Description,
Price:       req.Listing.Price,
AreaSqm:     req.Listing.AreaSqm,
Rooms:       req.Listing.Rooms,
Floor:       req.Listing.Floor,
FloorsTotal: req.Listing.FloorsTotal,
Address:     toEntityAddress(req.Listing.Address),
})
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, err.Error())
}

return &listingv1.CreateListingResponse{Listing: toProtoListing(listing)}, nil
}

func (h *ListingHandler) GetListing(ctx context.Context, req *listingv1.GetListingRequest) (*listingv1.GetListingResponse, error) {
id, err := uuid.Parse(req.Id)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный id")
}

listing, err := h.uc.GetByID(ctx, id)
if err != nil {
return nil, status.Errorf(codes.NotFound, err.Error())
}

return &listingv1.GetListingResponse{Listing: toProtoListing(listing)}, nil
}

func (h *ListingHandler) UpdateListing(ctx context.Context, req *listingv1.UpdateListingRequest) (*listingv1.UpdateListingResponse, error) {
id, err := uuid.Parse(req.Listing.Id)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный id")
}

sellerID, err := uuid.Parse(req.Listing.SellerId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный seller_id")
}

listing, err := h.uc.Update(ctx, usecase.UpdateInput{
ID:          id,
SellerID:    sellerID,
Title:       req.Listing.Title,
Description: req.Listing.Description,
Price:       req.Listing.Price,
AreaSqm:     req.Listing.AreaSqm,
Rooms:       req.Listing.Rooms,
Floor:       req.Listing.Floor,
FloorsTotal: req.Listing.FloorsTotal,
Address:     toEntityAddress(req.Listing.Address),
})
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &listingv1.UpdateListingResponse{Listing: toProtoListing(listing)}, nil
}

func (h *ListingHandler) DeleteListing(ctx context.Context, req *listingv1.DeleteListingRequest) (*listingv1.DeleteListingResponse, error) {
id, err := uuid.Parse(req.Id)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный id")
}

sellerID, err := uuid.Parse(req.SellerId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный seller_id")
}

if err := h.uc.Delete(ctx, id, sellerID); err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &listingv1.DeleteListingResponse{Success: true}, nil
}

func (h *ListingHandler) PublishListing(ctx context.Context, req *listingv1.PublishListingRequest) (*listingv1.PublishListingResponse, error) {
id, err := uuid.Parse(req.Id)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный id")
}

sellerID, err := uuid.Parse(req.SellerId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный seller_id")
}

listing, err := h.uc.Publish(ctx, id, sellerID)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, err.Error())
}

return &listingv1.PublishListingResponse{Listing: toProtoListing(listing)}, nil
}

func (h *ListingHandler) GetUploadURL(ctx context.Context, req *listingv1.GetUploadURLRequest) (*listingv1.GetUploadURLResponse, error) {
uploadURL, fileURL, err := h.uc.GetUploadURL(ctx, req.ListingId, req.Filename, req.ContentType)
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &listingv1.GetUploadURLResponse{
UploadUrl: uploadURL,
FileUrl:   fileURL,
}, nil
}

func (h *ListingHandler) SearchListings(ctx context.Context, req *listingv1.SearchListingsRequest) (*listingv1.SearchListingsResponse, error) {
return nil, status.Errorf(codes.Unimplemented, "поиск реализован в search-service")
}

func protoTypeToEntity(t listingv1.ListingType) string {
switch t {
case listingv1.ListingType_LISTING_TYPE_APARTMENT:
return "apartment"
case listingv1.ListingType_LISTING_TYPE_HOUSE:
return "house"
case listingv1.ListingType_LISTING_TYPE_COMMERCIAL:
return "commercial"
case listingv1.ListingType_LISTING_TYPE_LAND:
return "land"
default:
return ""
}
}

func toEntityAddress(a *listingv1.Address) entity.Address {
if a == nil {
return entity.Address{}
}
return entity.Address{
Country:  a.Country,
City:     a.City,
District: a.District,
Street:   a.Street,
Building: a.Building,
Lat:      a.Lat,
Lng:      a.Lng,
}
}

func toProtoListing(l *entity.Listing) *listingv1.Listing {
proto := &listingv1.Listing{
Id:          l.ID.String(),
SellerId:    l.SellerID.String(),
Title:       l.Title,
Description: l.Description,
Price:       l.Price,
Currency:    l.Currency,
AreaSqm:     l.AreaSqm,
Rooms:       l.Rooms,
Floor:       l.Floor,
FloorsTotal: l.FloorsTotal,
MediaUrls:   l.MediaURLs,
Address: &listingv1.Address{
Country:  l.Address.Country,
City:     l.Address.City,
District: l.Address.District,
Street:   l.Address.Street,
Building: l.Address.Building,
Lat:      l.Address.Lat,
Lng:      l.Address.Lng,
},
}

if l.AgentID != nil {
proto.AgentId = l.AgentID.String()
}

return proto
}
