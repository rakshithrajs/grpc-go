ALTER TABLE "userFiles"
DROP CONSTRAINT IF EXISTS userFiles_userID_fileID_unique,
ADD CONSTRAINT userFiles_pkey PRIMARY KEY ("userID", "fileID");
