ALTER TABLE "User" ADD CONSTRAINT valid_role CHECK (
    "role" IN ('client', 'artist', 'director', 'admin')
);

ALTER TABLE "User"
ADD COLUMN IF NOT EXISTS "twitter_handle" varchar(128);

ALTER TABLE "User"
ADD COLUMN IF NOT EXISTS "instagram_handle" varchar(128);

ALTER TABLE "User"
ADD COLUMN IF NOT EXISTS "creation_date" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;