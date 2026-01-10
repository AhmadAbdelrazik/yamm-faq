package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
)

type FaqCategoryRepository struct {
	db *sql.DB
}

func (r *FaqCategoryRepository) Create(category *models.FAQCategory) error {
	query := `INSERT INTO faq_categories(name) VALUES ($1) RETURNING id`

	if err := r.db.QueryRow(query, strings.ToLower(category.Name)).Scan(&category.ID); err != nil {
		switch {
		case strings.Contains(err.Error(), "faq_categories_name_key"):
			return fmt.Errorf("%w: category already exists", ErrDuplicate)
		default:
			return err
		}
	}

	return nil
}

func (r *FaqCategoryRepository) GetAll() ([]models.FAQCategory, error) {
	query := `SELECT id, name FROM faq_categories WHERE deleted_at IS NULL`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.FAQCategory{}

	for rows.Next() {
		var c models.FAQCategory
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return categories, nil
}

func (r *FaqCategoryRepository) FindByID(id int) (*models.FAQCategory, error) {
	query := `
	SELECT name 
	FROM faq_categories 
	WHERE id = $1 AND deleted_at IS NULL`

	category := &models.FAQCategory{ID: id}

	err := r.db.QueryRow(query, category.ID).Scan(&category.Name)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w: category doesn't exist", ErrNotFound)
		default:
			return nil, err
		}
	}

	return category, nil
}

func (r *FaqCategoryRepository) FindByCategoryName(name string) (*models.FAQCategory, error) {
	query := `
	SELECT id 
	FROM faq_categories 
	WHERE name = $1 AND deleted_at IS NULL`

	category := &models.FAQCategory{Name: name}

	err := r.db.QueryRow(query, category.Name).Scan(&category.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w: category doesn't exist", ErrNotFound)
		default:
			return nil, err
		}
	}

	return category, nil
}

func (r *FaqCategoryRepository) Update(category *models.FAQCategory) error {
	query := `
	UPDATE faq_categories
	SET name = $1, updated_at = NOW() 
	WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, category.Name, category.ID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "faq_categories_name_key"):
			return fmt.Errorf("%w: category already exists", ErrDuplicate)
		default:
			return err
		}
	}

	if n, err := result.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrEditConflict
	}

	return nil
}

// Delete soft deletes categories and any FAQ under it.
func (r *FaqCategoryRepository) Delete(categoryName string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `
	UPDATE faq_categories SET deleted_at = NOW()
	WHERE name = $1 AND deleted_at IS NULL`

	result, err := tx.Exec(query, categoryName)
	if err != nil {
		tx.Rollback()
		return err
	}

	if n, err := result.RowsAffected(); err != nil {
		tx.Rollback()
		return err
	} else if n == 0 {
		return ErrNotFound
	}

	query = `
	UPDATE faqs SET deleted_at = NOW()
	WHERE category = $1 AND deleted_at IS NULL`

	result, err = tx.Exec(query, categoryName)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// HardDelete Hard deletes FAQ Category and any associated FAQs
func (r *FaqCategoryRepository) HardDelete(categoryName string) error {
	query := `DELETE FROM faq_categories WHERE name = $1`

	result, err := r.db.Exec(query, categoryName)
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
