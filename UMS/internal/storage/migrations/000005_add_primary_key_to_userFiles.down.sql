ALTER TABLE "userFiles"
DROP CONSTRAINT IF EXISTS userFiles_pkey,
ADD CONSTRAINT userFiles_userID_fileID_unique UNIQUE ("userID", "fileID");
