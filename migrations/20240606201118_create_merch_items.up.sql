-- Add up migration script here
CREATE TABLE IF NOT EXISTS "MerchItem" (
    "id" uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4 (),
    "owner_id" uuid NOT NULL,
    "date_added" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "name" varchar NOT NULL,
    "description" varchar,
    "price" decimal NOT NULL,
    -- "images" text[],
    "inventory" int NOT NULL,
    "type" varchar NOT NULL,
    "height" decimal,
    "width" decimal,
    FOREIGN KEY ("owner_id") REFERENCES "User" ("id") ON DELETE CASCADE
);