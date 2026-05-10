package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	userv1 "realty/api/gen/user/v1"
	"realty/services/user-service/internal/domain/entity"
	"realty/services/user-service/internal/usecase"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	uc usecase.UserUseCase
}

func NewUserHandler(uc usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	out, err := h.uc.Register(ctx, usecase.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Role:     req.Role,
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &userv1.RegisterResponse{
		UserId:       out.User.ID.String(),
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	out, err := h.uc.Login(ctx, usecase.LoginInput{
		Email:      req.Email,
		Password:   req.Password,
		DeviceInfo: req.DeviceInfo,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &userv1.LoginResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		Profile:      toProtoProfile(out.User),
	}, nil
}

func (h *UserHandler) Logout(ctx context.Context, req *userv1.LogoutRequest) (*userv1.LogoutResponse, error) {
	claims, err := h.uc.ValidateToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "недействительный токен")
	}
	if err := h.uc.Logout(ctx, claims.UserID, req.RefreshToken); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &userv1.LogoutResponse{Success: true}, nil
}

func (h *UserHandler) ValidateToken(ctx context.Context, req *userv1.ValidateTokenRequest) (*userv1.ValidateTokenResponse, error) {
	claims, err := h.uc.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return &userv1.ValidateTokenResponse{Valid: false}, nil
	}
	return &userv1.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Role:   claims.Role,
	}, nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *userv1.GetProfileRequest) (*userv1.GetProfileResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "неверный формат user_id")
	}
	user, err := h.uc.GetProfile(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &userv1.GetProfileResponse{Profile: toProtoProfile(user)}, nil
}

func (h *UserHandler) UpdateProfile(ctx context.Context, req *userv1.UpdateProfileRequest) (*userv1.UpdateProfileResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "неверный формат user_id")
	}
	user, err := h.uc.UpdateProfile(ctx, usecase.UpdateProfileInput{
		UserID:    userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &userv1.UpdateProfileResponse{Profile: toProtoProfile(user)}, nil
}

func (h *UserHandler) RefreshToken(ctx context.Context, req *userv1.RefreshTokenRequest) (*userv1.RefreshTokenResponse, error) {
	out, err := h.uc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &userv1.RefreshTokenResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}, nil
}

func toProtoProfile(u *entity.User) *userv1.UserProfile {
	p := &userv1.UserProfile{
		UserId: u.ID.String(),
		Email:  u.Email,
		Role:   string(u.Role),
		Status: string(u.Status),
	}
	if u.Profile != nil {
		p.FirstName = u.Profile.FirstName
		p.LastName = u.Profile.LastName
		p.AvatarUrl = u.Profile.AvatarURL
	}
	return p
}
