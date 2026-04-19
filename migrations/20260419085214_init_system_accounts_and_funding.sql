-- +migrate Up

-- =========================================================
-- 1. CREATE SYSTEM ACCOUNTS
-- =========================================================

-- SYSTEM CASH (operational fund)
INSERT INTO accounts (
    id,
    owner_type,
    owner_id,
    account_code,
    currency,
    status,
    created_at,
    updated_at
)
VALUES (
    uuid_generate_v4(),
    'SYSTEM',
    NULL,
    'SYSTEM_CASH',
    'IDR',
    'ACTIVE',
    NOW(),
    NOW()
)
ON CONFLICT DO NOTHING;

-- SYSTEM CAPITAL (source of initial funding)
INSERT INTO accounts (
    id,
    owner_type,
    owner_id,
    account_code,
    currency,
    status,
    created_at,
    updated_at
)
VALUES (
    uuid_generate_v4(),
    'SYSTEM',
    NULL,
    'SYSTEM_CAPITAL',
    'IDR',
    'ACTIVE',
    NOW(),
    NOW()
)
ON CONFLICT DO NOTHING;


-- =========================================================
-- 2. CREATE INITIAL FUNDING TRANSACTION
-- =========================================================

WITH tx AS (
    INSERT INTO transactions (
        id,
        transaction_type,
        status,
        reference_type,
        reference_id,
        idempotency_key,
        created_at
    )
    VALUES (
        uuid_generate_v4(),
        'ADJUSTMENT',                -- system funding treated as adjustment
        'POSTED',
        'SYSTEM',
        NULL,
        'INIT_SYSTEM_FUNDING',       -- prevent duplicate execution
        NOW()
    )
    ON CONFLICT (idempotency_key) DO NOTHING
    RETURNING id
),

acc AS (
    SELECT
        (SELECT id FROM accounts 
         WHERE account_code = 'SYSTEM_CASH' 
         LIMIT 1) AS cash_id,

        (SELECT id FROM accounts 
         WHERE account_code = 'SYSTEM_CAPITAL' 
         LIMIT 1) AS capital_id
)

-- =========================================================
-- 3. LEDGER ENTRIES (SIGNED AMOUNT MODEL)
-- =========================================================
-- IMPORTANT:
-- amount > 0  = CREDIT
-- amount < 0  = DEBIT

INSERT INTO ledger_entries (
    id,
    transaction_id,
    account_id,
    amount,
    currency,
    created_at
)

-- SYSTEM_CASH gets money → DEBIT → negative
SELECT
    uuid_generate_v4(),
    tx.id,
    acc.cash_id,
    -10000000000, -- 10 billion IDR
    'IDR',
    NOW()
FROM tx, acc

UNION ALL

-- SYSTEM_CAPITAL is source → CREDIT → positive
SELECT
    uuid_generate_v4(),
    tx.id,
    acc.capital_id,
    10000000000,
    'IDR',
    NOW()
FROM tx, acc;


-- =========================================================
-- 4. INITIALIZE ACCOUNT BALANCE SNAPSHOT
-- =========================================================

-- Insert or update balance snapshot for SYSTEM_CASH
INSERT INTO account_balances (
    account_id,
    balance,
    updated_at
)
SELECT
    id,
    10000000000,
    NOW()
FROM accounts
WHERE account_code = 'SYSTEM_CASH'
ON CONFLICT (account_id) DO UPDATE
SET balance = EXCLUDED.balance,
    updated_at = NOW();


-- Insert or update balance snapshot for SYSTEM_CAPITAL
INSERT INTO account_balances (
    account_id,
    balance,
    updated_at
)
SELECT
    id,
    -10000000000,
    NOW()
FROM accounts
WHERE account_code = 'SYSTEM_CAPITAL'
ON CONFLICT (account_id) DO UPDATE
SET balance = EXCLUDED.balance,
    updated_at = NOW();