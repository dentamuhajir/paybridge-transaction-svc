-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS ledger_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL,
    account_id UUID NOT NULL,
    balance_type_id INT NOT NULL,
    amount BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    event_type VARCHAR(50) NOT NULL,
    reference_type VARCHAR(50),
    reference_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_ledger_transaction
        FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    CONSTRAINT fk_ledger_account
        FOREIGN KEY (account_id) REFERENCES accounts(id),
    CONSTRAINT fk_ledger_balance_type
        FOREIGN KEY (balance_type_id) REFERENCES balance_types(id),
    CONSTRAINT chk_ledger_amount_nonzero
        CHECK (amount <> 0)
);

CREATE INDEX idx_ledger_account_balance
    ON ledger_entries(account_id, balance_type_id);

CREATE INDEX idx_ledger_transaction
    ON ledger_entries(transaction_id);

CREATE INDEX idx_ledger_reference
    ON ledger_entries(reference_type, reference_id);

-- +migrate Down
DROP TABLE IF EXISTS ledger_entries;