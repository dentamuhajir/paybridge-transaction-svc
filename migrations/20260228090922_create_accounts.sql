-- +migrate Up
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    owner_type VARCHAR(20) NOT NULL, -- USER / SYSTEM
    owner_id UUID,                   -- null if SYSTEM

    account_code VARCHAR(50) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',

    reference_type VARCHAR(30),      -- e.g. LOAN
    reference_id UUID,               -- loan_id if applicable

    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_owner_type
        CHECK (owner_type IN ('USER','SYSTEM')),

    CONSTRAINT chk_owner_required
        CHECK (
            (owner_type = 'USER' AND owner_id IS NOT NULL)
            OR
            (owner_type = 'SYSTEM')
        )
);

CREATE UNIQUE INDEX idx_accounts_unique
ON accounts (owner_type, owner_id, account_code, currency, reference_id);