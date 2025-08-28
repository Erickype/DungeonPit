# DungeonPit — Server-Driven Procedural Dungeon

A Go-based gRPC service for procedural dungeon generation, driving rendering in Unreal Engine. 
This separation allows efficient algorithms, server-side determinism, and lightweight clients.

# Core Highlights
## 1. Server Architecture

- Runs as a Go gRPC server ``(server/server.go)``
- GenerateDungeon RPC; Unreal client receives structured dungeon data at startup

## 2. Modular Generation Pipeline

- Implemented inside internal/core
- Room placement (non-overlapping, buffered)
- Delaunay triangulation for room adjacency
- Minimum Spanning Tree (Prim’s algorithm) for path backbone
- A* pathfinding to carve hallways into a grid, complete with heuristics
- Line base data for rendering

## 3. Mapping Layer

- Located under pkg/mapper, cleanly translates between:
- Core domain types (e.g., Line2D, GridLine)
- Proto-generated messages (in pkg/gen/*)
- Avoids coupling and eases future refactoring.

## 4. Organized Protocol Buffers

- Maintained under proto/:
- Cleanly groups messages per domain (e.g., common, dungeon, game)
- go_package declarations align with generated directories (pkg/gen/...)

# Setup & Usage
```
# Clone the repo
git clone https://github.com/Erickype/DungeonPit.git
cd DungeonPit

# Generate/update gRPC code
make generate-proto

# Development database using docker
docker-compose up

# Run the dungeon generation server
go run server/server.go
```

# Why This Architecture?
| Requirement         | Design Choice                            | Benefit                                 |
|---------------------|------------------------------------------|-----------------------------------------|
| Reliable generation | Move logic server-side                   | Ensures determinism and version control |
| Clean separation    | Mapper layer between core and gRPC       | Easier domain evolution and testing     |
| Future-ready        | Organized proto files and codegen target | Maintainable and extendable structure   |
| Developer-friendly  | Makefile + structured packages           | Fast onboarding for teammates           |

# What's Next

- Unreal client integration: consume RPC, render dungeon tiles and hallways
- Art & modular assets: replace debug placeholders with polished mesh assets
- Multiplayer: server-prescribed dungeon with player network sync