
-- table dungeon
CREATE TABLE "dungeon" (
  "id"          uuid CONSTRAINT dungeon_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "name"        text NOT NULL,
  "description" text NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "dungeon_name_ck" CHECK (
    char_length("name") BETWEEN 1 AND 256
  )
);
COMMENT ON TABLE "dungeon" IS 'A dungeon is a set of locations that contain objects and monsters.';

-- TODO: effects

-- table object
CREATE TABLE "object" (
  "id"                   uuid CONSTRAINT object_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "name"                 text NOT NULL,
  "description"          text NOT NULL,
  "description_detailed" text NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "object_name_ck" CHECK (
    char_length("name") BETWEEN 1 AND 256
  ),
  CONSTRAINT "object_description_ck" CHECK (
    char_length("description") BETWEEN 1 AND 512
  ),
  CONSTRAINT "object_description_detailed_ck" CHECK (
    char_length("description_detailed") BETWEEN 1 AND 1024
  )
);
COMMENT ON TABLE "object" IS 'An object can be used, equipped, stashed or dropped.';

-- TODO: object effects

-- table monster
CREATE TABLE "monster" (
  "id"                uuid CONSTRAINT monster_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "name"              text NOT NULL,
  "description"       text NOT NULL,
  "strength"          integer NOT NULL DEFAULT 10,
  "dexterity"         integer NOT NULL DEFAULT 10,
  "intelligence"      integer NOT NULL DEFAULT 10,
  "health"            integer NOT NULL DEFAULT 0,
  "fatigue"           integer NOT NULL DEFAULT 0,
  "coins"             integer NOT NULL DEFAULT 0,
  "experience_points" integer NOT NULL DEFAULT 0,
  "attribute_points"  integer NOT NULL DEFAULT 0,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "monster_name_uq" UNIQUE("name"),
  CONSTRAINT "monster_name_ck" CHECK (
    char_length("name") BETWEEN 1 AND 256
  ),
  CONSTRAINT "monster_description_ck" CHECK (
    char_length("description") BETWEEN 1 AND 512
  )
);
COMMENT ON TABLE "monster" IS 'A monster can move, attack, converse with characters and interact objects.';

-- table monster_object
CREATE TABLE "monster_object" (
  "id"          uuid CONSTRAINT monster_object_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "monster_id"  uuid,
  "object_id"   uuid,
  "is_stashed"  boolean NOT NULL,
  "is_equipped" boolean NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "monster_object_monster_id_fk" FOREIGN KEY (monster_id) REFERENCES "monster"(id),
  CONSTRAINT "monster_object_object_id_fk" FOREIGN KEY (object_id) REFERENCES "object"(id),
  CONSTRAINT "monster_object_equipped_stashed_ck" CHECK (
    is_stashed != is_equipped
  )
);
COMMENT ON TABLE "monster_object" IS 'An object that is carried by a monster.';

-- table character
CREATE TABLE "character" (
  "id"                uuid CONSTRAINT character_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "name"              text NOT NULL,
  "strength"          integer NOT NULL DEFAULT 10,
  "dexterity"         integer NOT NULL DEFAULT 10,
  "intelligence"      integer NOT NULL DEFAULT 10,
  "health"            integer NOT NULL DEFAULT 0,
  "fatigue"           integer NOT NULL DEFAULT 0,
  "coins"             integer NOT NULL DEFAULT 0,
  "experience_points" integer NOT NULL DEFAULT 0,
  "attribute_points"  integer NOT NULL DEFAULT 0,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "character_name_uq" UNIQUE ("name"),
  CONSTRAINT "character_name_ck" CHECK (
    char_length("name") BETWEEN 1 AND 256
  )
);
COMMENT ON TABLE "monster" IS 'A character is controlled by a player and can move, attack, converse with monsters and interact with objects.';

