-- table dungeon
CREATE TABLE "dungeon" (
  "id"                    uuid CONSTRAINT dungeon_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "name"                  text NOT NULL,
  "description"           text NOT NULL,
  "created_at"            timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at"            timestamp WITH TIME ZONE,
  "deleted_at"            timestamp WITH TIME ZONE
);

-- table dungeon_location
CREATE TABLE "dungeon_location" (
  "id"                            uuid CONSTRAINT dungeon_location_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_id"                    uuid NOT NULL,
  "name"                          text NOT NULL,
  "description"                   text NOT NULL,
  "default"                       boolean NOT NULL DEFAULT FALSE,
  "north_dungeon_location_id"     uuid,
  "northeast_dungeon_location_id" uuid,
  "east_dungeon_location_id"      uuid,
  "southeast_dungeon_location_id" uuid,
  "south_dungeon_location_id"     uuid,
  "southwest_dungeon_location_id" uuid,
  "west_dungeon_location_id"      uuid,
  "northwest_dungeon_location_id" uuid,
  "up_dungeon_location_id"        uuid,
  "down_dungeon_location_id"      uuid,
  "created_at"                    timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at"                    timestamp WITH TIME ZONE,
  "deleted_at"                    timestamp WITH TIME ZONE,
  CONSTRAINT dungeon_location_dungeon_id_fk FOREIGN KEY (dungeon_id) REFERENCES dungeon(id),
  CONSTRAINT dungeon_location_north_location_id_fk FOREIGN KEY (north_dungeon_location_id) REFERENCES dungeon_location(id) INITIALLY DEFERRED,
  CONSTRAINT dungeon_location_northeast_location_id_fk FOREIGN KEY (northeast_dungeon_location_id) REFERENCES dungeon_location(id) INITIALLY DEFERRED,
  CONSTRAINT dungeon_location_east_location_id_fk FOREIGN KEY (east_dungeon_location_id) REFERENCES dungeon_location(id) INITIALLY DEFERRED,
  CONSTRAINT dungeon_location_southeast_location_id_fk FOREIGN KEY (southeast_dungeon_location_id) REFERENCES dungeon_location(id) INITIALLY DEFERRED,
  CONSTRAINT dungeon_location_south_location_id_fk FOREIGN KEY (south_dungeon_location_id) REFERENCES dungeon_location(id) INITIALLY DEFERRED,
  CONSTRAINT dungeon_location_southwest_location_id_fk FOREIGN KEY (southwest_dungeon_location_id) REFERENCES dungeon_location(id) INITIALLY DEFERRED,
  CONSTRAINT dungeon_location_west_location_id_fk FOREIGN KEY (west_dungeon_location_id) REFERENCES dungeon_location(id) INITIALLY DEFERRED,
  CONSTRAINT dungeon_location_northwest_location_id_fk FOREIGN KEY (northwest_dungeon_location_id) REFERENCES dungeon_location(id) INITIALLY DEFERRED,
  CONSTRAINT dungeon_location_up_location_id_fk FOREIGN KEY (up_dungeon_location_id) REFERENCES dungeon_location(id) INITIALLY DEFERRED,
  CONSTRAINT dungeon_location_down_location_id_fk FOREIGN KEY (down_dungeon_location_id) REFERENCES dungeon_location(id) INITIALLY DEFERRED
);

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
  "created_at"           timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at"           timestamp WITH TIME ZONE,
  "deleted_at"           timestamp WITH TIME ZONE,
  CONSTRAINT "character_name_uq" UNIQUE ("name")
);

