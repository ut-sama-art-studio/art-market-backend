CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    IF NOT EXISTS "User" (
        "id" uuid NOT NULL DEFAULT uuid_generate_v4 (),
        "oauth_id" varchar(128) UNIQUE,
        "role" varchar(128) NOT NULL DEFAULT 'regular', -- 'regular', 'merchant', 'admin', 'super_admin'
        "username" varchar(128) NOT NULL,
        "name" varchar(128) NOT NULL,
        "email" VARCHAR(128) NOT NULL UNIQUE,
        "password" VARCHAR(128), -- hashed password
        "profile_picture" varchar(256),
        "bio" varchar,
        PRIMARY KEY ("id")
    );