-- table character_object
CREATE TABLE "character_object" (
  "id"           uuid CONSTRAINT character_object_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "character_id" uuid,
  "object_id"    uuid,
  "is_stashed"   boolean NOT NULL,
  "is_equipped"  boolean NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "character_object_character_id_fk" FOREIGN KEY (character_id) REFERENCES "character"(id),
  CONSTRAINT "character_object_object_id_fk" FOREIGN KEY (object_id) REFERENCES "object"(id),
  CONSTRAINT "character_object_equipped_stashed_ck" CHECK (
    is_stashed != is_equipped
  )
);
COMMENT ON TABLE "character_object" IS 'An object that is carried by a character.';

-- table location
CREATE TABLE "location" (
  "id"                    uuid CONSTRAINT location_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_id"            uuid NOT NULL,
  "name"                  text NOT NULL,
  "description"           text NOT NULL,
  "is_default"            boolean NOT NULL DEFAULT FALSE,
  "north_location_id"     uuid,
  "northeast_location_id" uuid,
  "east_location_id"      uuid,
  "southeast_location_id" uuid,
  "south_location_id"     uuid,
  "southwest_location_id" uuid,
  "west_location_id"      uuid,
  "northwest_location_id" uuid,
  "up_location_id"        uuid,
  "down_location_id"      uuid,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "location_name_ck" CHECK (
    char_length("name") > 0
  ),
  CONSTRAINT "location_description_ck" CHECK (
    char_length("description") > 0
  ),
  CONSTRAINT location_dungeon_id_fk FOREIGN KEY (dungeon_id) REFERENCES dungeon(id),
  CONSTRAINT location_north_location_id_fk FOREIGN KEY (north_location_id) REFERENCES location(id) INITIALLY DEFERRED,
  CONSTRAINT location_northeast_location_id_fk FOREIGN KEY (northeast_location_id) REFERENCES location(id) INITIALLY DEFERRED,
  CONSTRAINT location_east_location_id_fk FOREIGN KEY (east_location_id) REFERENCES location(id) INITIALLY DEFERRED,
  CONSTRAINT location_southeast_location_id_fk FOREIGN KEY (southeast_location_id) REFERENCES location(id) INITIALLY DEFERRED,
  CONSTRAINT location_south_location_id_fk FOREIGN KEY (south_location_id) REFERENCES location(id) INITIALLY DEFERRED,
  CONSTRAINT location_southwest_location_id_fk FOREIGN KEY (southwest_location_id) REFERENCES location(id) INITIALLY DEFERRED,
  CONSTRAINT location_west_location_id_fk FOREIGN KEY (west_location_id) REFERENCES location(id) INITIALLY DEFERRED,
  CONSTRAINT location_northwest_location_id_fk FOREIGN KEY (northwest_location_id) REFERENCES location(id) INITIALLY DEFERRED,
  CONSTRAINT location_up_location_id_fk FOREIGN KEY (up_location_id) REFERENCES location(id) INITIALLY DEFERRED,
  CONSTRAINT location_down_location_id_fk FOREIGN KEY (down_location_id) REFERENCES location(id) INITIALLY DEFERRED
);
COMMENT ON TABLE "location" IS 'A location is a room or place within a dungeon.';

-- table location_object
CREATE TABLE "location_object" (
  "id"                   uuid CONSTRAINT location_object_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "location_id"          uuid,
  "object_id"            uuid,
  "spawn_minutes"        integer NOT NULL DEFAULT 0,
  "spawn_percent_chance" integer NOT NULL DEFAULT 100,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "location_object_location_id_fk" FOREIGN KEY (location_id) REFERENCES "location"(id),
  CONSTRAINT "location_object_object_id_fk" FOREIGN KEY (object_id) REFERENCES "object"(id),
  CONSTRAINT "location_object_spawn_minutes_ck" CHECK (
    spawn_minutes BETWEEN 0 AND 60
  ),
  CONSTRAINT "location_object_spawn_percent_chance_ck" CHECK (
    spawn_percent_chance BETWEEN 1 AND 100
  )
);
COMMENT ON TABLE "location_object" IS 'An object that spawns at a location.';

