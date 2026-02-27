-- +migrate Up

-- 1. Ensure balance type CASH exists (if not already inserted elsewhere)
INSERT INTO balance_types (id, code, description)
VALUES
  (1, 'CASH', 'User available balance')
ON CONFLICT (code) DO NOTHING;


-- 2. Create SYSTEM_TREASURY account
INSERT INTO accounts (id, owner_id, status)
VALUES (
    '00000000-0000-0000-0000-000000000010',
    '00000000-0000-0000-0000-000000000001',
    'ACTIVE'
)
ON CONFLICT (id) DO NOTHING;


-- 3. Initialize treasury balance = 10,000,000,000 IDR
INSERT INTO account_balances (account_id, balance_type_id, amount)
SELECT
    '00000000-0000-0000-0000-000000000010',
    bt.id,
    10000000000
FROM balance_types bt
WHERE bt.code = 'CASH'
ON CONFLICT (account_id, balance_type_id) DO NOTHING;



-- +migrate Down

DELETE FROM account_balances
WHERE account_id = '00000000-0000-0000-0000-000000000010';

DELETE FROM accounts
WHERE id = '00000000-0000-0000-0000-000000000010';