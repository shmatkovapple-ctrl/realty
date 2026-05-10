package grpc

import (
"context"
"time"

"github.com/google/uuid"
"google.golang.org/grpc/codes"
"google.golang.org/grpc/status"
dealv1 "realty/api/gen/deal/v1"
"realty/services/deal-service/internal/domain/entity"
"realty/services/deal-service/internal/usecase"
)

type DealHandler struct {
dealv1.UnimplementedDealServiceServer
dealUC    usecase.DealUseCase
viewingUC usecase.ViewingUseCase
favoriteUC usecase.FavoriteUseCase
}

func NewDealHandler(
dealUC usecase.DealUseCase,
viewingUC usecase.ViewingUseCase,
favoriteUC usecase.FavoriteUseCase,
) *DealHandler {
return &DealHandler{
dealUC:     dealUC,
viewingUC:  viewingUC,
favoriteUC: favoriteUC,
}
}

func (h *DealHandler) CreateViewingRequest(ctx context.Context, req *dealv1.CreateViewingRequestReq) (*dealv1.CreateViewingRequestRes, error) {
listingID, err := uuid.Parse(req.ListingId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный listing_id")
}
buyerID, err := uuid.Parse(req.BuyerId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный buyer_id")
}

var scheduledAt *time.Time
if req.ScheduledAt != nil {
t := req.ScheduledAt.AsTime()
scheduledAt = &t
}

viewing, err := h.viewingUC.CreateViewingRequest(ctx, usecase.CreateViewingInput{
ListingID:   listingID,
BuyerID:     buyerID,
Comment:     req.Comment,
ScheduledAt: scheduledAt,
})
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &dealv1.CreateViewingRequestRes{ViewingRequest: toProtoViewing(viewing)}, nil
}

func (h *DealHandler) UpdateViewingRequest(ctx context.Context, req *dealv1.UpdateViewingRequestReq) (*dealv1.UpdateViewingRequestRes, error) {
id, err := uuid.Parse(req.Id)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный id")
}

var scheduledAt *time.Time
if req.ScheduledAt != nil {
t := req.ScheduledAt.AsTime()
scheduledAt = &t
}

viewing, err := h.viewingUC.UpdateViewingRequest(ctx, usecase.UpdateViewingInput{
ID:          id,
Status:      protoToViewingStatus(req.Status),
ScheduledAt: scheduledAt,
})
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &dealv1.UpdateViewingRequestRes{ViewingRequest: toProtoViewing(viewing)}, nil
}

func (h *DealHandler) ListViewingRequests(ctx context.Context, req *dealv1.ListViewingRequestsReq) (*dealv1.ListViewingRequestsRes, error) {
input := usecase.ListViewingsInput{
Page:  int(req.Page),
Limit: int(req.Limit),
}

if req.ListingId != "" {
id, err := uuid.Parse(req.ListingId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный listing_id")
}
input.ListingID = &id
}

if req.BuyerId != "" {
id, err := uuid.Parse(req.BuyerId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный buyer_id")
}
input.BuyerID = &id
}

viewings, total, err := h.viewingUC.ListViewingRequests(ctx, input)
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

var protoViewings []*dealv1.ViewingRequest
for _, v := range viewings {
protoViewings = append(protoViewings, toProtoViewing(v))
}

return &dealv1.ListViewingRequestsRes{
ViewingRequests: protoViewings,
Total:           int32(total),
}, nil
}

func (h *DealHandler) CreateDeal(ctx context.Context, req *dealv1.CreateDealReq) (*dealv1.CreateDealRes, error) {
listingID, err := uuid.Parse(req.ListingId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный listing_id")
}
buyerID, err := uuid.Parse(req.BuyerId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный buyer_id")
}

var agentID *uuid.UUID
if req.AgentId != "" {
id, err := uuid.Parse(req.AgentId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный agent_id")
}
agentID = &id
}

deal, err := h.dealUC.CreateDeal(ctx, usecase.CreateDealInput{
ListingID:   listingID,
BuyerID:     buyerID,
AgentID:     agentID,
PriceAgreed: req.PriceAgreed,
Currency:    req.Currency,
})
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &dealv1.CreateDealRes{Deal: toProtoDeal(deal)}, nil
}

func (h *DealHandler) UpdateDeal(ctx context.Context, req *dealv1.UpdateDealReq) (*dealv1.UpdateDealRes, error) {
id, err := uuid.Parse(req.Id)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный id")
}

deal, err := h.dealUC.UpdateDeal(ctx, usecase.UpdateDealInput{
ID:          id,
Status:      protoToDealStatus(req.Status),
PriceAgreed: req.PriceAgreed,
})
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &dealv1.UpdateDealRes{Deal: toProtoDeal(deal)}, nil
}

func (h *DealHandler) GetDeal(ctx context.Context, req *dealv1.GetDealReq) (*dealv1.GetDealRes, error) {
id, err := uuid.Parse(req.Id)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный id")
}

deal, err := h.dealUC.GetDeal(ctx, id)
if err != nil {
return nil, status.Errorf(codes.NotFound, err.Error())
}

return &dealv1.GetDealRes{Deal: toProtoDeal(deal)}, nil
}

