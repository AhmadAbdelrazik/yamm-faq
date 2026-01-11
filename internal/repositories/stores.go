package repositories

import (
	"database/sql"
	"errors"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
)

// StoresRepository Implement Retrieval and Updating of Stores. Creation is
// handled by UsersRepository merchant creation method. Deletion is also
// handled by the UsersRepository since there can't be a merchant Account with
// no store associated with it.
type StoresRepository struct {
	db *sql.DB
}

func (r *StoresRepository) FindByID(id int) (*models.Store, error) {
	query := `
	SELECT merchant_id, name 
	FROM stores 
	WHERE id = $1 AND deleted_at IS NULL`

	store := &models.Store{
		ID: id,
	}

	err := r.db.QueryRow(query, id).Scan(&store.MerchantID, &store.Name)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return store, nil
}

func (r *StoresRepository) FindByMerchantID(merchantID int) (*models.Store, error) {
	query := `
	SELECT id, name 
	FROM stores 
	WHERE merchant_id = $1 AND deleted_at IS NULL`

	store := &models.Store{
		MerchantID: merchantID,
	}

	err := r.db.QueryRow(query, merchantID).Scan(&store.ID, &store.Name)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return store, nil
}

func (r *StoresRepository) Update(store *models.Store) error {
	query := `
	UPDATE stores
	SET name = $1, updated_at = NOW()
	WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, store.ID, store.Name)
	if err != nil {
		return err
	}

	if n, err := result.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrNotFound
	}

	return nil
}
