CREATE TABLE IF NOT EXISTS files (
    "ID" UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    "userID" UUID NOT NULL,
    "name" VARCHAR(150) NOT NULL,
    "path" VARCHAR(255) NOT NULL,
    "size" BIGINT NOT NULL,
    "mimeType" VARCHAR(30) NOT NULL,
    "createdAtUTC" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updatedAtUTC" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT files_user_name_unique UNIQUE ("userID", "name"),
    CONSTRAINT files_user_path_unique UNIQUE ("userID", "path")
);

CREATE INDEX IF NOT EXISTS idx_files_userID ON files ("userID");
CREATE INDEX IF NOT EXISTS idx_files_name ON files ("name");
