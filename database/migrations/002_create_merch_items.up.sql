-- Add up migration script here
CREATE TABLE
    IF NOT EXISTS "MerchItem" (
        "timestamp" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "id" uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4 (),
        "owner_id" uuid NOT NULL,
        "name" varchar(128) NOT NULL,
        "description" varchar(255),
        "price" decimal NOT NULL,
        "inventory" int,
        "type" varchar(128) NOT NULL, -- "postcard" etc
        "height" decimal,
        "width" decimal,
        "unit" varchar(10) NOT NULL, -- "cm" or "in"
        image_url1 varchar(255) NOT NULL,
        image_url2 varchar(255),
        image_url3 varchar(255),
        image_url4 varchar(255),
        image_url5 varchar(255),
        FOREIGN KEY ("owner_id") REFERENCES "User" ("id") ON DELETE CASCADE
    );