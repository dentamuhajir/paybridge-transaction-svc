-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL,
    transaction_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    total_amount BIGINT NOT NULL,
    reference_type VARCHAR(50),
    reference_id UUID,
    idempotency_key VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_transactions_account
        FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE INDEX idx_transactions_account_id
    ON transactions(account_id);

CREATE INDEX idx_transactions_reference
    ON transactions(reference_type, reference_id);

-- +migrate Down
DROP TABLE IF EXISTS transactions;