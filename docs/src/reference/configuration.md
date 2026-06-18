# Configuration Reference

binmate is configured through a JSON file located at
`~/.config/.binmate/config.json`. This page documents every field. For a guided
introduction, see the [configuration guide](/guide/configuration).

## Top-level fields

| Field       | Type   | Description                                                       |
| ----------- | ------ | ---------------------------------------------------------------- |
| `$schema`   | string | URL of the JSON schema used for validation.                      |
| `version`   | number | Configuration schema version (currently `1`).                    |
| `global`    | object | Optional global defaults applied to all binaries.                |
| `binaries`  | array  | The list of binaries managed by binmate.                         |

## Global configuration

Global settings apply to all binaries unless overridden on an individual binary.

| Field                                       | Type    | Description                                                  |
| ------------------------------------------- | ------- | ----------------------------------------------------------- |
| `global.installPath`                        | string  | Default installation path for all binaries (e.g. `/usr/local/bin`). |
| `global.providers.<provider>.authenticated` | boolean | Default authentication setting for a provider.              |

## Binary configuration

Each entry in the `binaries` array describes a single binary.

| Field           | Type    | Required | Description                                                       |
| --------------- | ------- | -------- | ---------------------------------------------------------------- |
| `id`            | string  | Yes      | Unique identifier for the binary.                                |
| `name`          | string  | Yes      | Display name of the binary.                                      |
| `provider`      | string  | Yes      | Provider type (currently only `github` is supported).            |
| `path`          | string  | Yes      | Repository path (e.g. `owner/repo`).                             |
| `format`        | string  | Yes      | Archive format (`.tar.gz`, `.zip`, `.tgz`).                      |
| `installPath`   | string  | No       | Custom installation path (overrides `global.installPath`).       |
| `assetRegex`    | string  | No       | Regex to filter release assets.                                  |
| `releaseRegex`  | string  | No       | Regex to filter releases.                                        |
| `authenticated` | boolean | No       | Use authentication for API calls (overrides the provider default). |

## Supported values

- **Providers**: `github` (currently the only supported provider).
- **Formats**: `.tar.gz`, `.zip`, `.tgz`.

## Example

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
