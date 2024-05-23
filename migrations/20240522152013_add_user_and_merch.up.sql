CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "USER" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "email" varchar NOT NULL UNIQUE,
    "profilePicture" varchar,
    "password" VARCHAR (127) NOT NULL,
    "bio" varchar
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "MerchItem" (
    "id" uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "ownerId" uuid NOT NULL,
    "name" varchar NOT NULL,
    "description" varchar,
    "price" decimal NOT NULL,
    "images" text[],
    "inventory" int NOT NULL,
    "type" varchar NOT NULL,
    FOREIGN KEY ("ownerId") REFERENCES "USER" ("id") ON DELETE CASCADE
);