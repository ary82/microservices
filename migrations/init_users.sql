CREATE TABLE IF NOT EXISTS users(
  id UUID PRIMARY KEY,
  username VARCHAR(100),
  email TEXT UNIQUE,
  salt TEXT,
  pass_hash TEXT,
  user_type INT,
  created_at TIMESTAMP
);
