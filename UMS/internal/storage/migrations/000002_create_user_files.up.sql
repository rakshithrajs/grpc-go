CREATE TABLE IF NOT EXISTS "userFiles" (
    "userID" UUID NOT NULL REFERENCES users("ID") ON DELETE CASCADE,
    "fileID" UUID NOT NULL,
    "createdAtUTC" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT userFiles_userID_fileID_unique UNIQUE ("userID", "fileID")
);

CREATE INDEX IF NOT EXISTS idx_userFiles_userID ON "userFiles" ("userID");
CREATE INDEX IF NOT EXISTS idx_userFiles_fileID ON "userFiles" ("fileID");
