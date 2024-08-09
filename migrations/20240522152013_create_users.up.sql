CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "User" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "oauth_id" varchar(128) UNIQUE, -- New field for Discord ID
    "name" varchar(127) NOT NULL,
    "email" VARCHAR (127) NOT NULL UNIQUE,
    "password" VARCHAR (127) NOT NULL, -- hashed password
    "profile_picture" varchar(256),
    "bio" varchar,
    PRIMARY KEY ("id")
);