func (h *DealHandler) ListDeals(ctx context.Context, req *dealv1.ListDealsReq) (*dealv1.ListDealsRes, error) {
input := usecase.ListDealsInput{
Page:  int(req.Page),
Limit: int(req.Limit),
}

if req.BuyerId != "" {
id, err := uuid.Parse(req.BuyerId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный buyer_id")
}
input.BuyerID = &id
}

if req.ListingId != "" {
id, err := uuid.Parse(req.ListingId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный listing_id")
}
input.ListingID = &id
}

deals, total, err := h.dealUC.ListDeals(ctx, input)
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

var protoDeals []*dealv1.Deal
for _, d := range deals {
protoDeals = append(protoDeals, toProtoDeal(d))
}

return &dealv1.ListDealsRes{Deals: protoDeals, Total: int32(total)}, nil
}

func (h *DealHandler) AddToFavorites(ctx context.Context, req *dealv1.AddToFavoritesReq) (*dealv1.AddToFavoritesRes, error) {
userID, err := uuid.Parse(req.UserId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный user_id")
}
listingID, err := uuid.Parse(req.ListingId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный listing_id")
}

favorite, err := h.favoriteUC.AddToFavorites(ctx, userID, listingID)
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &dealv1.AddToFavoritesRes{Favorite: toProtoFavorite(favorite)}, nil
}

func (h *DealHandler) RemoveFromFavorites(ctx context.Context, req *dealv1.RemoveFromFavoritesReq) (*dealv1.RemoveFromFavoritesRes, error) {
userID, err := uuid.Parse(req.UserId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный user_id")
}
listingID, err := uuid.Parse(req.ListingId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный listing_id")
}

if err := h.favoriteUC.RemoveFromFavorites(ctx, userID, listingID); err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

return &dealv1.RemoveFromFavoritesRes{Success: true}, nil
}

func (h *DealHandler) ListFavorites(ctx context.Context, req *dealv1.ListFavoritesReq) (*dealv1.ListFavoritesRes, error) {
userID, err := uuid.Parse(req.UserId)
if err != nil {
return nil, status.Errorf(codes.InvalidArgument, "неверный user_id")
}

favorites, total, err := h.favoriteUC.ListFavorites(ctx, userID, int(req.Page), int(req.Limit))
if err != nil {
return nil, status.Errorf(codes.Internal, err.Error())
}

var protoFavorites []*dealv1.Favorite
for _, f := range favorites {
protoFavorites = append(protoFavorites, toProtoFavorite(f))
}

return &dealv1.ListFavoritesRes{Favorites: protoFavorites, Total: int32(total)}, nil
}

func protoToDealStatus(s dealv1.DealStatus) entity.DealStatus {
switch s {
case dealv1.DealStatus_DEAL_STATUS_REVIEW:
return entity.DealStatusReview
case dealv1.DealStatus_DEAL_STATUS_APPROVED:
return entity.DealStatusApproved
case dealv1.DealStatus_DEAL_STATUS_CLOSED:
return entity.DealStatusClosed
case dealv1.DealStatus_DEAL_STATUS_REJECTED:
return entity.DealStatusRejected
default:
return entity.DealStatusNew
}
}

func protoToViewingStatus(s dealv1.ViewingStatus) entity.ViewingStatus {
switch s {
case dealv1.ViewingStatus_VIEWING_STATUS_CONFIRMED:
return entity.ViewingStatusConfirmed
case dealv1.ViewingStatus_VIEWING_STATUS_CANCELLED:
return entity.ViewingStatusCancelled
case dealv1.ViewingStatus_VIEWING_STATUS_COMPLETED:
return entity.ViewingStatusCompleted
default:
return entity.ViewingStatusPending
}
}

func toProtoDeal(d *entity.Deal) *dealv1.Deal {
proto := &dealv1.Deal{
Id:          d.ID.String(),
ListingId:   d.ListingID.String(),
BuyerId:     d.BuyerID.String(),
PriceAgreed: d.PriceAgreed,
Currency:    d.Currency,
}
if d.AgentID != nil {
proto.AgentId = d.AgentID.String()
}
switch d.Status {
case entity.DealStatusNew:
proto.Status = dealv1.DealStatus_DEAL_STATUS_NEW
case entity.DealStatusReview:
proto.Status = dealv1.DealStatus_DEAL_STATUS_REVIEW
case entity.DealStatusApproved:
proto.Status = dealv1.DealStatus_DEAL_STATUS_APPROVED
case entity.DealStatusClosed:
proto.Status = dealv1.DealStatus_DEAL_STATUS_CLOSED
case entity.DealStatusRejected:
proto.Status = dealv1.DealStatus_DEAL_STATUS_REJECTED
}
return proto
}

func toProtoViewing(v *entity.ViewingRequest) *dealv1.ViewingRequest {
proto := &dealv1.ViewingRequest{
Id:        v.ID.String(),
ListingId: v.ListingID.String(),
BuyerId:   v.BuyerID.String(),
Comment:   v.Comment,
}
switch v.Status {
case entity.ViewingStatusPending:
proto.Status = dealv1.ViewingStatus_VIEWING_STATUS_PENDING
case entity.ViewingStatusConfirmed:
proto.Status = dealv1.ViewingStatus_VIEWING_STATUS_CONFIRMED
case entity.ViewingStatusCancelled:
proto.Status = dealv1.ViewingStatus_VIEWING_STATUS_CANCELLED
case entity.ViewingStatusCompleted:
proto.Status = dealv1.ViewingStatus_VIEWING_STATUS_COMPLETED
}
return proto
}

func toProtoFavorite(f *entity.Favorite) *dealv1.Favorite {
return &dealv1.Favorite{
Id:        f.ID.String(),
UserId:    f.UserID.String(),
ListingId: f.ListingID.String(),
}
}
