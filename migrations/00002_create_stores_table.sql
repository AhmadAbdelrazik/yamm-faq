-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stores (
  id SERIAL PRIMARY KEY,
  merchant_id INT NOT NULL REFERENCES users(id),
  name VARCHAR(30) NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stores;
-- +goose StatementEnd
