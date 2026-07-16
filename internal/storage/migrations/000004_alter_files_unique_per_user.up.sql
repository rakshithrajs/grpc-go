ALTER TABLE files DROP CONSTRAINT IF EXISTS files_name_key;
ALTER TABLE files DROP CONSTRAINT IF EXISTS files_path_key;

ALTER TABLE files ADD CONSTRAINT files_user_name_unique UNIQUE ("userID", name);
ALTER TABLE files ADD CONSTRAINT files_user_path_unique UNIQUE ("userID", path);
