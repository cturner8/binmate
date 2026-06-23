# Plan: Extend GitHub Authentication Options (Issue #66)

## Summary

Extend `internal/providers/github/auth.go` to support a token priority chain (`BINMATE_GITHUB_TOKEN` > `GITHUB_TOKEN` > `BINMATE_GITHUB_ASKPASS` > `GITHUB_ASKPASS`) and update the provider function signatures to accept a pre-built `*http.Client`, enabling the TUI to resolve the token once on startup (avoiding repeated askpass invocations).

---

## Phase 1 — Core Auth Changes

**1.1 Rewrite `internal/providers/github/auth.go`**
- Add private `resolveAskpass(scriptPath string) (string, error)`: executes the script (no args), captures stdout, trims whitespace. Returns error on non-zero exit or empty output.
- Add exported `ResolveToken() (string, error)`: tries sources in order, returns first non-empty token or `""` if none found:
  1. `resolveAskpass(os.Getenv("BINMATE_GITHUB_ASKPASS"))` — only if env var is set
  2. `resolveAskpass(os.Getenv("GITHUB_ASKPASS"))` — only if env var is set
  3. `os.Getenv("BINMATE_GITHUB_TOKEN")`
  4. `os.Getenv("GITHUB_TOKEN")`
- Change `CreateHTTPClient(authenticated bool)` signature to `CreateHTTPClient(token string) (*http.Client, error)` — empty token returns a plain unauthenticated client.
- Add `NewClientForBinary(binary *database.Binary) (*http.Client, error)` convenience function: if `!binary.Authenticated` returns plain client, otherwise calls `ResolveToken()` and `CreateHTTPClient(token)`. Returns error if auth requested but token empty.

**1.2 Create `internal/providers/github/auth_test.go`** (new file)
- Test token priority: BINMATE_GITHUB_TOKEN > GITHUB_TOKEN
- Test askpass priority: BINMATE_GITHUB_ASKPASS > GITHUB_ASKPASS
- Test `resolveAskpass` with a small shell script that echoes a known token
- Test `CreateHTTPClient("")` returns plain client, `CreateHTTPClient("tok")` sets Authorization header
- Test `NewClientForBinary` with `Authenticated=false` and `Authenticated=true`

---

## Phase 2 — Provider Function Signature Changes

Remove internal `CreateHTTPClient` calls from provider functions. All public functions now accept `client *http.Client` as their first parameter:

**2.1 `internal/providers/github/release_info.go`**
- `FetchReleaseNotes(client *http.Client, binary *database.Binary, version string) (ReleaseInfo, error)` — remove internal `CreateHTTPClient(binary.Authenticated)` call
- `ListAvailableVersions(client *http.Client, binary *database.Binary, limit int) ([]ReleaseInfo, error)` — same
- `GetRepositoryInfo(client *http.Client, binary *database.Binary) (RepositoryInfo, error)` — same
- `StarRepository(client *http.Client, binary *database.Binary) error` — same (was hard-coded `authenticated=true`; caller now responsible for passing authenticated client)

**2.2 `internal/providers/github/fetch_release_asset.go`**
- `FetchReleaseAsset(client *http.Client, binary *database.Binary, version string) (Release, ReleaseAsset, error)`

**2.3 `internal/providers/github/download_asset.go`**
- `DownloadAsset(client *http.Client, providerPath string, assetId int, assetName string) (string, error)` — remove `authenticated bool` param, use provided client directly

---

## Phase 3 — Update CLI Callers

**3.1 `internal/core/install/service.go`**
- Before calling provider functions, call `github.NewClientForBinary(binaryConfig)` to get `*http.Client`
- Pass client to `github.FetchReleaseAsset(client, binaryConfig, version)` and `github.DownloadAsset(client, ...)`

**3.2 `internal/cli/check/command.go`**
- Same pattern: call `github.NewClientForBinary(binaryConfig)` before the two `FetchReleaseAsset` call sites (lines 64 and 99)

---

## Phase 4 — Config: askpassMode Option

**4.1 `internal/core/config/config.go`**
- Add `AskpassMode string \`mapstructure:"askpassMode"\`` to `ProviderDefaults` struct
- Document valid values: `"startup"` (default) and `"always"`

**4.2 `docs/src/public/schema.json`**
- Add `askpassMode` property to the `providerDefaults` definition with `"type": "string"`, `"enum": ["startup", "always"]`, description, and `"default": "startup"`

---

## Phase 5 — TUI Startup Token Resolution

**5.1 `internal/tui/messages.go`**
- Add `githubTokenResolvedMsg struct { token string; err error }` message type

**5.2 `internal/tui/model.go`**
- Add `resolvedGithubToken string` field (empty = no token/not yet resolved)
- Add `githubTokenResolved bool` flag to track whether resolution has been attempted