-- table location_monster
CREATE TABLE "location_monster" (
  "id"                   uuid CONSTRAINT location_monster_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "location_id"          uuid NOT NULL,
  "monster_id"           uuid NOT NULL,
  "spawn_minutes"        integer NOT NULL DEFAULT 0,
  "spawn_percent_chance" integer NOT NULL DEFAULT 100,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "location_monster_location_id_fk" FOREIGN KEY (location_id) REFERENCES "location"(id),
  CONSTRAINT "location_monster_monster_id_fk" FOREIGN KEY (monster_id) REFERENCES "monster"(id),
  CONSTRAINT "location_monster_spawn_minutes_ck" CHECK (
    spawn_minutes BETWEEN 0 AND 60
  ),
  CONSTRAINT "location_monster_spawn_percent_chance_ck" CHECK (
    spawn_percent_chance BETWEEN 1 AND 100
  )
);
COMMENT ON TABLE "location_monster" IS 'A monster that spawns at a location.';

-- --
-- -- Instance objects
-- --

-- dungeon_instance
CREATE TABLE "dungeon_instance" (
  "id"         uuid CONSTRAINT dungeon_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_id" uuid NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT dungeon_instance_dungeon_id_fk FOREIGN KEY (dungeon_id) REFERENCES dungeon(id)
);

-- table location_instance
CREATE TABLE "location_instance" (
  "id"                             uuid CONSTRAINT location_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "location_id"                    uuid NOT NULL,
  "dungeon_instance_id"            uuid NOT NULL,
  "north_location_instance_id"     uuid,
  "northeast_location_instance_id" uuid,
  "east_location_instance_id"      uuid,
  "southeast_location_instance_id" uuid,
  "south_location_instance_id"     uuid,
  "southwest_location_instance_id" uuid,
  "west_location_instance_id"      uuid,
  "northwest_location_instance_id" uuid,
  "up_location_instance_id"        uuid,
  "down_location_instance_id"      uuid,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "location_instance_location_id_fk" FOREIGN KEY (location_id) REFERENCES location(id),
  CONSTRAINT "location_instance_dungeon_instance_id_fk" FOREIGN KEY (dungeon_instance_id) REFERENCES dungeon_instance(id),
  CONSTRAINT "location_instance_north_location_instance_id_fk" FOREIGN KEY (north_location_instance_id) REFERENCES location_instance(id) INITIALLY DEFERRED,
  CONSTRAINT "location_instance_northeast_location_instance_id_fk" FOREIGN KEY (northeast_location_instance_id) REFERENCES location_instance(id) INITIALLY DEFERRED,
  CONSTRAINT "location_instance_east_location_instance_id_fk" FOREIGN KEY (east_location_instance_id) REFERENCES location_instance(id) INITIALLY DEFERRED,
  CONSTRAINT "location_instance_southeast_location_instance_id_fk" FOREIGN KEY (southeast_location_instance_id) REFERENCES location_instance(id) INITIALLY DEFERRED,
  CONSTRAINT "location_instance_south_location_instance_id_fk" FOREIGN KEY (south_location_instance_id) REFERENCES location_instance(id) INITIALLY DEFERRED,
  CONSTRAINT "location_instance_southwest_location_instance_id_fk" FOREIGN KEY (southwest_location_instance_id) REFERENCES location_instance(id) INITIALLY DEFERRED,
  CONSTRAINT "location_instance_west_location_instance_id_fk" FOREIGN KEY (west_location_instance_id) REFERENCES location_instance(id) INITIALLY DEFERRED,
  CONSTRAINT "location_instance_northwest_location_instance_id_fk" FOREIGN KEY (northwest_location_instance_id) REFERENCES location_instance(id) INITIALLY DEFERRED,
  CONSTRAINT "location_instance_up_location_instance_id_fk" FOREIGN KEY (up_location_instance_id) REFERENCES location_instance(id) INITIALLY DEFERRED,
  CONSTRAINT "location_instance_down_location_instance_id_fk" FOREIGN KEY (down_location_instance_id) REFERENCES location_instance(id) INITIALLY DEFERRED
);

