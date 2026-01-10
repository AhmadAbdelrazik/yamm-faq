-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS faq_categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(20) NOT NULL UNIQUE,

  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS faqs (
  id SERIAL PRIMARY KEY,
  category VARCHAR(20) NOT NULL REFERENCES faq_categories(name) ON DELETE CASCADE,
  default_language VARCHAR(2) NOT NULL,
  store_id INT REFERENCES stores(id),
  is_global BOOLEAN NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE INDEX faqs_category_idx ON faqs(category);
CREATE INDEX faqs_store_id_idx ON faqs(store_id);
CREATE INDEX faqs_language_idx ON faqs(default_language);

CREATE TABLE IF NOT EXISTS translations (
  id SERIAL PRIMARY KEY,
  faq_id INT NOT NULL REFERENCES faqs(id) ON DELETE CASCADE,
  language VARCHAR(2) NOT NULL,
  question TEXT NOT NULL,
  answer TEXT NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE INDEX translations_faq_id_idx ON translations(faq_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS translations;
DROP TABLE IF EXISTS faqs;
DROP TABLE IF EXISTS faq_categories;
-- +goose StatementEnd
