ALTER TABLE "User" ADD CONSTRAINT valid_role CHECK (
    "role" IN ('client', 'artist', 'director', 'admin')
);

ALTER TABLE "User"
ADD COLUMN IF NOT EXISTS "timestamp" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;