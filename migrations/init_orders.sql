CREATE TABLE IF NOT EXISTS orders(
  id UUID PRIMARY KEY,
  user_id UUID,
  total_quantity INT,
  price_total INT,
  created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_products(
  id UUID PRIMARY KEY,
  order_id UUID,
  product_id UUID,
  price INT,
  quantity INT,
);


CREATE INDEX order_products_order_id
ON order_products(order_id); 

ALTER TABLE order_products
ADD FOREIGN KEY (order_id) REFERENCES orders(id); 
