
CREATE TABLE IF NOT EXISTS domains (
    id SERIAL PRIMARY KEY,
    domain TEXT UNIQUE NOT NULL,
    legal_name TEXT NOT NULL,
    country TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
