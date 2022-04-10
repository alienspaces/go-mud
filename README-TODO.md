# Go M.U.D -TODO

- [Getting Started](README.md)
- [How to Play](README-HOWTOPLAY.md)
- [API](README-API.md)
- [Design](README-DESIGN.md)
- [TODO](README-TODO.md)

## Now

- Server: Dungeon instances
  - Spawn a new dungeon instance that can occupy max N characters
  - Delete old dungeon instances 5 minutes after they are empty of characters

## Next

- Server: Implement 2 second turns
- Client: Implement 2 second progress bar and retries

## Later

- Server: Implement monster movement
- Server: Effects
  - Object effects (passive)
  - Object and spells damage effects (active)
- Server: Object and monster respawning
- Client: Icons for room contents
- Server: Support multiple different servers that belong to the same service to support
  building a separate server for managing monster movement, dungeon cleanup etc
  - Rename `server/service/game/internal/cli/runner` to `server/service/game/internal/cli/api`
  - Rename `server/service/game/internal/server/runner` to `server/service/game/internal/server/api`
- All: Docker Compose
