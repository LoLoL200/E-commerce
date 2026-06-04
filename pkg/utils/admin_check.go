package utils

import (
	models "ecommers/internal/domain"
	"fmt"
	"log"
	"net/http"
)

func AdminCheck(r *http.Request, user *models.User) error {
	if user == nil {
		return fmt.Errorf("unauthorized: user data is missing")
	}

	if user.Role != "admin" {
		log.Printf("[SECURITY WARNING] User ID %s (Email: %s) attempted to access admin route %s from IP %s",
			user.ID, user.Email, r.URL.Path, r.RemoteAddr)

		return fmt.Errorf("access denied: user %s does not have administrator privileges", user.ID)
	}

	return nil
}
