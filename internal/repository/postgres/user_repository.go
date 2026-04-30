package repository

import (
	"context"
	"database/sql"
	models "ecommers/internal/domin"
	"errors"

	//models "e_commerce/internal/domin"
	"fmt"

	database "ecommers/internal/db"

	//models

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
}
type userRepo struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepo{db: db}
}

// Create implements [UserRepository].
func (u *userRepo) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users(id,email,password_hash,first_name,surname,role)
		VALUES($1,$2,$3,$4,$5,$6)
		RETURNING id, create_at, updated_at
	`
	user.ID = uuid.New()
	return u.db.QueryRowContext(
		ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.Surname, user.Role,
	).Scan(&user.ID, &user.CreateAt, &user.UpdateAt)
}

// Delete implements [UserRepository].
func (u *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users WHERE id=$1
	`
	result, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("filed to delete user : %w", err)
	}
	row, err := result.RowsAffected()
	if row == 0 {
		return fmt.Errorf("user not found ")
	}
	return nil
}

// GetByID implements [UserRepository].
func (u *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User

	query := `
	SELECT id,email,password_hash,first_name,surname,role
	FROM users
	WHERE id = $1
	`
	err := u.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("filed to get user by id:%w", err)
	}
	return &user, nil
}

// GetEmail implements [UserRepository].
func (u *userRepo) GetEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `
        SELECT id, email, password_hash, first_name, surname, role
        FROM users
        WHERE email = $1
    `
	err := u.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // пользователь не найден — это не ошибка
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// List implements [UserRepository].
func (u *userRepo) List(ctx context.Context, limit int, offset int) ([]*models.User, error) {
	query := `
	SELECT id,email,password_hash,first_name,surname,role
	FROM users
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2

	`
	var users []*models.User

	err := u.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("filed to list userd: %w", err)
	}
	return users, nil
}

// Update implements [UserRepository].
func (u *userRepo) Update(ctx context.Context, user *models.User) error {
	query := `
	UPDATE users
	SET first_name = $1, surname = $2, update_at = NOW()
	WHERE id = $3
	RETURNING update_at

	`

	return u.db.QueryRowContext(
		ctx, query,
		user.FirstName, user.Surname, user.ID,
	).Scan(&user.UpdateAt)
}
