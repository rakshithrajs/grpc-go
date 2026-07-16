ALTER TABLE files DROP CONSTRAINT IF EXISTS files_user_name_unique;
ALTER TABLE files DROP CONSTRAINT IF EXISTS files_user_path_unique;

ALTER TABLE files ADD CONSTRAINT files_name_key UNIQUE (name);
ALTER TABLE files ADD CONSTRAINT files_path_key UNIQUE (path);
