ALTER TABLE "User"
DROP CONSTRAINT IF EXISTS valid_role;

ALTER TABLE "User"
DROP COLUMN IF EXISTS "twitter_handle";

ALTER TABLE "User"
DROP COLUMN IF EXISTS "instagram_handle";

ALTER TABLE "User"
DROP COLUMN IF EXISTS "creation_date";