-- +migrate Up
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    transaction_type VARCHAR(50) NOT NULL, 
    -- DISBURSEMENT / REPAYMENT / FEE / ADJUSTMENT

    status VARCHAR(20) NOT NULL DEFAULT 'POSTED',
    -- PENDING / POSTED / REVERSED

    reference_type VARCHAR(50),
    reference_id UUID,

    idempotency_key VARCHAR(100) UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);