**5.3 `internal/tui/init.go`**
- In `Init()`, check if `askpassMode == "startup"` (or default empty → treat as startup). If so, add `resolveGithubTokenCmd(m.config)` to the `tea.Batch`
- `resolveGithubTokenCmd` is a new `tea.Cmd` function that calls `github.ResolveToken()` and returns `githubTokenResolvedMsg`
- If askpass mode is not "startup", skip (token will be resolved per-call in helpers)

**5.4 `internal/tui/update.go`**
- Handle `githubTokenResolvedMsg` in the main `Update()` switch: set `m.resolvedGithubToken = msg.token` and `m.githubTokenResolved = true`. If `msg.err != nil`, log it but don't crash (unauthenticated binaries still work).
- Add helper `(m model) clientForBinary(binary *database.Binary) (*http.Client, error)`:
  - If `!binary.Authenticated`: return `&http.Client{}`
  - If `m.githubTokenResolved`: return `github.CreateHTTPClient(m.resolvedGithubToken)`
  - Otherwise (per-call mode or startup not yet complete): call `github.NewClientForBinary(binary)`
- Update TUI command functions to accept and pass `*http.Client`:
  - `fetchReleaseNotes(client *http.Client, ...)` → `github.FetchReleaseNotes(client, binary, version)`
  - `fetchRepositoryInfo(client *http.Client, ...)` → `github.GetRepositoryInfo(client, binary)`
  - `fetchAvailableVersions(client *http.Client, ...)` → `github.ListAvailableVersions(client, binary, 20)`
  - `starRepository(client *http.Client, ...)` → `github.StarRepository(client, binary)`
- At each call site (keyReleaseNotes, keyRepoInfo, keyAvailVersions, keySwitch), resolve client via `m.clientForBinary(m.selectedBinary)` before dispatching the command

---

## Phase 6 — Documentation Updates

- New dedicated "authentication" page in "Guide" section of docs, explaining:
  - Token priority chain
  - Askpass script usage
  - `askpassMode` config option
  - Token API requirements (scopes, etc.)
- Update existing pages following updated config schema
- Add a "Tip" box to top of "Usage" page explaining that GitHub authentication is supported and recommended to avoid rate limits, should link to new authentication page

---


## Relevant Files

- `internal/providers/github/auth.go` — Core rewrite: new token priority chain + signature changes
- `internal/providers/github/auth_test.go` — New: comprehensive tests for auth functions
- `internal/providers/github/release_info.go` — Add `client *http.Client` param, remove internal client creation
- `internal/providers/github/fetch_release_asset.go` — Add `client *http.Client` param
- `internal/providers/github/download_asset.go` — Replace `authenticated bool` with `client *http.Client`
- `internal/core/install/service.go` — Call `NewClientForBinary` before provider calls (lines 36, 44)
- `internal/cli/check/command.go` — Same pattern at lines 64, 99
- `internal/core/config/config.go` — Add `AskpassMode` to `ProviderDefaults`
- `docs/src/public/schema.json` — Add `askpassMode` to `providerDefaults` schema
- `internal/tui/model.go` — Add `resolvedGithubToken`, `githubTokenResolved` fields
- `internal/tui/messages.go` — Add `githubTokenResolvedMsg`
- `internal/tui/init.go` — Conditionally resolve token on startup
- `internal/tui/update.go` — Handle resolution message, update TUI helper signatures and call sites

---

## Verification

1. `go build ./...` — no compile errors after all signature changes are propagated
2. `go test ./internal/providers/github/...` — all new auth tests pass, existing filter/download tests unaffected
3. `go test ./...` — full test suite green
4. `go fmt ./internal/providers/github/... ./internal/core/install/... ./internal/cli/check/... ./internal/tui/...`
5. Manual: `BINMATE_GITHUB_TOKEN=testtoken ./binmate install gh` — token used in request headers
6. Manual: `GITHUB_TOKEN=testtoken BINMATE_GITHUB_TOKEN=override ./binmate install gh` — `override` takes precedence
7. Manual: create a simple shell script `echo "my-token"`, set `BINMATE_GITHUB_ASKPASS=/path/to/script`, verify token resolved in TUI startup
8. Manual: set `providers.github.askpassMode: "always"` in config, verify askpass script is invoked on each GitHub API action in TUI
9. `go vet ./...`

---

## Decisions

- `askpassMode` only relevant for TUI (CLI is stateless per invocation so always resolves fresh per command)
- `askpassMode` default is `"startup"` per issue guidance ("could default to once on startup")
- If askpass resolution fails at startup, TUI continues normally; error surfaces only when an authenticated binary API call is made
- `StarRepository` previously hard-coded `authenticated=true`; the caller (TUI `starRepository` helper) now creates an authenticated client via `clientForBinary`, preserving correct behaviour
- `resolveAskpass` calls the script with no arguments (binmate doesn't need interactive prompting, just the token)

## Excluded Scope

- No support for other providers (GitLab, etc.) in this issue
- No config file option for specifying token value inline (env vars / askpass only)
- No UI feedback in TUI for successful token resolution (silent background init)
