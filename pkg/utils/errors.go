package utils

import (
	"errors"
)

// Errors
var (
	ErrUserNotfound       = errors.New("User not found")
	ErrUserAlreadyExists  = errors.New("User already exists")
	ErrInvalidEmail       = errors.New("Invalid email address")
	ErrEmailRequired      = errors.New("Please enter your email addres")
	ErrPasswordRequired   = errors.New("Please enter your password addres")
	ErrWeakPassword       = errors.New("The password is too simpless")
	ErrInvalidCredentials = errors.New("Invalid email or password")
	ErrInvalidToken       = errors.New("Invalid or malformed token")
	ErrTokenExpired       = errors.New("Token has expired")
	ErrCartNotFound       = errors.New("It`s CART not found")
	ErrItemNotFound       = errors.New("It`s ITEM not found")
)
