CREATE TABLE IF NOT EXISTS users (
    "ID" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    phone VARCHAR(15) UNIQUE NOT NULL,
    "createdAtUTC" TIMESTAMPTZ DEFAULT current_timestamp,
    "updatedAtUTC" TIMESTAMPTZ DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
