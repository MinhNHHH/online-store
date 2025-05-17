-- Add your down migration here
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_status;
DROP INDEX IF EXISTS idx_products_price;
DROP INDEX IF EXISTS idx_products_created_at;

DROP INDEX IF EXISTS idx_categories_name;
DROP INDEX IF EXISTS idx_categories_created_at;

DROP INDEX IF EXISTS idx_product_categories_product_id;
DROP INDEX IF EXISTS idx_product_categories_category_id;

DROP TABLE IF EXISTS wishlist;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;
