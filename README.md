# Go M.U.D - Getting Started

- [Getting Started](README.md)
- [How to Play](README-HOWTOPLAY.md)
- [API](README-API.md)
- [Design](README-DESIGN.md)
- [TODO](README-TODO.md)

A M.U.D (multi user dungeon) with a back end API built with [`Go`](https://go.dev/) and a front end UI built with [`Flutter`](https://flutter.dev/docs).

❗ No Longer Maintained ~ March 18th 2024 - See [Go Little M.U.D.](https://gitlab.com/alienspaces/go-little-mud
) for more information❗

## Quick Start

📝 _Might be a good idea to look at what these tools do before running them!_

### Running the server

Starts a [`postgres`](https://www.postgresql.org/) database in a [`docker`](https://www.docker.com/) container, runs database migrations with [`db-migrate`](https://db-migrate.readthedocs.io/en/latest/), loads game data and starts the [`Go`](https://go.dev/) API server.

```bash
cd backend
./tools/start
```

API server will be available at [http://localhost:8082/](http://localhost:8082/)

### Running the client

Generate client configuration code.

```bash
cd frontend
./tools/generate.dart
```

Use your favourite Flutter project IDE and Android/iOS emulator.
