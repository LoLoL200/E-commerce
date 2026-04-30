package auth

import (
	"context"
	models "ecommers/internal/domin"
	repository "ecommers/internal/repository/postgres"
	"ecommers/pkg/utils"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}
type AuthService interface {
	Register(ctx context.Context, email, password string) (uuid.UUID, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
	ValidateToken(token string) (*TokenClaims, error)
}
type authService struct {
	userRepo repository.UserRepository
	secret   string
}

func (a *authService) Register(ctx context.Context, email, password string) (uuid.UUID, error) {
	user, err := a.userRepo.GetEmail(ctx, email)
	if err != nil && err != utils.ErrUserNotfound {
		return uuid.Nil, err
	}

	if user != nil {
		return uuid.Nil, utils.ErrUserAlreadyExists
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	// CREATE USER (вот этого у тебя не было)
	newUser := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		Role:         "customer",
		CreateAt:     time.Now(),
	}

	// SAVE TO DB
	if err := a.userRepo.Create(ctx, newUser); err != nil {
		return uuid.Nil, err
	}

	return newUser.ID, nil
}
func (a *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := a.userRepo.GetEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", utils.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", utils.ErrInvalidCredentials
	}

	access, err := a.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refresh, err := a.generateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (a *authService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	// 1. валидируем refresh token
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		return []byte(a.secret), nil
	})
	if err != nil || !token.Valid {
		log.Printf("❌ refresh token invalid: %v", err)
		return "", utils.ErrInvalidCredentials
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", utils.ErrInvalidCredentials
	}

	userIDStr, _ := claims["user_id"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", err
	}

	// 2. достаём пользователя (лучше из БД)
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	// 3. выдаём новый access token
	return a.generateAccessToken(user)
}

//	func (a *authService) ValidateToken(token string) (*TokenClaims, error) {
//		panic("implement me")
//	}
func (a *authService) ValidateToken(tokenString string) (*TokenClaims, error) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(a.secret), nil
	})

	if err != nil {
		log.Printf("🔴 JWT Parsing Error: %v", err)
		return nil, err // Здесь вернется "token is expired" или "signature is invalid"
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, utils.ErrInvalidCredentials
	}

	// Извлекаем данные точно по тем ключам, которые использовали при генерации
	userID, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	return &TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
	}, nil
}

func NewAuthService(repo repository.UserRepository, secret string) *authService {
	return &authService{
		userRepo: repo,
		secret:   secret,
	}
}
func (a *authService) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secret))
}
func (a *authService) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secret))
}