-- table character_instance
CREATE TABLE "character_instance" (
  "id"                   uuid CONSTRAINT character_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "character_id"         uuid NOT NULL,
  "dungeon_instance_id"  uuid NOT NULL,
  "location_instance_id" uuid NOT NULL,
  "strength"             integer NOT NULL DEFAULT 10,
  "dexterity"            integer NOT NULL DEFAULT 10,
  "intelligence"         integer NOT NULL DEFAULT 10,
  "health"               integer NOT NULL DEFAULT 0,
  "fatigue"              integer NOT NULL DEFAULT 0,
  "coins"                integer NOT NULL DEFAULT 0,
  "experience_points"    integer NOT NULL DEFAULT 0,
  "attribute_points"     integer NOT NULL DEFAULT 0,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "character_instance_character_id_fk" FOREIGN KEY (character_id) REFERENCES character(id),
  CONSTRAINT "character_instance_dungeon_instance_id_fk" FOREIGN KEY (dungeon_instance_id) REFERENCES dungeon_instance(id),
  CONSTRAINT "character_instance_location_instance_id_fk" FOREIGN KEY (location_instance_id) REFERENCES location_instance(id)
);

-- table monster_instance
CREATE TABLE "monster_instance" (
  "id"                   uuid CONSTRAINT monster_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "monster_id"           uuid NOT NULL,
  "dungeon_instance_id"  uuid NOT NULL,
  "location_instance_id" uuid NOT NULL,
  "strength"             integer NOT NULL DEFAULT 10,
  "dexterity"            integer NOT NULL DEFAULT 10,
  "intelligence"         integer NOT NULL DEFAULT 10,
  "health"               integer NOT NULL DEFAULT 0,
  "fatigue"              integer NOT NULL DEFAULT 0,
  "coins"                integer NOT NULL DEFAULT 0,
  "experience_points"    integer NOT NULL DEFAULT 0,
  "attribute_points"     integer NOT NULL DEFAULT 0,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "monster_instance_monster_id_fk" FOREIGN KEY (monster_id) REFERENCES monster(id),
  CONSTRAINT "monster_instance_dungeon_instance_id_fk" FOREIGN KEY (dungeon_instance_id) REFERENCES dungeon_instance(id),
  CONSTRAINT "monster_instance_location_instance_id_fk" FOREIGN KEY (location_instance_id) REFERENCES location_instance(id)
);

-- table object_instance
CREATE TABLE "object_instance" (
  "id"                    uuid CONSTRAINT object_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "object_id"             uuid NOT NULL,
  "dungeon_instance_id"   uuid NOT NULL,
  "location_instance_id"  uuid,
  "character_instance_id" uuid,
  "monster_instance_id"   uuid,
  "is_stashed"            boolean NOT NULL DEFAULT FALSE,
  "is_equipped"           boolean NOT NULL DEFAULT FALSE,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "object_instance_object_id_fk" FOREIGN KEY (object_id) REFERENCES object(id),
  CONSTRAINT "object_instance_dungeon_instance_id_fk" FOREIGN KEY (dungeon_instance_id) REFERENCES dungeon_instance(id),
  CONSTRAINT "object_instance_location_instance_id_fk" FOREIGN KEY (location_instance_id) REFERENCES location_instance(id),
  CONSTRAINT "object_instance_character_instance_id_fk" FOREIGN KEY (character_instance_id) REFERENCES character_instance(id),
  CONSTRAINT "object_instance_monster_instance_id_fk" FOREIGN KEY (monster_instance_id) REFERENCES monster_instance(id),
  CONSTRAINT "object_instance_location_character_monster_ck" CHECK (
    num_nonnulls(location_instance_id, character_instance_id, monster_instance_id) = 1
  )
);

-- --
-- -- Views
-- --

-- view dungeon_instance_view
CREATE OR REPLACE VIEW "dungeon_instance_view" AS
SELECT 
  di.id,
  di.dungeon_id,
  d.name,
  d.description,
  di.created_at,
  di.updated_at,
  di.deleted_at
FROM "dungeon_instance" di
JOIN "dungeon" d on d.id = di.dungeon_id
;

