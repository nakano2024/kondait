
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    code UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sub TEXT NOT NULL,
    synced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE cooking_items (
    code UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_code UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    cook_count INTEGER NOT NULL DEFAULT 0,
    last_cooked_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_cooking_items_owner
        FOREIGN KEY (owner_code)
        REFERENCES users(code)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE INDEX idx_cooking_items_owner_code
ON cooking_items(owner_code);

-- updated_at を自動更新するトリガー
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_cooking_items_updated_at
BEFORE UPDATE ON cooking_items
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
