INSERT INTO users (code, sub, created_at, updated_at) VALUES
('11111111-1111-1111-1111-111111111111', 'test-user-1', now(), now()),
('22222222-2222-2222-2222-222222222222', 'test-user-2', now(), now());

INSERT INTO cooking_items (code, owner_code, name, description, cook_count, last_cooked_date, created_at, updated_at) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', '11111111-1111-1111-1111-111111111111', 'Rice', NULL, 0, NULL, now(), now()),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2', '11111111-1111-1111-1111-111111111111', 'Pasta', NULL, 0, DATE '2024-01-01', now(), now()),
('cccccccc-cccc-cccc-cccc-ccccccccccc3', '11111111-1111-1111-1111-111111111111', 'Salad', NULL, 1, DATE '2024-01-02', now(), now()),
('dddddddd-dddd-dddd-dddd-ddddddddddd4', '11111111-1111-1111-1111-111111111111', 'Soup', NULL, 1, DATE '2024-01-03', now(), now()),
('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeee5', '11111111-1111-1111-1111-111111111111', 'Bread', NULL, 2, DATE '2024-01-01', now(), now()),
('99999999-9999-9999-9999-999999999999', '11111111-1111-1111-1111-111111111111', 'Curry', NULL, 3, DATE '2024-01-01', now(), now()),
('ffffffff-ffff-ffff-ffff-fffffffffff6', '22222222-2222-2222-2222-222222222222', 'Other', NULL, 0, NULL, now(), now());
