-- +goose Up
-- +goose StatementBegin
----------

WITH inserted_users AS (
    INSERT INTO users (user_id, email, password, name, role, created_at, updated_at)
    VALUES
        (gen_random_uuid(), 'user1@example.com', '$2a$10$cSYbO9l6f.EzZzzz/Z7hOOPbQ9XQiIMyQMvZSLk10LIoY8c.xUbhC', 'User 1', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
        (gen_random_uuid(), 'user2@example.com', '$2a$10$JAjhZkZ0ajVT77JpN8a/E.kq7p9U8Xfwr29g2vbWLnfsS/KibHNoW', 'User 2', 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
        (gen_random_uuid(), 'user3@example.com', '$2a$10$pb0EMw4Eyk4qDVRzYWIVZO6GeWP4jovFWU8PsmoqPxqjk25nLO92u', 'User 3', 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    RETURNING user_id, email
),
inserted_products AS (
    INSERT INTO products (product_id, name, description, price, stock, created_at, updated_at)
    VALUES
        (gen_random_uuid(), 'Laptop', 'A high-performance laptop with 16GB RAM and 512GB SSD', 1200.99, 50, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
        (gen_random_uuid(), 'Smartphone', 'A flagship smartphone with excellent camera', 799.99, 100, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
        (gen_random_uuid(), 'Headphones', 'Noise-cancelling wireless headphones', 199.99, 150, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
        (gen_random_uuid(), 'Keyboard', 'Mechanical keyboard with RGB lighting', 99.99, 200, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
        (gen_random_uuid(), 'Mouse', 'Wireless mouse with ergonomic design', 49.99, 300, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    RETURNING product_id, name
),
inserted_orders AS (
    INSERT INTO orders (order_id, user_id, order_date, total_price, status, address, created_at, updated_at)
    SELECT
        gen_random_uuid(), user_id, CURRENT_TIMESTAMP, 1499.98, 'Pending', '123 Main St, Cityville, NY', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
    FROM inserted_users
    WHERE email = 'user1@example.com'
    UNION ALL
    SELECT
        gen_random_uuid(), user_id, CURRENT_TIMESTAMP, 799.99, 'Shipped', '456 Oak Ave, Townsville, TX', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
    FROM inserted_users
    WHERE email = 'user2@example.com'
    UNION ALL
    SELECT
        gen_random_uuid(), user_id, CURRENT_TIMESTAMP, 2499.96, 'Delivered', '789 Pine Blvd, Villagetown, CA', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
    FROM inserted_users
    WHERE email = 'user3@example.com'
    RETURNING order_id, user_id
),
inserted_order_items AS (
    INSERT INTO order_items (order_item_id, order_id, product_id, quantity, unit_price, total_price)
    SELECT
        gen_random_uuid(), o.order_id, p.product_id, 1, 1200.99, 1200.99
    FROM inserted_orders o
    JOIN inserted_products p ON p.name = 'Laptop'
    WHERE o.order_id IN (SELECT order_id FROM inserted_orders WHERE user_id IN (SELECT user_id FROM inserted_users WHERE email = 'user1@example.com'))
    UNION ALL
    SELECT
        gen_random_uuid(), o.order_id, p.product_id, 1, 199.99, 199.99
    FROM inserted_orders o
    JOIN inserted_products p ON p.name = 'Headphones'
    WHERE o.order_id IN (SELECT order_id FROM inserted_orders WHERE user_id IN (SELECT user_id FROM inserted_users WHERE email = 'user1@example.com'))
    UNION ALL
    SELECT
        gen_random_uuid(), o.order_id, p.product_id, 1, 799.99, 799.99
    FROM inserted_orders o
    JOIN inserted_products p ON p.name = 'Smartphone'
    WHERE o.order_id IN (SELECT order_id FROM inserted_orders WHERE user_id IN (SELECT user_id FROM inserted_users WHERE email = 'user2@example.com'))
    UNION ALL
    SELECT
        gen_random_uuid(), o.order_id, p.product_id, 2, 1200.99, 2401.98
    FROM inserted_orders o
    JOIN inserted_products p ON p.name = 'Laptop'
    WHERE o.order_id IN (SELECT order_id FROM inserted_orders WHERE user_id IN (SELECT user_id FROM inserted_users WHERE email = 'user3@example.com'))
    UNION ALL
    SELECT
        gen_random_uuid(), o.order_id, p.product_id, 1, 99.99, 99.99
    FROM inserted_orders o
    JOIN inserted_products p ON p.name = 'Keyboard'
    WHERE o.order_id IN (SELECT order_id FROM inserted_orders WHERE user_id IN (SELECT user_id FROM inserted_users WHERE email = 'user3@example.com'))
    RETURNING order_item_id, order_id, product_id
)
-- Execute the insertion of order items
SELECT * FROM inserted_order_items;

----------
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
----------

-- Delete all order items first
DELETE FROM order_items;

-- Delete all orders
DELETE FROM orders;

-- Delete all products
DELETE FROM products;

-- Delete all users
DELETE FROM users;

----------
-- +goose StatementEnd
