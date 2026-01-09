package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
	INSERT INTO users(email, password, role)
	VALUES($1, $2, $3)
	RETURNING id`

	args := []any{user.Email, user.Password.Hash}
	if err := r.db.QueryRow(query, args...).Scan(&user.ID); err != nil {
		switch {
		case strings.Contains(err.Error(), "users_email_key"):
			return fmt.Errorf("%w: user with this email already exists", ErrDuplicate)
		default:
			return err
		}
	}

	return nil
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	query := `SELECT email, role, hash WHERE id = $1 AND deleted_at IS NULL`

	user := &models.User{
		ID:       id,
		Password: &models.Password{},
	}

	err := r.db.QueryRow(query, id).Scan(
		&user.Email,
		&user.Role,
		&user.Password.Hash,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `SELECT id, role, hash WHERE email = $1 AND deleted_at IS NULL`

	user := &models.User{
		Email:    email,
		Password: &models.Password{},
	}

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Role,
		&user.Password.Hash,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
