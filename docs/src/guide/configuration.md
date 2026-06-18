# Configuration

Configuration is stored in `~/.config/.binmate/config.json`. It defines the binaries
binmate knows about, along with optional global defaults.

A full breakdown of every field is available in the
[configuration reference](/reference/configuration).

## Basic configuration

```json
{
  "$schema": "https://binmate.cturner8.dev/schema.json",
  "version": 1,
  "binaries": [
    {
      "id": "gh",
      "name": "gh",
      "provider": "github",
      "path": "cli/cli",
      "format": ".tar.gz"
    }
  ]
}
```

::: tip VS Code schema validation
If using VS Code, you'll need to allow use of the binmate JSON schema in
`.vscode/settings.json`:

```jsonc
// .vscode/settings.json
{
  "json.schemaDownload.trustedDomains": {
    "https://binmate.cturner8.dev": true
  }
}
```
:::

## Global configuration

You can define global defaults that apply to all binaries unless overridden:

```json
{
  "$schema": "https://binmate.cturner8.dev/schema.json",
  "version": 1,
  "global": {
    "installPath": "/usr/local/bin",
    "providers": {
      "github": {
        "authenticated": true
      }
    }
  },
  "binaries": [
    {
      "id": "gh",
      "name": "gh",
      "provider": "github",
      "path": "cli/cli",
      "format": ".tar.gz"
    },
    {
      "id": "fzf",
      "name": "fzf",
      "provider": "github",
      "path": "junegunn/fzf",
      "format": ".tar.gz",
      "installPath": "/opt/bin"
    }
  ]
}
```

In this example:

- All binaries use `/usr/local/bin` as the install path by default.
- All binaries use GitHub authentication by default to avoid rate limits.
- The `fzf` binary overrides the global install path with `/opt/bin`.

## Syncing changes

After editing the configuration file by hand, run `binmate sync` to reconcile the
configuration with the database:

```bash
binmate sync
```

## Next steps

- See the [configuration reference](/reference/configuration) for every available
  field.
- Learn how binmate stores state in the [database](/reference/database).
