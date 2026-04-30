package auth

import (
	"context"
	models "ecommers/internal/domin"
	repository "ecommers/internal/repository/postgres"
	"ecommers/pkg/utils"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetProfile(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, request UpdateProfileRequest, isAdmin bool) (*models.User, error)
	ListUser(ctx context.Context, fillter UserFillter) (*UserListResponce, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) UserService {
	return &service{userRepo: userRepo}
}

// Delete implements [UserService].
func (s *service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetProfile implements [UserService].
func (s *service) GetProfile(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, utils.ErrUserNotfound
	}
	return user, nil
}

// ListUser implements [UserService].
func (s *service) ListUser(ctx context.Context, fillter UserFillter) (*UserListResponce, error) {
	if fillter.Limit <= 0 {
		fillter.Limit = 20
	}
	if fillter.Limit > 100 {
		fillter.Limit = 100
	}
	if fillter.Offset < 0 {
		fillter.Offset = 0
	}
	users, err := s.userRepo.List(ctx, fillter.Limit, fillter.Offset)
	if err != nil {
		return nil, fmt.Errorf("filead to list users: %w", err)
	}
	fillted := []*models.User{}
	for _, u := range users {
		if fillter.Role != nil && u.Role != *fillter.Role {
			continue
		}
		if fillter.Search != "" && !strings.Contains(
			strings.ToLower(u.FirstName),
			strings.ToLower(fillter.Search)) && !strings.Contains(
			strings.ToLower(u.Surname),
			strings.ToLower(fillter.Search)) {
			continue
		}

		fillted = append(fillted, u)

	}
	resp := &UserListResponce{
		Users: fillted,
		Total: len(fillted),
	}
	return resp, nil
}

// UpdateProfile implements [UserService].
func (s *service) UpdateProfile(ctx context.Context, id uuid.UUID, request UpdateProfileRequest, isAdmin bool) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to go user: %w", err)
	}
	if request.FirstName != nil {
		user.FirstName = *request.FirstName
	}
	if request.LastName != nil {
		user.Surname = *request.LastName
	}
	if request.Email != nil {
		user.Email = *request.Email
	}
	if request.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash: %w", err)
		}
		user.PasswordHash = string(hash)
	}
	if request.Role != nil {
		user.Role = *request.Role

	}
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("filed to update user: %w", err)
	}
	return user, nil
}
