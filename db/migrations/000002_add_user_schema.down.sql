-- ...existing code...
ALTER TABLE IF EXISTS "accounts"
DROP CONSTRAINT "accounts_owner_currency_idx";

ALTER TABLE IF EXISTS "accounts"
DROP CONSTRAINT "accounts_owner_fkey";

-- ...existing code...
DROP TABLE IF EXISTS "users";