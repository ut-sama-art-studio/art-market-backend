ALTER TABLE "User"
    ADD CONSTRAINT valid_role CHECK ("role" IN ('client', 'artist', 'director', 'senpai', 'admin'));

ALTER TABLE "User"
    ADD COLUMN IF NOT EXISTS "timestamp" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE "User"
    ADD COLUMN IF NOT EXISTS "socials" TEXT[] DEFAULT '{}'; 
    -- in format: 'twitter:@kev_in', 'instagram:@kev_in', 'pixiv:@kev_in'