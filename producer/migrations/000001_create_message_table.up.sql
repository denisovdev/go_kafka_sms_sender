CREATE TABLE IF NOT EXISTS message (
  id          SERIAL PRIMARY KEY,
  payload     JSONB NOT NULL,
  status      VARCHAR(4) NOT NULL DEFAULT 'new' CHECK (status in ('new', 'done')),
  reserved_to TIMESTAMP DEFAULT NOW(),
  created_at  TIMESTAMP DEFAULT NOW()
);
