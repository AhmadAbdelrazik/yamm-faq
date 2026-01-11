package repositories

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
)

type TranslationsRepository struct {
	db *sql.DB
}

func (r *TranslationsRepository) Create(translation *models.Translation) error {
	query := `
	INSERT INTO translations(faq_id, language, question,
	answer)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	args := []any{
		translation.FAQID,
		translation.Language,
		translation.Question,
		translation.Answer,
	}

	err := r.db.QueryRow(query, args...).Scan(&translation.ID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "translations_faq_id_language_idx"):
			return ErrDuplicate
		default:
			return err
		}
	}

	return nil
}

// GetAll Returns all translations related to one FAQ
func (r *TranslationsRepository) GetAll(faqID int) ([]models.Translation, error) {
	query := `
	SELECT id, language, question, answer
	FROM translations
	WHERE faq_id = $1 AND deleted_at IS NULL
	`

	rows, err := r.db.Query(query, faqID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	translations := make([]models.Translation, 0)

	for rows.Next() {
		var t models.Translation
		err := rows.Scan(
			&t.ID,
			&t.Language,
			&t.Question,
			&t.Answer,
		)

		if err != nil {
			return nil, err
		}

		translations = append(translations, t)
	}

	return translations, nil
}

// Get Returns specific translation of an FAQ
func (r *TranslationsRepository) Get(faqID int, language models.Language) (*models.Translation, error) {
	query := `
	SELECT id, question, answer
	FROM translations
	WHERE faq_id = $1 AND language = $2 AND deleted_at IS NULL
	`

	translation := &models.Translation{
		FAQID:    faqID,
		Language: language,
	}

	err := r.db.QueryRow(query, faqID, language).Scan(
		&translation.ID,
		&translation.Question,
		&translation.Answer,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return translation, nil
}

func (r *TranslationsRepository) Update(translation *models.Translation) error {
	query := `
	UPDATE translations 
	SET language = $1, question = $2, answer = $3
	WHERE id = $4 AND deleted_at IS NULL`
	args := []any{
		translation.Language,
		translation.Question,
		translation.Answer,
		translation.ID,
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "translations_faq_id_language_idx"):
			return ErrDuplicate
		default:
			return err
		}
	}

	return nil
}

func (r *TranslationsRepository) Delete(id int) error {
	query := `DELETE FROM translations WHERE id = $1`

	result, err := r.db.Exec(query, id)
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
