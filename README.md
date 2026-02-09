# Venvi - EU Event Suggestion Platform

A PocketBase-powered event aggregator that discovers hackathons, meetups, and cultural events across Europe.

## Features

- 🔄 **Multi-source Aggregation**: Pulls events from Open Data Hub and Euro Hackathons
- 🚀 **Single Executable**: Built with PocketBase for easy deployment
- 🎨 **HTMX Frontend**: Dynamic, server-rendered UI with minimal JavaScript
- 📅 **Scheduled Sync**: Automatic event updates every 6 hours
- 🔌 **Extensible**: Easy to add new event providers

## Quick Start

### Prerequisites

- [Nix](https://nixos.org/download.html) with Flakes enabled.

### Development Rules

> [!IMPORTANT]
> **Strict Nix Usage Enforced**
> 1.  **Use Nix**: Always run `nix develop` to enter the development environment. This ensures you have the correct versions of all tools.
> 2.  **Manage Tools**: If you need a new tool, **add it to `flake.nix`**. Do not install tools globally or manually.

### Development

```bash
# Enter development environment
nix develop

# Build the application
go build -o venvi .

# Run the server
./venvi serve --http=localhost:8090
```

### First Run

1. Open http://localhost:8090/_/ to access the admin panel
2. Create your admin account
3. The `events` collection will be created automatically
4. Visit http://localhost:8090/ to see the frontend
5. Click "Sync & Refresh" to fetch events

## Project Structure

```
venvi/
├── main.go              # Application entry point
├── providers/           # Event data source implementations
│   ├── provider.go      # Base interface
│   ├── odh.go           # Open Data Hub provider
│   ├── euro_hackathons.go
│   └── sync.go          # Sync orchestrator
├── routes/              # HTTP route handlers
│   ├── web.go           # HTMX/web routes
│   └── api.go           # JSON API routes
├── views/               # Go html/templates
│   ├── layout.html
│   ├── index.html
│   └── partials/
├── pb_public/           # Static assets
├── pb_migrations/       # Database migrations
└── scripts/
    └── validate.sh      # Verification script
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Homepage with HTMX |
| GET | `/partials/events` | Event list partial |
| GET | `/api/venvi/events` | List events (JSON) |
| GET | `/api/venvi/events?category=hackathon` | Filter by category |
| GET | `/api/venvi/events?source=odh` | Filter by source |
| POST | `/api/venvi/sync` | Trigger manual sync |
| GET | `/api/venvi/health` | Health check |

## Adding a New Provider

1. Create `providers/new_source.go` implementing `EventProvider`
2. Add to `Providers` slice in `providers/sync.go`
3. Run tests: `go test ./providers/...`

## Development Commands

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run tests
go test -v -cover ./...

# Full validation
./scripts/validate.sh
```

## License

MIT
