# GitHub Authentication

binmate can authenticate with the GitHub API to avoid rate limits and access private
repositories. This page explains the supported authentication methods and how to
configure them.

## Why authenticate?

The GitHub API enforces rate limits on unauthenticated requests (60 requests per
hour per IP address). When browsing binaries in the TUI, checking for new versions,
or listing available releases, binmate may issue several API calls in quick
succession. Authenticating with a personal access token raises the limit to 5 000
requests per hour.

## Token priority chain

When authentication is required, binmate resolves a token by trying the following
sources in order. The first non-empty result is used:

1. **`BINMATE_GITHUB_ASKPASS` script** - executes the script and captures its
   standard output as the token.
2. **`GITHUB_ASKPASS` script** - same behaviour; used as a fallback when
   `BINMATE_GITHUB_ASKPASS` is not set.
3. **`BINMATE_GITHUB_TOKEN`** - a static token specific to binmate.
4. **`GITHUB_TOKEN`** - the standard GitHub token variable (commonly set in CI
   environments).

If none of the sources produces a token, and a binary is configured with
`"authenticated": true`, binmate returns an error rather than silently making an
unauthenticated request.

## Setting a static token

Export one of the token environment variables before running binmate:

```bash
export BINMATE_GITHUB_TOKEN=ghp_...
binmate
```

Or add it to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) so it is available
in every session.

`BINMATE_GITHUB_TOKEN` takes precedence over `GITHUB_TOKEN`, which makes it useful
when you need a more restrictive token for binmate alongside a broader token used by
other tools.

## Using an askpass script

An askpass script is an executable that prints a token to standard output. This
allows binmate to retrieve credentials from a password manager or secrets store at
run time - without storing the token in a plain-text environment variable.

Create a script that prints a token and nothing else:

```bash
#!/usr/bin/env bash
# ~/.local/bin/binmate-askpass
# Retrieve a GitHub token from your password manager.
op item get "GitHub Token" --fields password
```

Make the script executable and point binmate to it:

```bash
chmod +x ~/.local/bin/binmate-askpass
export BINMATE_GITHUB_ASKPASS=~/.local/bin/binmate-askpass
```

binmate runs the script with no arguments and uses the trimmed output as the token.
The script must exit with status 0 and print a non-empty token; any other outcome is
treated as an error.

::: tip Prefer `BINMATE_GITHUB_ASKPASS`
Use `BINMATE_GITHUB_ASKPASS` rather than `GITHUB_ASKPASS` when the token is
specifically for binmate. This avoids conflicts with other tools that also respect
`GITHUB_ASKPASS`.
:::

## Configuring authentication in `config.json`

Set `authenticated` to `true` on a binary (or globally) to tell binmate that API
calls for that binary should always use an authenticated client:

```json
{
  "$schema": "https://binmate.cturner8.dev/schema.json",
  "version": 1,
  "global": {
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
    }
  ]
}
```

### `askpassMode`

The `askpassMode` option under `global.providers.github` controls when the askpass
script is invoked during TUI sessions:

| Value | Behaviour |
| --- | --- |
| `"startup"` (default) | The token is resolved once when the TUI starts. The result is cached and reused for all subsequent API calls in that session. |
| `"always"` | The token is resolved fresh on every GitHub API action. Use this if your askpass script rotates short-lived tokens. |

```json
{
  "$schema": "https://binmate.cturner8.dev/schema.json",
  "version": 1,
  "global": {
    "providers": {
      "github": {
        "authenticated": true,
        "askpassMode": "startup"
      }
    }
  },
  "binaries": []
}
```

::: tip CLI commands are always per-invocation
`askpassMode` only affects TUI sessions. Each CLI command (e.g. `binmate install`)
is a separate process and resolves the token once per invocation regardless of this
setting.
:::

## Required token scopes

Create a fine-grained personal access token (or a classic token) with the following
permissions:

| Permission | Reason |
| --- | --- |
| `public_repo` (classic) or "Contents" read (fine-grained) | Required to list releases and download assets from public repositories. |
| `repo` (classic) or "Contents" read + "Metadata" read (fine-grained) | Required for private repositories. |
| `public_repo` or `repo` write | Required only if you use `binmate star` to star repositories. |

For read-only usage with public repositories, a classic token with `public_repo` is
sufficient.

## Next steps

- Return to the [usage guide](/guide/usage) for everyday commands.
- See the [configuration reference](/reference/configuration) for a complete field
  listing.
