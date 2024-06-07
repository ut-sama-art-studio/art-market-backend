-- Add up migration script here

CREATE TABLE IF NOT EXISTS "MerchItem" (
    "id" uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "ownerId" uuid NOT NULL,
    "name" varchar NOT NULL,
    "description" varchar,
    "price" decimal NOT NULL,
    -- "images" text[],
    "inventory" int NOT NULL,
    "type" varchar NOT NULL,
    FOREIGN KEY ("ownerId") REFERENCES "USER" ("id") ON DELETE CASCADE
);