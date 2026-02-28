-- +migrate Up
CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    transaction_id UUID NOT NULL,
    account_id UUID NOT NULL,

    amount BIGINT NOT NULL, 
    -- signed: + credit, - debit

    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_ledger_transaction
        FOREIGN KEY (transaction_id) REFERENCES transactions(id),

    CONSTRAINT fk_ledger_account
        FOREIGN KEY (account_id) REFERENCES accounts(id),

    CONSTRAINT chk_amount_nonzero
        CHECK (amount <> 0)
);

CREATE INDEX idx_ledger_account
ON ledger_entries(account_id);

CREATE INDEX idx_ledger_transaction
ON ledger_entries(transaction_id);