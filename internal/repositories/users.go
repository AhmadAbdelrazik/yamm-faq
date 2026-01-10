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

// Create creates a new user with no regard of the user type. For merchants you
// should use CreateMerchant instead
func (r *UserRepository) Create(user *models.User) error {
	query := `
	INSERT INTO users(email, hash, role)
	VALUES($1, $2, $3)
	RETURNING id`

	args := []any{user.Email, user.Password.Hash, string(user.Role)}
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

// Create creates a new merchant and creates a relative store for the merchant.
func (r *UserRepository) CreateMerchant(merchant *models.User, store *models.Store) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := `
	INSERT INTO users(email, hash, role)
	VALUES($1, $2, 'merchant')
	RETURNING id`
	args := []any{merchant.Email, merchant.Password.Hash}

	if err := tx.QueryRow(query, args...).Scan(&merchant.ID); err != nil {
		tx.Rollback()
		switch {
		case strings.Contains(err.Error(), "users_email_key"):
			return fmt.Errorf("%w: user with this email already exists", ErrDuplicate)
		default:
			return err
		}
	}

	store.MerchantID = merchant.ID

	query = `INSERT INTO stores(merchant_id, name) VALUES ($1, $2) RETURNING id`
	args = []any{store.MerchantID, store.Name}

	if err := tx.QueryRow(query, args...).Scan(&store.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	query := `
	SELECT email, role, hash 
	FROM users 
	WHERE id = $1 AND deleted_at IS NULL`

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
	query := `
	SELECT id, role, hash 
	FROM users 
	WHERE email = $1 AND deleted_at IS NULL`

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
