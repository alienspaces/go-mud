-- type player provider 
CREATE TYPE "provider" AS ENUM (
  'anonymous',
  'apple',
  'facebook',
  'github',
  'google',
  'twitter'
);

-- type player role
CREATE TYPE "role" AS ENUM (
  'default',
  'administrator'
);

-- table player
CREATE TABLE "player" (
  "id"                  uuid CONSTRAINT player_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "name"                text NOT NULL,
  "email"               text NOT NULL,
  "provider"            provider NOT NULL,
  "provider_account_id" text NOT NULL,
  "created_at"          timestamp NOT NULL DEFAULT (current_timestamp),
  "updated_at"          timestamp DEFAULT null,
  "deleted_at"          timestamp DEFAULT null
);

-- table player role
CREATE TABLE "player_role" (
  "id"         uuid CONSTRAINT player_role_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "player_id"  uuid NOT NULL,
  "role"       role NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp DEFAULT null,
  "deleted_at" timestamp DEFAULT null
);

ALTER TABLE "player_role" ADD CONSTRAINT "player_role_player_id_fk" FOREIGN KEY ("player_id") REFERENCES "player" ("id");
