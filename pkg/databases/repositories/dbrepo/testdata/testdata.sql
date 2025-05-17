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




-- Add your up migration here
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	description TEXT,
	price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
	stock_quantity INT NOT NULL CHECK (stock_quantity >= 0),
	status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('in_stock', 'out_of_stock', 'draft')),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL CHECK (name <> ''),
	description TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product_categories (
	product_id INT NOT NULL,
	category_id INT NOT NULL,
	PRIMARY KEY (product_id, category_id),
	FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
	FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE reviews (
	id SERIAL PRIMARY KEY,
	product_id INT NOT NULL,
	user_id INT NOT NULL,
	rating INT NOT NULL CHECK (
		rating >= 1
		AND rating <= 5
	),
	comment TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE wishlist (
	user_id INT NOT NULL,
	product_id INT NOT NULL,
	added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, product_id),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Products
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_products_created_at ON products(created_at);

-- Categories
CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_categories_created_at ON categories(created_at);

-- product_categories (join table)
CREATE INDEX idx_product_categories_product_id ON product_categories(product_id);
CREATE INDEX idx_product_categories_category_id ON product_categories(category_id);


-- USERS
INSERT INTO users (name, email, password) VALUES
('Alice Nguyen', 'alice@example.com', 'hashed_pwd_1'),
('Bob Tran', 'bob@example.com', 'hashed_pwd_2'),
('Carol Pham', 'carol@example.com', 'hashed_pwd_3');
