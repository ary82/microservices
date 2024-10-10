CREATE TABLE IF NOT EXISTS products(
  id UUID PRIMARY KEY,
  name VARCHAR(100),
  description TEXT,
  price INT,
  stock INT,
  created_at TIMESTAMP
);

ALTER TABLE products ADD CONSTRAINT chk_stock CHECK (stock >= 0);
