# Nest M.U.D - API

- [Getting Started](README.md)
- [How to Play](README-HOWTOPLAY.md)
- [API](README-API.md)
- [Design](README-DESIGN.md)
- [TODO](README-TODO.md)

## Dungeons

**List dungeons:**

- [Response Schema](server/schema/docs/dungeon/response.schema.json)

```bash
GET /api/v1/dungeons
```

**Get dungeon:**

- [Response Schema](server/schema/docs/dungeon/response.schema.json)

```bash
GET /api/v1/dungeons/{:dungeon_id}
```

## Locations

**List dungeon locations:**

📝 _Unimplemented_

- [Response Schema](server/schema/docs/location/response.schema.json)

```bash
GET /api/v1/dungeons/{:dungeon_id}/locations
```

**Get dungeon location:**

📝 _Unimplemented_

- [Response Schema](server/schema/docs/location/response.schema.json)

```bash
GET /api/v1/dungeons/{:dungeon_id}/locations/{:location_id}
```

## Characters

**Create character:**

- [Request Schema](server/schema/docs/character/create.request.schema.json)
- [Response Schema](server/schema/docs/character/response.schema.json)

```bash
POST /api/v1/characters
```

**List characters:**

- [Response Schema](server/schema/docs/character/response.schema.json)

```bash
GET /api/v1/characters
```

**Get character:**

- [Response Schema](server/schema/docs/character/response.schema.json)

```bash
GET /api/v1/characters/{:character_id}
```

**Update a character:**

📝 _Unimplemented_

- [Request Schema](server/schema/docs/character/response.schema.json)
- [Response Schema](server/schema/docs/character/response.schema.json)

```bash
PUT /api/v1/characters/{:character_id}
```

## Dungeon Instances

Dungeon instances are created to accomodate a maximum number of characters per dungeon.

**Enter dungeon:**

Entering a dungeon returns the specific dungeon instance the character entered.

- [Response Schema](server/schema/docs/dungeoninstance/response.schema.json)

```bash
POST /api/v1/dungeons/{:dungeon_id}/enter
```

**List dungeon instances:**

Lists the currently running dungeon instances.

- [Response Schema](server/schema/docs/dungeoninstance/response.schema.json)

```bash
GET /api/v1/dungeons/{:dungeon_id}/instances
```

**Get dungeon instances:**

Returns a specific running dungeon instance.

- [Response Schema](server/schema/docs/dungeoninstance/response.schema.json)

```bash
GET /api/v1/dungeons/{:dungeon_id}/instances/{:dungeon_instance_id}
```

## Actions

Characters are controlled by performing actions.

**Create a character action:**

- [Request Schema](server/schema/docs/action/create.request.schema.json)
- [Response Schema](server/schema/docs/action/response.schema.json)

```bash
POST /api/v1/dungeons/{:dungeon_id}/instances/{:dungeon_instance_id}/actions
```
