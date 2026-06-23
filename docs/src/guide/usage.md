# Usage

binmate can be driven either through its interactive TUI or via individual CLI
commands. This page covers everyday usage; see [the interface](/guide/interface) for a
tour of the TUI and the [CLI reference](/reference/cli) for the full command list.

::: tip GitHub authentication recommended
Authenticating with GitHub avoids API rate limits and unlocks private repository
access. See the [authentication guide](/guide/authentication) to set up a token or
askpass script.
:::

## Interactive mode

Launch the TUI for interactive management by running binmate with no arguments:

```bash
binmate
```

This opens the [Binary List View](/guide/interface#binary-list-view) where you can
browse, search, filter, and manage your binaries.

## CLI commands

### Add a binary

Add a binary from a GitHub release URL or from config:

```bash
# Add from URL
binmate add https://github.com/cli/cli/releases/download/v2.30.0/gh_2.30.0_linux_amd64.tar.gz

# Add from config
binmate add gh
```

### List binaries

List all registered binaries:

```bash
binmate list
```

List versions of a specific binary:

```bash
binmate list gh
```

### Install a binary

Install a specific version of a binary:

```bash
binmate install --binary gh --version v2.30.0
```

Install the latest version:

```bash
binmate install --binary gh --version latest
```

### Switch versions

Switch to a different installed version:

```bash
binmate switch gh v2.29.0
```

### Update to latest

Update a binary to the latest version:

```bash
binmate update gh
```

### Remove a binary

Remove a binary from the database:

```bash
binmate remove gh
```

Remove a binary and its files:

```bash
binmate remove gh --files
```

### View configuration

Display the current configuration:

```bash
binmate config
```

Display configuration as JSON:

```bash
binmate config --json
```

### Sync configuration

Sync the configuration file with the database:

```bash
binmate sync
```

### Version information

Show the current binmate version:

```bash
binmate --version
binmate version
```

Show detailed build metadata:

```bash
binmate version --verbose
```

## Next steps

- Browse the full [CLI reference](/reference/cli).
- Learn how to [configure binmate](/guide/configuration).
