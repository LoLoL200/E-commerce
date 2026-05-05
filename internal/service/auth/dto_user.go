package auth

import models "ecommers/internal/domain"

type UserFillter struct {
	Search  string
	Role    *string
	Limit   int
	Offset  int
	OrderBy string
}

type UserListResponce struct {
	Users  []*models.User
	Total  int
	Limit  int
	Offset int
}

type UpdateProfileRequest struct {
	Email     *string `json:"email"      validate:"omitempty,email"`
	Password  *string `json:"password"   validate:"omitempty,min=6"`
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=256"`
	LastName  *string `json:"last_name"  validate:"omitempty,min=2,max=256"`
	Role      *string `json:"role"       validate:"omitempty,oneof=admin user"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterRequest struct {
	Email     string `json:"email"      validate:"required,email"`
	Password  string `json:"password"   validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=2,max=256"`
	Surname   string `json:"surname"    validate:"required,min=2,max=256"`
}
