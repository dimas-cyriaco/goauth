CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL,
    user_agent VARCHAR(255),
    ip_address VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sessions_accounts FOREIGN KEY (account_id)
    REFERENCES accounts (id) ON DELETE CASCADE
);
