INSERT INTO users (code, sub, created_at, updated_at) VALUES
('11111111-1111-1111-1111-111111111111', 'test-user-1', now(), now()),
('22222222-2222-2222-2222-222222222222', 'test-user-2', now(), now());

INSERT INTO cooking_items (code, owner_code, name, description, created_at, updated_at) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', '11111111-1111-1111-1111-111111111111', 'Rice', NULL, now(), now()),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2', '11111111-1111-1111-1111-111111111111', 'Pasta', NULL, now(), now()),
('cccccccc-cccc-cccc-cccc-ccccccccccc3', '11111111-1111-1111-1111-111111111111', 'Salad', NULL, now(), now()),
('dddddddd-dddd-dddd-dddd-ddddddddddd4', '11111111-1111-1111-1111-111111111111', 'Soup', NULL, now(), now()),
('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeee5', '11111111-1111-1111-1111-111111111111', 'Bread', NULL, now(), now()),
('99999999-9999-9999-9999-999999999999', '11111111-1111-1111-1111-111111111111', 'Curry', NULL, now(), now()),
('ffffffff-ffff-ffff-ffff-fffffffffff6', '22222222-2222-2222-2222-222222222222', 'Other', NULL, now(), now());

INSERT INTO cooking_histories (cooking_item_code, cooked_at, created_at) VALUES
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2', DATE '2024-01-01', now()),
('cccccccc-cccc-cccc-cccc-ccccccccccc3', DATE '2024-01-02', now()),
('dddddddd-dddd-dddd-dddd-ddddddddddd4', DATE '2024-01-03', now()),
('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeee5', DATE '2023-12-31', now()),
('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeee5', DATE '2024-01-01', now()),
('99999999-9999-9999-9999-999999999999', DATE '2023-12-30', now()),
('99999999-9999-9999-9999-999999999999', DATE '2023-12-31', now()),
('99999999-9999-9999-9999-999999999999', DATE '2024-01-01', now());