-- table monster
CREATE TABLE "monster" (
  "id"                   uuid CONSTRAINT monster_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_id"           uuid NOT NULL,
  "dungeon_location_id"  uuid NOT NULL,
  "name"                 text NOT NULL,
  "strength"             integer NOT NULL DEFAULT 10,
  "dexterity"            integer NOT NULL DEFAULT 10,
  "intelligence"         integer NOT NULL DEFAULT 10,
  "current_strength"     integer NOT NULL DEFAULT 10,
  "current_dexterity"    integer NOT NULL DEFAULT 10,
  "current_intelligence" integer NOT NULL DEFAULT 10,
  "health"               integer NOT NULL DEFAULT 0,
  "fatigue"              integer NOT NULL DEFAULT 0,
  "current_health"       integer NOT NULL DEFAULT 0,
  "current_fatigue"      integer NOT NULL DEFAULT 0,
  "coins"                integer NOT NULL DEFAULT 0,
  "experience_points"    integer NOT NULL DEFAULT 0,
  "attribute_points"     integer NOT NULL DEFAULT 0,
  "created_at"           timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at"           timestamp WITH TIME ZONE,
  "deleted_at"           timestamp WITH TIME ZONE,
  CONSTRAINT "monster_dungeon_id_fk" FOREIGN KEY (dungeon_id) REFERENCES dungeon(id),
  CONSTRAINT "monster_dungeon_location_id_fk" FOREIGN KEY (dungeon_location_id) REFERENCES dungeon_location(id),
  CONSTRAINT "monster_dungeon_id_name_uq" UNIQUE (dungeon_id, "name")
);

