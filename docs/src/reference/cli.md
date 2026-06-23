# CLI Commands

This page documents the binmate command-line interface. Running binmate with no
arguments launches the [interactive TUI](/guide/interface).

## `binmate add`

Add a binary from a GitHub release URL or from config.

```bash
# Add from URL
binmate add https://github.com/cli/cli/releases/download/v2.30.0/gh_2.30.0_linux_amd64.tar.gz

# Add from config
binmate add gh
```

## `binmate list`

List all registered binaries, or the versions of a specific binary.

```bash
# List all binaries
binmate list

# List versions of a specific binary
binmate list gh
```

## `binmate install`

Install a specific version of a binary.

```bash
# Install a specific version
binmate install --binary gh --version v2.30.0

# Install the latest version
binmate install --binary gh --version latest
```

| Flag        | Description                                       |
| ----------- | ------------------------------------------------- |
| `--binary`  | The id of the binary to install.                  |
| `--version` | The version to install (use `latest` for newest). |

## `binmate switch`

Switch the active version to a different installed version.

```bash
binmate switch gh v2.29.0
```

## `binmate update`

Update a binary to the latest version.

```bash
binmate update gh
```

## `binmate remove`

Remove a binary from the database.

```bash
# Remove from the database only
binmate remove gh

# Remove the binary and its files
binmate remove gh --files
```

| Flag      | Description                                        |
| --------- | ------------------------------------------------- |
| `--files` | Also remove the installed files from the filesystem. |

## `binmate config`

Display the current configuration.

```bash
# Human-readable output
binmate config

# JSON output
binmate config --json
```

| Flag     | Description                       |
| -------- | --------------------------------- |
| `--json` | Display the configuration as JSON. |

## `binmate sync`

Sync the configuration file with the database.

```bash
binmate sync
```

## `binmate version`

Show version and build information.

```bash
# Show the current version
binmate --version
binmate version

# Show detailed build metadata
binmate version --verbose
```

| Flag        | Description                              |
| ----------- | ---------------------------------------- |
| `--verbose` | Show detailed build metadata.            |
