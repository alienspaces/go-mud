-- table dungeon
CREATE TABLE "dungeon" (
  "id"                    uuid CONSTRAINT dungeon_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "name"                  text NOT NULL,
  "description"           text NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE
);

-- dungeon_instance
CREATE TABLE "dungeon_instance" (
  "id"                    uuid CONSTRAINT dungeon_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_id"            uuid NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT dungeon_instance_dungeon_id_fk FOREIGN KEY (dungeon_id) REFERENCES dungeon(id)
);

-- dungeon_instance_view
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

-- location_instance_view
CREATE OR REPLACE VIEW "location_instance_view" AS
SELECT 
  dli.id,
  dl.dungeon_id,
  dli.location_id,
  dli.dungeon_instance_id,
  dl.name,
  dl.description,
  dl.is_default,
  dli.north_location_instance_id,
  dli.northeast_location_instance_id,
  dli.east_location_instance_id,
  dli.southeast_location_instance_id,
  dli.south_location_instance_id,
  dli.southwest_location_instance_id,
  dli.west_location_instance_id,
  dli.northwest_location_instance_id,
  dli.up_location_instance_id,
  dli.down_location_instance_id,  
  dli.created_at,
  dli.updated_at,
  dli.deleted_at
FROM "location_instance" dli
JOIN "location" dl on dl.id = dli.location_id
;

-- table character
CREATE TABLE "character" (
  "id"                   uuid CONSTRAINT character_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "name"                 text NOT NULL,
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
  CONSTRAINT "character_name_uq" UNIQUE ("name")
);

-- table character_instance
CREATE TABLE "character_instance" (
  "id"                           uuid CONSTRAINT character_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "character_id"                 uuid NOT NULL,
  "dungeon_instance_id"          uuid NOT NULL,
  "location_instance_id"         uuid NOT NULL,
  "strength"                     integer NOT NULL DEFAULT 10,
  "dexterity"                    integer NOT NULL DEFAULT 10,
  "intelligence"                 integer NOT NULL DEFAULT 10,
  "health"                       integer NOT NULL DEFAULT 0,
  "fatigue"                      integer NOT NULL DEFAULT 0,
  "coins"                        integer NOT NULL DEFAULT 0,
  "experience_points"            integer NOT NULL DEFAULT 0,
  "attribute_points"             integer NOT NULL DEFAULT 0,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "character_instance_character_id_fk" FOREIGN KEY (character_id) REFERENCES character(id),
  CONSTRAINT "character_instance_dungeon_instance_id_fk" FOREIGN KEY (dungeon_instance_id) REFERENCES dungeon_instance(id),
  CONSTRAINT "character_instance_location_instance_id_fk" FOREIGN KEY (location_instance_id) REFERENCES location_instance(id)
);

-- character_instance_view
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

-- table monster
CREATE TABLE "monster" (
  "id"                   uuid CONSTRAINT monster_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_id"           uuid NOT NULL,
  "location_id"          uuid NOT NULL,
  "name"                 text NOT NULL,
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
  CONSTRAINT "monster_dungeon_id_fk" FOREIGN KEY (dungeon_id) REFERENCES dungeon(id),
  CONSTRAINT "monster_location_id_fk" FOREIGN KEY (location_id) REFERENCES location(id),
  CONSTRAINT "monster_dungeon_id_name_uq" UNIQUE (dungeon_id, "name")
);

-- table monster_instance
CREATE TABLE "monster_instance" (
  "id"                           uuid CONSTRAINT monster_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "monster_id"                   uuid NOT NULL,
  "dungeon_instance_id"          uuid NOT NULL,
  "location_instance_id"         uuid NOT NULL,
  "strength"                     integer NOT NULL DEFAULT 10,
  "dexterity"                    integer NOT NULL DEFAULT 10,
  "intelligence"                 integer NOT NULL DEFAULT 10,
  "health"                       integer NOT NULL DEFAULT 0,
  "fatigue"                      integer NOT NULL DEFAULT 0,
  "coins"                        integer NOT NULL DEFAULT 0,
  "experience_points"            integer NOT NULL DEFAULT 0,
  "attribute_points"             integer NOT NULL DEFAULT 0,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "monster_instance_monster_id_fk" FOREIGN KEY (monster_id) REFERENCES monster(id),
  CONSTRAINT "monster_instance_dungeon_instance_id_fk" FOREIGN KEY (dungeon_instance_id) REFERENCES dungeon_instance(id),
  CONSTRAINT "monster_instance_location_instance_id_fk" FOREIGN KEY (location_instance_id) REFERENCES location_instance(id)
);

-- monster_instance_view
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

-- table object
CREATE TABLE "object" (
  "id"                   uuid CONSTRAINT object_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_id"           uuid NOT NULL,
  "location_id"          uuid,
  "character_id"         uuid,
  "monster_id"           uuid,
  "name"                 text NOT NULL,
  "description"          text NOT NULL,
  "description_detailed" text NOT NULL,
  "is_stashed"           boolean NOT NULL DEFAULT FALSE,
  "is_equipped"          boolean NOT NULL DEFAULT FALSE,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "object_dungeon_id_fk" FOREIGN KEY (dungeon_id) REFERENCES dungeon(id),
  CONSTRAINT "object_location_id_fk" FOREIGN KEY (location_id) REFERENCES location(id),
  CONSTRAINT "object_character_id_fk" FOREIGN KEY (character_id) REFERENCES character(id),
  CONSTRAINT "object_monster_id_fk" FOREIGN KEY (monster_id) REFERENCES monster(id),
  CONSTRAINT "object_dungeon_id_name_uq" UNIQUE (dungeon_id, "name"),
  CONSTRAINT "object_location_character_monster_ck" CHECK (
    num_nonnulls(location_id, character_id, monster_id) = 1
  ),
  CONSTRAINT "object_name_ck" CHECK (
    char_length("name") > 0
  ),
  CONSTRAINT "object_description_ck" CHECK (
    char_length("description") > 0
  ),
  CONSTRAINT "object_description_detailed_ck" CHECK (
    char_length("description_detailed") > 0
  )
  -- CHECK only stashed or equipped, not both..
);

-- table object_instance
CREATE TABLE "object_instance" (
  "id"                           uuid CONSTRAINT object_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "object_id"                    uuid NOT NULL,
  "dungeon_instance_id"          uuid NOT NULL,
  "location_instance_id"         uuid,
  "character_instance_id"        uuid,
  "monster_instance_id"          uuid,
  "is_stashed"                   boolean NOT NULL DEFAULT FALSE,
  "is_equipped"                  boolean NOT NULL DEFAULT FALSE,
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

-- object_instance_view
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
  oi.is_equipped,
  oi.is_stashed,
  oi.created_at,
  oi.updated_at,
  oi.deleted_at
FROM "object_instance" oi
JOIN "object" o on o.id = oi.object_id
;

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