-- table object
CREATE TABLE "object" (
  "id"                   uuid CONSTRAINT object_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_id"           uuid NOT NULL,
  "dungeon_location_id"  uuid,
  "character_id"         uuid,
  "monster_id"           uuid,
  "name"                 text NOT NULL,
  "description"          text NOT NULL,
  "description_detailed" text NOT NULL,
  "is_stashed"           boolean NOT NULL DEFAULT FALSE,
  "is_equipped"          boolean NOT NULL DEFAULT FALSE,
  "created_at"           timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at"           timestamp WITH TIME ZONE,
  "deleted_at"           timestamp WITH TIME ZONE,
  CONSTRAINT "object_dungeon_id_fk" FOREIGN KEY (dungeon_id) REFERENCES dungeon(id),
  CONSTRAINT "object_dungeon_location_id_fk" FOREIGN KEY (dungeon_location_id) REFERENCES dungeon_location(id),
  CONSTRAINT "object_character_id_fk" FOREIGN KEY (character_id) REFERENCES character(id),
  CONSTRAINT "object_monster_id_fk" FOREIGN KEY (monster_id) REFERENCES monster(id),
  CONSTRAINT "object_dungeon_id_name_uq" UNIQUE (dungeon_id, "name"),
  CONSTRAINT "object_location_character_monster_ck" CHECK (
    num_nonnulls(dungeon_location_id, character_id, monster_id) = 1
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

-- dungeon_instance
CREATE TABLE "dungeon_instance" (
  "id"                    uuid CONSTRAINT dungeon_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_id"            uuid NOT NULL,
  "created_at"            timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at"            timestamp WITH TIME ZONE,
  "deleted_at"            timestamp WITH TIME ZONE,
  CONSTRAINT dungeon_instance_dungeon_id_fk FOREIGN KEY (dungeon_id) REFERENCES dungeon(id)
);

-- table dungeon_location_instance
CREATE TABLE "dungeon_location_instance" (
  "id"                    uuid CONSTRAINT dungeon_location_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_location_id"   uuid NOT NULL,
  "created_at"            timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at"            timestamp WITH TIME ZONE,
  "deleted_at"            timestamp WITH TIME ZONE,
  CONSTRAINT dungeon_location_instance_dungeon_location_id_fk FOREIGN KEY (dungeon_location_id) REFERENCES dungeon_location(id)
);

-- table character_instance
CREATE TABLE "character_instance" (
  "id"                           uuid CONSTRAINT character_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "character_id"                 uuid NOT NULL,
  "dungeon_instance_id"          uuid NOT NULL,
  "dungeon_location_instance_id" uuid NOT NULL,
  "current_strength"             integer NOT NULL DEFAULT 10,
  "current_dexterity"            integer NOT NULL DEFAULT 10,
  "current_intelligence"         integer NOT NULL DEFAULT 10,
  "current_health"               integer NOT NULL DEFAULT 0,
  "current_fatigue"              integer NOT NULL DEFAULT 0,
  "coins"                        integer NOT NULL DEFAULT 0,
  "experience_points"            integer NOT NULL DEFAULT 0,
  "attribute_points"             integer NOT NULL DEFAULT 0,
  "created_at"                   timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at"                   timestamp WITH TIME ZONE,
  "deleted_at"                   timestamp WITH TIME ZONE,
  CONSTRAINT "character_instance_character_id_fk" FOREIGN KEY (character_id) REFERENCES character(id),
  CONSTRAINT "character_instance_dungeon_instance_id_fk" FOREIGN KEY (dungeon_instance_id) REFERENCES dungeon_instance(id),
  CONSTRAINT "character_instance_dungeon_location_instance_id_fk" FOREIGN KEY (dungeon_location_instance_id) REFERENCES dungeon_location_instance(id)
);

-- table monster_instance
CREATE TABLE "monster_instance" (
  "id"                           uuid CONSTRAINT monster_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "monster_id"                   uuid NOT NULL,
  "dungeon_instance_id"          uuid NOT NULL,
  "dungeon_location_instance_id" uuid NOT NULL,
  "current_strength"             integer NOT NULL DEFAULT 10,
  "current_dexterity"            integer NOT NULL DEFAULT 10,
  "current_intelligence"         integer NOT NULL DEFAULT 10,
  "current_health"               integer NOT NULL DEFAULT 0,
  "current_fatigue"              integer NOT NULL DEFAULT 0,
  "coins"                        integer NOT NULL DEFAULT 0,
  "experience_points"            integer NOT NULL DEFAULT 0,
  "attribute_points"             integer NOT NULL DEFAULT 0,
  "created_at"                   timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at"                   timestamp WITH TIME ZONE,
  "deleted_at"                   timestamp WITH TIME ZONE,
  CONSTRAINT "monster_instance_monster_id_fk" FOREIGN KEY (monster_id) REFERENCES monster(id),
  CONSTRAINT "monster_instance_dungeon_instance_id_fk" FOREIGN KEY (dungeon_instance_id) REFERENCES dungeon_instance(id),
  CONSTRAINT "monster_instance_dungeon_location_instance_id_fk" FOREIGN KEY (dungeon_location_instance_id) REFERENCES dungeon_location_instance(id)
);

-- table object_instance
CREATE TABLE "object_instance" (
  "id"                           uuid CONSTRAINT object_instance_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "object_id"                    uuid NOT NULL,
  "dungeon_instance_id"          uuid NOT NULL,
  "dungeon_location_instance_id" uuid,
  "character_instance_id"        uuid,
  "monster_instance_id"          uuid,
  "is_stashed"                   boolean NOT NULL DEFAULT FALSE,
  "is_equipped"                  boolean NOT NULL DEFAULT FALSE,
  "created_at"                   timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at"                   timestamp WITH TIME ZONE,
  "deleted_at"                   timestamp WITH TIME ZONE,
  CONSTRAINT "object_instance_object_id_fk" FOREIGN KEY (object_id) REFERENCES object(id),
  CONSTRAINT "object_instance_dungeon_instance_id_fk" FOREIGN KEY (dungeon_instance_id) REFERENCES dungeon_instance(id),
  CONSTRAINT "object_instance_dungeon_location_instance_id_fk" FOREIGN KEY (dungeon_location_instance_id) REFERENCES dungeon_location_instance(id),
  CONSTRAINT "object_instance_character_instance_id_fk" FOREIGN KEY (character_instance_id) REFERENCES character_instance(id),
  CONSTRAINT "object_instance_monster_instance_id_fk" FOREIGN KEY (monster_instance_id) REFERENCES monster_instance(id),
  CONSTRAINT "object_instance_location_character_monster_ck" CHECK (
    num_nonnulls(dungeon_location_instance_id, character_instance_id, monster_instance_id) = 1
  )
);

-- table dungeon_action
CREATE TABLE "dungeon_action" (
  "id"                                           uuid CONSTRAINT dungeon_action_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_instance_id"                          uuid NOT NULL,
  "dungeon_location_instance_id"                 uuid NOT NULL,
  "character_id"                                 uuid,
  "monster_id"                                   uuid,
  "serial_id"                                    SERIAL,
  "resolved_command"                             text NOT NULL,
  "resolved_equipped_object_instance_id"         uuid,
  "resolved_stashed_object_instance_id"          uuid,
  "resolved_dropped_object_instance_id"          uuid,
  "resolved_target_object_instance_id"           uuid,
  "resolved_target_character_instance_id"        uuid,
  "resolved_target_monster_instance_id"          uuid,
  "resolved_target_dungeon_location_direction"   text,
  "resolved_target_dungeon_location_instance_id" uuid,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "dungeon_action_dungeon_instance_id_fk" FOREIGN KEY (dungeon_instance_id) REFERENCES dungeon_instance(id),
  CONSTRAINT "dungeon_action_dungeon_location_instance_id_fk" FOREIGN KEY (dungeon_location_instance_id) REFERENCES dungeon_location_instance(id),
  CONSTRAINT "dungeon_action_character_id_fk" FOREIGN KEY (character_id) REFERENCES character(id),
  CONSTRAINT "dungeon_action_resolved_equipped_object_instance_id_fk" FOREIGN KEY (resolved_equipped_object_instance_id) REFERENCES object_instance(id),
  CONSTRAINT "dungeon_action_resolved_stashed_object_instance_id_fk" FOREIGN KEY (resolved_stashed_object_instance_id) REFERENCES   object_instance(id),
  CONSTRAINT "dungeon_action_resolved_dropped_object_instance_id_fk" FOREIGN KEY (resolved_dropped_object_instance_id) REFERENCES   object_instance(id),
  CONSTRAINT "dungeon_action_resolved_target_object_instance_id_fk" FOREIGN KEY (resolved_target_object_instance_id) REFERENCES object_instance(id),
  CONSTRAINT "dungeon_action_resolved_target_character_instance_id_fk" FOREIGN KEY (resolved_target_character_instance_id) REFERENCES character_instance(id),
  CONSTRAINT "dungeon_action_resolved_target_monster_instance_id_fk" FOREIGN KEY (resolved_target_monster_instance_id) REFERENCES monster_instance(id),
  CONSTRAINT "dungeon_action_resolved_target_dungeon_location_instance_id_fk" FOREIGN KEY (resolved_target_dungeon_location_instance_id) REFERENCES dungeon_location_instance(id),
  CONSTRAINT "dungeon_action_character_or_monster_ck" CHECK 
  (
      ( CASE WHEN character_id IS NULL THEN 0 ELSE 1 END
      + CASE WHEN monster_id IS NULL THEN 0 ELSE 1 END
      ) = 1
  ),
  CONSTRAINT "dungeon_action_target_instance_id_ck" CHECK (
    num_nonnulls(resolved_target_object_instance_id, resolved_target_character_instance_id, resolved_target_monster_instance_id, resolved_target_dungeon_location_instance_id) = 1
  )
);

-- table dungeon_action_character
CREATE TABLE "dungeon_action_character" (
  "id"                           uuid CONSTRAINT dungeon_action_character_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "record_type"                  text NOT NULL,
  "dungeon_action_id"            uuid NOT NULL,
  "dungeon_location_instance_id" uuid NOT NULL,
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
  CONSTRAINT "dungeon_action_character_dungeon_action_id_fk" FOREIGN KEY (dungeon_action_id) REFERENCES dungeon_action(id),
  CONSTRAINT "dungeon_action_character_dungeon_location_instance_id_fk" FOREIGN KEY (dungeon_location_instance_id) REFERENCES dungeon_location_instance(id),
  CONSTRAINT "dungeon_action_character_character_instance_id_fk" FOREIGN KEY (character_instance_id) REFERENCES character_instance(id),
  CONSTRAINT "dungeon_action_character_record_type_ck" CHECK (
    record_type = 'source' OR record_type = 'target' OR record_type = 'occupant'
  )
);

-- table dungeon_action_character_object
CREATE TABLE "dungeon_action_character_object" (
  "id"                    uuid CONSTRAINT dungeon_action_character_object_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_action_id"     uuid NOT NULL,
  "character_instance_id" uuid NOT NULL,
  "object_instance_id"    uuid NOT NULL,
  "name"                  text NOT NULL,
  "is_stashed"            boolean NOT NULL,
  "is_equipped"           boolean NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "dungeon_action_character_object_dungeon_action_id_fk" FOREIGN KEY (dungeon_action_id) REFERENCES dungeon_action(id),
  CONSTRAINT "dungeon_action_character_object_character_instance_id_fk" FOREIGN KEY (character_instance_id) REFERENCES character_instance(id),
  CONSTRAINT "dungeon_action_character_object_object_instance_id_fk" FOREIGN KEY (object_instance_id) REFERENCES object_instance(id),
  CONSTRAINT "dungeon_action_character_object_equipped_stashed_ck" CHECK (
    is_stashed != is_equipped
  )
);

-- table dungeon_action_monster
CREATE TABLE "dungeon_action_monster" (
  "id"                           uuid CONSTRAINT dungeon_action_monster_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "record_type"                  text NOT NULL,
  "dungeon_action_id"            uuid NOT NULL,
  "dungeon_location_instance_id" uuid NOT NULL,
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
  CONSTRAINT "dungeon_action_monster_dungeon_action_id_fk" FOREIGN KEY (dungeon_action_id) REFERENCES dungeon_action(id),
  CONSTRAINT "dungeon_action_monster_dungeon_location_instance_id_fk" FOREIGN KEY (dungeon_location_instance_id) REFERENCES dungeon_location_instance(id),
  CONSTRAINT "dungeon_action_monster_monster_instance_id_fk" FOREIGN KEY (monster_instance_id) REFERENCES monster_instance(id),
  CONSTRAINT "dungeon_action_monster_record_type_ck" CHECK (record_type = 'source' OR record_type = 'target' OR record_type = 'occupant')
);

-- table dungeon_action_monster_object
CREATE TABLE "dungeon_action_monster_object" (
  "id"                  uuid CONSTRAINT dungeon_action_monster_object_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "dungeon_action_id"   uuid NOT NULL,
  "monster_instance_id" uuid NOT NULL,
  "object_instance_id"  uuid NOT NULL,
  "name"                text NOT NULL,
  "is_stashed"          boolean NOT NULL,
  "is_equipped"         boolean NOT NULL,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "dungeon_action_monster_object_dungeon_action_id_fk" FOREIGN KEY (dungeon_action_id) REFERENCES dungeon_action(id),
  CONSTRAINT "dungeon_action_monster_object_monster_instance_id_fk" FOREIGN KEY (monster_instance_id) REFERENCES monster_instance(id),
  CONSTRAINT "dungeon_action_monster_object_object_instance_id_fk" FOREIGN KEY (object_instance_id) REFERENCES object_instance(id),
  CONSTRAINT "dungeon_action_monster_object_equipped_stashed_ck" CHECK (
    is_stashed != is_equipped
  )
);

-- table dungeon_action_object
CREATE TABLE "dungeon_action_object" (
  "id"                           uuid CONSTRAINT dungeon_action_object_pk PRIMARY KEY DEFAULT gen_random_uuid(),
  "record_type"                  text NOT NULL,
  "dungeon_action_id"            uuid NOT NULL,
  "dungeon_location_instance_id" uuid NOT NULL,
  "object_instance_id"           uuid NOT NULL,
  "name"                         text NOT NULL,
  "description"                  text NOT NULL,
  "is_stashed"                   boolean NOT NULL DEFAULT FALSE,
  "is_equipped"                  boolean NOT NULL DEFAULT FALSE,
  "created_at" timestamp WITH TIME ZONE NOT NULL DEFAULT (current_timestamp),
  "updated_at" timestamp WITH TIME ZONE,
  "deleted_at" timestamp WITH TIME ZONE,
  CONSTRAINT "dungeon_action_object_dungeon_action_id_fk" FOREIGN KEY (dungeon_action_id) REFERENCES dungeon_action(id),
  CONSTRAINT "dungeon_action_object_dungeon_location_instance_id_fk" FOREIGN KEY (dungeon_location_instance_id) REFERENCES dungeon_location_instance(id),
  CONSTRAINT "dungeon_action_object_object_instance_id_fk" FOREIGN KEY (object_instance_id) REFERENCES object_instance(id),
  CONSTRAINT "dungeon_action_object_record_type_ck" CHECK (
    record_type = 'equipped' OR 
    record_type = 'stashed' OR 
    record_type = 'dropped' OR 
    record_type = 'target' OR 
    record_type = 'occupant'
  )
);
