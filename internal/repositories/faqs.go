package repositories

import (
	"database/sql"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
)

type FAQRepository struct {
	db *sql.DB
}

func (r *FAQRepository) Create(faq *models.FAQ) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var query string
	var args []any

	if faq.IsGlobal {
		query = `
	INSERT INTO faqs(category, default_language, is_global)
	VALUES ($1, $2, $3) RETURNING id`
		args = []any{
			faq.Category.Name,
			faq.DefaultLanguage,
			faq.IsGlobal,
		}
	} else {
		query = `
	INSERT INTO faqs(category, default_language, store_id, is_global)
	VALUES ($1, $2, $3, $4) RETURNING id`
		args = []any{
			faq.Category.Name,
			faq.DefaultLanguage,
			faq.StoreID,
			faq.IsGlobal,
		}

	}

	if err := tx.QueryRow(query, args...).Scan(&faq.ID); err != nil {
		tx.Rollback()
		return err
	}

	query = `INSERT INTO translations(faq_id, language, question, answer)
	VALUES ($1, $2, $3, $4)`
	args = []any{
		faq.ID,
		faq.Translations[0].Language,
		faq.Translations[0].Question,
		faq.Translations[0].Answer,
	}

	if err := tx.QueryRow(query, args...).Scan(&faq.Translations[0].ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

