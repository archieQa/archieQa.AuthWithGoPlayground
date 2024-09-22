-- Create tokens table
CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    token_string TEXT NOT NULL UNIQUE,
    expiration_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE
);

-- Create index on user_id for faster lookups
CREATE INDEX idx_tokens_user_id ON tokens(user_id);

-- Create index on token_string for faster token validation
CREATE INDEX idx_tokens_token_string ON tokens(token_string);

-- Create index on expiration_date for efficient cleanup of expired tokens
CREATE INDEX idx_tokens_expiration_date ON tokens(expiration_date);