-- view location_instance_view
CREATE OR REPLACE VIEW "location_instance_view" AS
SELECT 
  li.id,
  l.dungeon_id,
  li.location_id,
  li.dungeon_instance_id,
  l.name,
  l.description,
  l.is_default,
  li.north_location_instance_id,
  li.northeast_location_instance_id,
  li.east_location_instance_id,
  li.southeast_location_instance_id,
  li.south_location_instance_id,
  li.southwest_location_instance_id,
  li.west_location_instance_id,
  li.northwest_location_instance_id,
  li.up_location_instance_id,
  li.down_location_instance_id,  
  li.created_at,
  li.updated_at,
  li.deleted_at
FROM "location_instance" li
JOIN "location" l on l.id = li.location_id
;

-- view character_instance_view
CREATE OR REPLACE VIEW "character_instance_view" AS
SELECT 
  ci.id,
  ci.character_id,
  ci.dungeon_instance_id,
  ci.location_instance_id,
  c.name,
  c.strength,
  c.dexterity,
  c.intelligence,
  ci.strength     as "current_strength",
  ci.dexterity    as "current_dexterity",
  ci.intelligence as "current_intelligence",
  c.health,
  c.fatigue,
  ci.health  as "current_health",
  ci.fatigue as "current_fatigue",
  ci.coins,
  ci.experience_points,
  ci.attribute_points,
  ci.created_at,
  ci.updated_at,
  ci.deleted_at
FROM "character_instance" ci
JOIN "character" c on c.id = ci.character_id
;

-- view monster_instance_view
CREATE OR REPLACE VIEW "monster_instance_view" AS
SELECT 
  mi.id,
  mi.monster_id,
  mi.dungeon_instance_id,
  mi.location_instance_id,
  m.name,
  m.strength,
  m.dexterity,
  m.intelligence,
  mi.strength     as "current_strength",
  mi.dexterity    as "current_dexterity",
  mi.intelligence as "current_intelligence",
  m.health,
  m.fatigue,
  mi.health  as "current_health",
  mi.fatigue as "current_fatigue",
  mi.coins,
  mi.experience_points,
  mi.attribute_points,
  mi.created_at,
  mi.updated_at,
  mi.deleted_at
FROM "monster_instance" mi
JOIN "monster" m on m.id = mi.monster_id
;


-- view object_instance_view
CREATE OR REPLACE VIEW "object_instance_view" AS
SELECT 
  oi.id,
  oi.object_id,
  oi.dungeon_instance_id,
  oi.location_instance_id,
  oi.character_instance_id,
  oi.monster_instance_id,
  o.name,
  o.description,
  o.description_detailed,
  oi.is_stashed,
  oi.is_equipped,
  oi.created_at,
  oi.updated_at,
  oi.deleted_at
FROM "object_instance" oi
JOIN "object" o on o.id = oi.object_id
;

-- --
-- -- Actions
-- --

