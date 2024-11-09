DROP INDEX IF EXISTS accounts_owner_idx;
DROP INDEX IF EXISTS entries_account_id_idx;
DROP INDEX IF EXISTS transfers_from_account_id_idx;
DROP INDEX IF EXISTS transfers_to_account_id_idx;
DROP INDEX IF EXISTS transfers_from_account_id_to_account_id_idx;

DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS transfers;