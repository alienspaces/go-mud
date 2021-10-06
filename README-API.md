# Nest M.U.D - API

- [Getting Started](README.md)
- [How to Play](README-HOWTOPLAY.md)
- [API](README-API.md)
- [Design](README-DESIGN.md)

## Dungeons

**List dungeons:**

- [Response Schema](server/src/controllers/dungeon/schema/dungeon.schema.json)

```bash
GET /api/v1/dungeons
```

**Get a dungeon:**

- [Response Schema](server/src/controllers/dungeon/schema/dungeon.schema.json)

```bash
GET /api/v1/dungeons/{:dungeon_id}
```

## Characters

**Create a character:**

- [Request Schema](server/src/controllers/dungeon-character/schema/create-dungeon-character.schema.json)
- [Response Schema](server/src/controllers/dungeon-character/schema/dungeon-character.schema.json)

```bash
POST /api/v1/dungeons/{:dungeon_id}/characters
```

**Update a character:**

📝 _Unimplemented_

- [Request Schema](server/src/controllers/dungeon-character/schema/update-dungeon-character.schema.json)
- [Response Schema](server/src/controllers/dungeon-character/schema/dungeon-character.schema.json)

```bash
PUT /api/v1/dungeons/{:dungeon_id}/characters/{:character_id}
```

**Get a character:**

- [Response Schema](server/src/controllers/dungeon-character/schema/dungeon-character.schema.json)

```bash
GET /api/v1/dungeons/{:dungeon_id}/characters/{:character_id}
```

## Locations

**List dungeon locations:**

- [Response Schema](server/src/controllers/dungeon-location/schema/dungeon-location.schema.json)

```bash
GET /api/v1/dungeons/{:dungeon_id}/locations
```

**Get a dungeon location:**

- [Response Schema](server/src/controllers/dungeon-location/schema/dungeon-location.schema.json)

```bash
GET /api/v1/dungeons/{:dungeon_id}/locations/{:location_id}
```

## Actions

Characters are controlled with actions that are simple sentences.

**Create a character action:**

- [Request Schema](server/src/controllers/dungeon-character-action/schema/create-dungeon-character-action.schema.json)
- [Response Schema](server/src/controllers/dungeon-character-action/schema/dungeon-character-action.schema.json)

```bash
POST /api/v1/dungeons/{:dungeon_id}/actions
```
