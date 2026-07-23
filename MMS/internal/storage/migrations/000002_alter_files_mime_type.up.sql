CREATE TYPE mime_type AS ENUM (
    'image/png',
    'image/jpeg',
    'image/gif',
    'image/webp',
    'image/svg+xml',
    'application/pdf',
    'text/plain',
    'text/markdown',
    'application/json'
);

ALTER TABLE files
ALTER COLUMN "mimeType" TYPE mime_type
USING "mimeType"::mime_type;
