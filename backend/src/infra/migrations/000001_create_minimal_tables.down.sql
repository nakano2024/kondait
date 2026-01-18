DROP TRIGGER IF EXISTS trg_cooking_items_updated_at ON cooking_items;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP FUNCTION IF EXISTS set_updated_at();

DROP INDEX IF EXISTS idx_cooking_items_owner_code;

DROP TABLE IF EXISTS cooking_items;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "pgcrypto";
