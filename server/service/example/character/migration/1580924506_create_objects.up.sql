-- Character
CREATE TABLE "character" (
  "id"                uuid CONSTRAINT character_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "player_id"         uuid NOT NULL,
  "name"              text NOT NULL,
  "avatar"            text NOT NULL,
  "strength"          int NOT NULL DEFAULT 0,
  "dexterity"         int NOT NULL DEFAULT 0,
  "intelligence"      int NOT NULL DEFAULT 0,
  "attribute_points"  int NOT NULL DEFAULT 0,
  "experience_points" bigint NOT NULL DEFAULT 0,
  "coins"             bigint NOT NULL DEFAULT 0,
  "created_at"        timestamp NOT NULL DEFAULT (current_timestamp),
  "updated_at"        timestamp DEFAULT null,
  "deleted_at"        timestamp DEFAULT null
);
