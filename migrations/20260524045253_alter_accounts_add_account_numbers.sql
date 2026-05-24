-- +migrate Up

ALTER TABLE accounts
ADD COLUMN account_number VARCHAR(30);

UPDATE accounts
SET account_number =
      '122'
   || '01'
   || LPAD(subquery.row_num::TEXT, 8, '0')
FROM (
    SELECT
        id,
        ROW_NUMBER() OVER (
            ORDER BY created_at ASC, id ASC
        ) AS row_num
    FROM accounts
    WHERE account_number IS NULL
) AS subquery
WHERE accounts.id = subquery.id;

ALTER TABLE accounts
ALTER COLUMN account_number SET NOT NULL;

CREATE UNIQUE INDEX idx_accounts_account_number
ON accounts(account_number);

CREATE SEQUENCE account_number_seq START 100000001;