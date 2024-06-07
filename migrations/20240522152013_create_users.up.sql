CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "User" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(127) NOT NULL,
    "email" VARCHAR (127) NOT NULL UNIQUE,
    "password" VARCHAR (127) NOT NULL, -- hashed password
    "profilePicture" varchar(256),
    "bio" varchar,
    PRIMARY KEY ("id")
);

