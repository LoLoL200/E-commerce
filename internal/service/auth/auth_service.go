package auth

import (
	"context"
	models "ecommers/internal/domain"
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

// SERVICE
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

// REGISTER
func (a *authService) Register(ctx context.Context, email, password string) (uuid.UUID, error) {
	user, err := a.userRepo.GetEmail(ctx, email)
	if err != nil && err != utils.ErrUserNotfound {
		return uuid.Nil, err
	}

	if user != nil {
		return uuid.Nil, utils.ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	newUser := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		Role:         "customer",
		CreateAt:     time.Now(),
	}

	if err := a.userRepo.Create(ctx, newUser); err != nil {
		return uuid.Nil, err
	}

	return newUser.ID, nil
}

// LOGIN
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

// REFRESH TOKEN
func (a *authService) Refresh(ctx context.Context, refreshToken string) (string, error) {

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

	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	return a.generateAccessToken(user)
}

// Validate token
func (a *authService) ValidateToken(tokenString string) (*TokenClaims, error) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(a.secret), nil
	})

	if err != nil {
		log.Printf("🔴 JWT Parsing Error: %v", err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, utils.ErrInvalidCredentials
	}

	userID, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	return &TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
	}, nil
}

// Create auth service
func NewAuthService(repo repository.UserRepository, secret string) *authService {
	return &authService{
		userRepo: repo,
		secret:   secret,
	}
}

// Generate Access token
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

// Generate Refresh Token
func (a *authService) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secret))
}
