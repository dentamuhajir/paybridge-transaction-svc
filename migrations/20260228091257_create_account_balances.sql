-- +migrate Up
CREATE TABLE account_balances (
    account_id UUID PRIMARY KEY,
    balance BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (account_id) REFERENCES accounts(id)
);