-- table action
CREATE TABLE "action" (
  "id"                                           uuid CONSTRAINT action_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_instance_id"                          uuid NOT NULL,
  "location_instance_id"                         uuid NOT NULL,
  "character_instance_id"                        uuid,
  "monster_instance_id"                          uuid,
  "serial_id"                                    SERIAL,
  "resolved_command"                             text NOT NULL,
  "resolved_equipped_object_instance_id"         uuid,
  "resolved_stashed_object_instance_id"          uuid,
  "resolved_dropped_object_instance_id"          uuid,
  "resolved_target_object_instance_id"           uuid,
  "resolved_target_character_instance_id"        uuid,
  "resolved_target_monster_instance_id"          uuid,
  "resolved_target_location_direction"   text,
  "resolved_target_location_instance_id" uuid,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "action_dungeon_instance_id_fk" FOREIGN KEY (dungeon_instance_id) REFERENCES dungeon_instance(id),
  CONSTRAINT "action_location_instance_id_fk" FOREIGN KEY (location_instance_id) REFERENCES location_instance(id),
  CONSTRAINT "action_character_instance_id_fk" FOREIGN KEY (character_instance_id) REFERENCES character_instance(id),
  CONSTRAINT "action_monster_instance_id_fk" FOREIGN KEY (monster_instance_id) REFERENCES monster_instance(id),
  CONSTRAINT "action_resolved_equipped_object_instance_id_fk" FOREIGN KEY (resolved_equipped_object_instance_id) REFERENCES object_instance(id),
  CONSTRAINT "action_resolved_stashed_object_instance_id_fk" FOREIGN KEY (resolved_stashed_object_instance_id) REFERENCES   object_instance(id),
  CONSTRAINT "action_resolved_dropped_object_instance_id_fk" FOREIGN KEY (resolved_dropped_object_instance_id) REFERENCES   object_instance(id),
  CONSTRAINT "action_resolved_target_object_instance_id_fk" FOREIGN KEY (resolved_target_object_instance_id) REFERENCES object_instance(id),
  CONSTRAINT "action_resolved_target_character_instance_id_fk" FOREIGN KEY (resolved_target_character_instance_id) REFERENCES character_instance(id),
  CONSTRAINT "action_resolved_target_monster_instance_id_fk" FOREIGN KEY (resolved_target_monster_instance_id) REFERENCES monster_instance(id),
  CONSTRAINT "action_resolved_target_location_instance_id_fk" FOREIGN KEY (resolved_target_location_instance_id) REFERENCES location_instance(id),
  CONSTRAINT "action_character_or_monster_ck" CHECK 
  (
      ( CASE WHEN character_instance_id IS NULL THEN 0 ELSE 1 END
      + CASE WHEN monster_instance_id IS NULL THEN 0 ELSE 1 END
      ) = 1
  ),
  CONSTRAINT "action_target_instance_id_ck" CHECK (
    num_nonnulls(resolved_target_object_instance_id, resolved_target_character_instance_id, resolved_target_monster_instance_id, resolved_target_location_instance_id) = 1
  )
);

-- table action_character
CREATE TABLE "action_character" (
  "id"                           uuid CONSTRAINT action_character_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "record_type"                  text NOT NULL,
  "action_id"                    uuid NOT NULL,
  "location_instance_id"         uuid NOT NULL,
  "character_instance_id"        uuid NOT NULL,
  "name"                         text NOT NULL,
  "strength"                     integer NOT NULL,
  "dexterity"                    integer NOT NULL,
  "intelligence"                 integer NOT NULL,
  "current_strength"             integer NOT NULL,
  "current_dexterity"            integer NOT NULL,
  "current_intelligence"         integer NOT NULL,
  "health"                       integer NOT NULL,
  "fatigue"                      integer NOT NULL,
  "current_health"               integer NOT NULL,
  "current_fatigue"              integer NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "action_character_action_id_fk" FOREIGN KEY (action_id) REFERENCES action(id),
  CONSTRAINT "action_character_location_instance_id_fk" FOREIGN KEY (location_instance_id) REFERENCES location_instance(id),
  CONSTRAINT "action_character_character_instance_id_fk" FOREIGN KEY (character_instance_id) REFERENCES character_instance(id),
  CONSTRAINT "action_character_record_type_ck" CHECK (
    record_type = 'source' OR record_type = 'target' OR record_type = 'occupant'
  )
);

-- table action_character_object
CREATE TABLE "action_character_object" (
  "id"                    uuid CONSTRAINT action_character_object_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "action_id"             uuid NOT NULL,
  "character_instance_id" uuid NOT NULL,
  "object_instance_id"    uuid NOT NULL,
  "name"                  text NOT NULL,
  "is_stashed"            boolean NOT NULL,
  "is_equipped"           boolean NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "action_character_object_action_id_fk" FOREIGN KEY (action_id) REFERENCES action(id),
  CONSTRAINT "action_character_object_character_instance_id_fk" FOREIGN KEY (character_instance_id) REFERENCES character_instance(id),
  CONSTRAINT "action_character_object_object_instance_id_fk" FOREIGN KEY (object_instance_id) REFERENCES object_instance(id),
  CONSTRAINT "action_character_object_equipped_stashed_ck" CHECK (
    is_stashed != is_equipped
  )
);

