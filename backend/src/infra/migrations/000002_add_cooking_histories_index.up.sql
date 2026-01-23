CREATE INDEX idx_cooking_histories_item_code_cooked_at
ON cooking_histories (cooking_item_code, cooked_at DESC);
