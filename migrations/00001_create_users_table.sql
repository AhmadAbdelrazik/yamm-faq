-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(200) UNIQUE NOT NULL,
  role VARCHAR(10) NOT NULL,
  hash BYTEA NOT NULL,
  
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- admin user with password: admin123
INSERT INTO users(email, role, hash) VALUES('admin@test.com', 'admin',
  '\x24326124313024716969504a792e312f36565547386847794b664a792e6731553433336e5135513479736f396336514c6431704f78766d55777a6553')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