-- table action_monster
CREATE TABLE "action_monster" (
  "id"                           uuid CONSTRAINT action_monster_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "record_type"                  text NOT NULL,
  "action_id"                    uuid NOT NULL,
  "location_instance_id"         uuid NOT NULL,
  "monster_instance_id"          uuid NOT NULL,
  "name"                         text NOT NULL,
  "strength"                     integer NOT NULL,
  "dexterity"                    integer NOT NULL,
  "intelligence"                 integer NOT NULL,
  "current_strength"             integer NOT NULL,
  "current_dexterity"            integer NOT NULL,
  "current_intelligence"         integer NOT NULL,
  "health"                       integer NOT NULL,
  "fatigue"                      integer NOT NULL,
  "current_health"               integer NOT NULL,
  "current_fatigue"              integer NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "action_monster_action_id_fk" FOREIGN KEY (action_id) REFERENCES action(id),
  CONSTRAINT "action_monster_location_instance_id_fk" FOREIGN KEY (location_instance_id) REFERENCES location_instance(id),
  CONSTRAINT "action_monster_monster_instance_id_fk" FOREIGN KEY (monster_instance_id) REFERENCES monster_instance(id),
  CONSTRAINT "action_monster_record_type_ck" CHECK (record_type = 'source' OR record_type = 'target' OR record_type = 'occupant')
);

-- table action_monster_object
CREATE TABLE "action_monster_object" (
  "id"                  uuid CONSTRAINT action_monster_object_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "action_id"           uuid NOT NULL,
  "monster_instance_id" uuid NOT NULL,
  "object_instance_id"  uuid NOT NULL,
  "name"                text NOT NULL,
  "is_stashed"          boolean NOT NULL,
  "is_equipped"         boolean NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "action_monster_object_action_id_fk" FOREIGN KEY (action_id) REFERENCES action(id),
  CONSTRAINT "action_monster_object_monster_instance_id_fk" FOREIGN KEY (monster_instance_id) REFERENCES monster_instance(id),
  CONSTRAINT "action_monster_object_object_instance_id_fk" FOREIGN KEY (object_instance_id) REFERENCES object_instance(id),
  CONSTRAINT "action_monster_object_equipped_stashed_ck" CHECK (
    is_stashed != is_equipped
  )
);

-- table action_object
CREATE TABLE "action_object" (
  "id"                           uuid CONSTRAINT action_object_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "record_type"                  text NOT NULL,
  "action_id"                    uuid NOT NULL,
  "location_instance_id"         uuid NOT NULL,
  "object_instance_id"           uuid NOT NULL,
  "name"                         text NOT NULL,
  "description"                  text NOT NULL,
  "is_stashed"                   boolean NOT NULL DEFAULT FALSE,
  "is_equipped"                  boolean NOT NULL DEFAULT FALSE,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "action_object_action_id_fk" FOREIGN KEY (action_id) REFERENCES action(id),
  CONSTRAINT "action_object_location_instance_id_fk" FOREIGN KEY (location_instance_id) REFERENCES location_instance(id),
  CONSTRAINT "action_object_object_instance_id_fk" FOREIGN KEY (object_instance_id) REFERENCES object_instance(id),
  CONSTRAINT "action_object_record_type_ck" CHECK (
    record_type = 'equipped' OR 
    record_type = 'stashed' OR 
    record_type = 'dropped' OR 
    record_type = 'target' OR 
    record_type = 'occupant'
  )
);


CREATE TABLE "turn" (
  "id"                  uuid CONSTRAINT turn_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_instance_id" uuid NOT NULL,
  "turn_count"          integer,
  "incremented_at"      timestamp WITH TIME ZONE, 
  "created_at"          timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at"          timestamp WITH TIME ZONE,
  "deleted_at"          timestamp WITH TIME ZONE,
  CONSTRAINT "turn_dungeon_instance_id_fk" FOREIGN KEY (dungeon_instance_id) REFERENCES dungeon_instance(id)
);
