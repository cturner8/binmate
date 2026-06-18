# Architecture

binmate follows a layered architecture that separates the command-line interface, core
business logic, persistence, provider integrations, and the terminal UI.

## Project structure

```
cmd/                    # Command entry point
internal/
  cli/                  # CLI command definitions
    add/                # Add binary command
    config/             # Config command
    import/             # Import command
    install/            # Install command
    list/               # List command
    remove/             # Remove command
    switch/             # Switch version command
    sync/               # Sync config command
    update/             # Update command
  core/                 # Core business logic
    binary/             # Binary management service
    config/             # Configuration management
    crypto/             # Checksum verification
    install/            # Installation and extraction
    url/                # GitHub URL parsing
    version/            # Version management service
  database/             # SQLite data layer
    repository/         # Data access repositories
  providers/            # External provider integrations
    github/             # GitHub releases API
  tui/                  # Terminal UI (Bubble Tea)
```

## Layers

- **CLI (`internal/cli`)** — Defines each command and parses user input. The root
  command launches the TUI; the remaining commands provide scriptable equivalents.
- **Core (`internal/core`)** — The business logic: managing binaries and versions,
  parsing GitHub URLs, extracting archives, and verifying checksums.
- **Database (`internal/database`)** — A SQLite data layer with a repository per
  entity. See the [database reference](/reference/database) for details.
- **Providers (`internal/providers`)** — Integrations with external sources. GitHub is
  currently the only supported provider.
- **TUI (`internal/tui`)** — The interactive Terminal UI, built with
  [Bubble Tea](https://github.com/charmbracelet/bubbletea) and following the
  Elm Architecture (Model–Update–View).

## Technology

binmate is written in [Go](https://go.dev/) and builds on a small set of well-known
libraries:

- **[Cobra](https://github.com/spf13/cobra)** — CLI framework for command structure.
- **[Viper](https://github.com/spf13/viper)** — Configuration management.
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** — Terminal UI framework.
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** — Terminal styling.
- **SQLite** — Embedded database for state persistence.
