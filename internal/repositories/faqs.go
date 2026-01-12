package repositories

import (
	"database/sql"
	"errors"

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

	query = `
	INSERT INTO translations(faq_id, language, question, answer)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
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

// FindDefault Returns FAQ with its default Language
func (r *FAQRepository) FindDefault(id int) (*models.FAQ, error) {
	query := `
	SELECT f.id, f.category, f.default_language, f.store_id, t.id, t.faq_id, t.language,
	t.question, t.answer
	FROM faqs AS f
	JOIN translations AS t ON t.faq_id = f.id AND t.language = f.default_language
	WHERE f.id = $1`

	var c models.FAQCategory
	var t models.Translation
	var storeID sql.NullInt32

	faq := &models.FAQ{ID: id}

	err := r.db.QueryRow(query, id).Scan(
		&faq.ID,
		&c.Name,
		&faq.DefaultLanguage,
		&storeID,
		&t.ID,
		&t.FAQID,
		&t.Language,
		&t.Question,
		&t.Answer,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	faq.Category = c
	faq.Translations = append(faq.Translations, t)

	if storeID.Valid {
		faq.StoreID = int(storeID.Int32)
	}

	return faq, nil
}

// Find Returns FAQ with all available translations
func (r *FAQRepository) Find(id int) (*models.FAQ, error) {
	query := `
	SELECT f.id, f.category, f.default_language, f.store_id, t.id, t.faq_id, t.language,
	t.question, t.answer
	FROM faqs AS f
	JOIN translations AS t ON t.faq_id = f.id
	WHERE f.id = $2`

	faq := &models.FAQ{ID: id}
	var storeID sql.NullInt32

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t models.Translation
		err := rows.Scan(
			&faq.ID,
			&faq.Category.Name,
			&faq.DefaultLanguage,
			&storeID,
			&t.ID,
			&t.FAQID,
			&t.Language,
			&t.Question,
			&t.Answer,
		)
		if err != nil {
			return nil, err
		}

		faq.Translations = append(faq.Translations, t)
	}

	if storeID.Valid {
		faq.StoreID = int(storeID.Int32)
	}

	return faq, nil
}

// GetAllByCategory Returns all FAQs that belongs to category. This includes
// global and store-specific FAQs
func (r *FAQRepository) GetAllByCategory(category string) ([]models.FAQ, error) {
	query := `
	SELECT f.id, f.category, f.default_language, t.id, t.faq_id, t.language,
	t.question, t.answer
	FROM faqs AS f
	JOIN translations AS t ON t.faq_id = f.id AND t.language = f.default_language
	WHERE f.category = $1`

	faqs := make([]models.FAQ, 0)

	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c models.FAQCategory
		var f models.FAQ
		var t models.Translation

		err := rows.Scan(
			&f.ID,
			&c.Name,
			&f.DefaultLanguage,
			&t.ID,
			&t.FAQID,
			&t.Language,
			&t.Question,
			&t.Answer,
		)

		if err != nil {
			return nil, err
		}

		f.Category = c
		f.Translations = append(f.Translations, t)

		faqs = append(faqs, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return faqs, nil
}

// GetAllByStore Returns all FAQs that belongs to specific store.
func (r *FAQRepository) GetAllByStore(storeID int) ([]models.FAQ, error) {
	query := `
	SELECT f.id, f.category, f.default_language, t.id, t.faq_id, t.language,
	t.question, t.answer
	FROM faqs AS f
	JOIN translations AS t ON t.faq_id = f.id AND t.language = f.default_language
	WHERE f.store_id = $1`

	faqs := make([]models.FAQ, 0)

	rows, err := r.db.Query(query, storeID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c models.FAQCategory
		var f models.FAQ
		var t models.Translation

		err := rows.Scan(
			&f.ID,
			&c.Name,
			&f.DefaultLanguage,
			&t.ID,
			&t.FAQID,
			&t.Language,
			&t.Question,
			&t.Answer,
		)

		if err != nil {
			return nil, err
		}

		f.Category = c
		f.Translations = append(f.Translations, t)

		faqs = append(faqs, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return faqs, nil
}

func (r *FAQRepository) Update(faq *models.FAQ) error {
	query := `
	UPDATE faqs
	SET category = $1, default_language = $2, is_global = $3
	`
	args := []any{
		faq.Category.Name,
		faq.DefaultLanguage,
		faq.IsGlobal,
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); err != nil || n == 0 {
		return err
	}

	return nil
}

// SetGlobal Sets FAQ to global
func (r *FAQRepository) SetGlobal(faqID int) error {
	query := `UPDATE faqs SET is_global = TRUE WHERE id = $1`

	res, err := r.db.Exec(query, faqID)
	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); err != nil || n == 0 {
		return err
	}

	return nil
}

func (r *FAQRepository) Delete(faqID int) error {
	query := `DELETE FROM faqs WHERE id = $1`

	result, err := r.db.Exec(query, faqID)
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

// MerchantDelete Delete FAQ with ID and Store ID. Used for merchants
func (r *FAQRepository) MerchantDelete(faqID, storeID int) error {
	query := `DELETE FROM faqs WHERE id = $1 AND store_id = $2`

	result, err := r.db.Exec(query, faqID, storeID)
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
