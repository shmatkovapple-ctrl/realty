package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"realty/services/user-service/internal/domain/entity"
	"realty/services/user-service/internal/domain/repository"
	redisInfra "realty/services/user-service/internal/infrastructure/redis"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

type userUseCase struct {
	userRepo   repository.UserRepository
	tokenCache *redisInfra.TokenCache
	jwtSecret  string
}

func NewUserUseCase(
	userRepo repository.UserRepository,
	tokenCache *redisInfra.TokenCache,
	jwtSecret string,
) UserUseCase {
	return &userUseCase{
		userRepo:   userRepo,
		tokenCache: tokenCache,
		jwtSecret:  jwtSecret,
	}
}

func (uc *userUseCase) Register(ctx context.Context, input RegisterInput) (*AuthOutput, error) {
	existing, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err == nil && existing != nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("хэширование пароля: %w", err)
	}

	user, err := entity.NewUser(input.Email, string(hash), input.Phone, entity.Role(input.Role))
	if err != nil {
		return nil, fmt.Errorf("создание пользователя: %w", err)
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("сохранение пользователя: %w", err)
	}

	return uc.buildAuthOutput(ctx, user)
}

func (uc *userUseCase) Login(ctx context.Context, input LoginInput) (*AuthOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	if !user.IsActive() {
		return nil, errors.New("аккаунт заблокирован или не активирован")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	return uc.buildAuthOutput(ctx, user)
}

func (uc *userUseCase) Logout(ctx context.Context, userID, refreshToken string) error {
	return uc.tokenCache.DeleteRefreshToken(ctx, userID, refreshToken)
}

func (uc *userUseCase) RefreshToken(ctx context.Context, refreshToken string) (*AuthOutput, error) {
	claims, err := uc.parseToken(refreshToken)
	if err != nil {
		return nil, errors.New("недействительный refresh токен")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("недействительный идентификатор пользователя")
	}

	valid, err := uc.tokenCache.ValidateRefreshToken(ctx, claims.UserID, refreshToken)
	if err != nil || !valid {
		return nil, errors.New("refresh токен не найден или истёк")
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}

	if err := uc.tokenCache.DeleteRefreshToken(ctx, claims.UserID, refreshToken); err != nil {
		return nil, fmt.Errorf("удаление старого токена: %w", err)
	}

	return uc.buildAuthOutput(ctx, user)
}

func (uc *userUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}
	return user, nil
}

func (uc *userUseCase) UpdateProfile(ctx context.Context, input UpdateProfileInput) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("пользователь не найден: %w", err)
	}

	user.Profile.FirstName = input.FirstName
	user.Profile.LastName = input.LastName
	user.Profile.Bio = input.Bio

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("обновление профиля: %w", err)
	}

	return user, nil
}

func (uc *userUseCase) ValidateToken(ctx context.Context, accessToken string) (*TokenClaims, error) {
	claims, err := uc.parseToken(accessToken)
	if err != nil {
		return nil, errors.New("недействительный токен")
	}
	return claims, nil
}

func (uc *userUseCase) buildAuthOutput(ctx context.Context, user *entity.User) (*AuthOutput, error) {
	accessToken, err := uc.generateToken(user, accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("генерация access токена: %w", err)
	}

	refreshToken, err := uc.generateToken(user, refreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("генерация refresh токена: %w", err)
	}

	if err := uc.tokenCache.SaveRefreshToken(ctx, user.ID.String(), refreshToken, refreshTokenTTL); err != nil {
		return nil, fmt.Errorf("сохранение refresh токена: %w", err)
	}

	return &AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (uc *userUseCase) generateToken(user *entity.User, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    string(user.Role),
		"exp":     time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

func (uc *userUseCase) parseToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", t.Header["alg"])
		}
		return []byte(uc.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("недействительный токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("неверный формат claims")
	}

	return &TokenClaims{
		UserID: claims["user_id"].(string),
		Role:   claims["role"].(string),
	}, nil
}
