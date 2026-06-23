# Installation

There are several ways to install binmate. Using the install script is the
recommended approach for most users.

## Using the install script (recommended)

### Unix (Linux/macOS)

Install the latest version using the install script:

```bash
curl -fsSL https://binmate.cturner8.dev/install.sh | bash
```

By default, the installer auto-imports the installed `binmate` binary for
self-management using its resolved release URL and version.

Install a specific version:

```bash
curl -fsSL https://binmate.cturner8.dev/install.sh | BINMATE_VERSION=v1.0.0 bash
```

Install to a custom directory:

```bash
curl -fsSL https://binmate.cturner8.dev/install.sh | BINMATE_INSTALL_DIR=$HOME/.local/bin bash
```

Skip automatic post-install self-import:

```bash
curl -fsSL https://binmate.cturner8.dev/install.sh | BINMATE_SKIP_AUTO_IMPORT=1 bash
```

### Windows (PowerShell)

Install the latest version using the PowerShell install script:

```powershell
irm https://binmate.cturner8.dev/install.ps1 | iex
```

Install a specific version:

```powershell
$env:BINMATE_VERSION = "v1.0.0"
irm https://binmate.cturner8.dev/install.ps1 | iex
```

Install to a custom directory:

```powershell
$env:BINMATE_INSTALL_DIR = "$env:LOCALAPPDATA\binmate\bin"
irm https://binmate.cturner8.dev/install.ps1 | iex
```

Skip automatic post-install self-import:

```powershell
$env:BINMATE_SKIP_AUTO_IMPORT = "1"
irm https://binmate.cturner8.dev/install.ps1 | iex
```

## Install script options

The install scripts respect the following environment variables:

| Variable                  | Description                                                        |
| ------------------------- | ----------------------------------------------------------------- |
| `BINMATE_VERSION`         | Install a specific version (e.g. `v1.0.0`). Defaults to `latest`. |
| `BINMATE_INSTALL_DIR`     | Directory to install the `binmate` binary into.                   |
| `BINMATE_SKIP_AUTO_IMPORT`| Set to `1` to skip importing `binmate` for self-management.       |

## Manual installation

Download the appropriate binary for your platform from the
[releases page](https://github.com/cturner8/binmate/releases), extract it, and place
it somewhere on your `PATH`.

## Building from source

binmate is written in Go. To build it yourself:

```bash
git clone https://github.com/cturner8/binmate.git
cd binmate
go build -o binmate .
```

::: tip Requirements
Building from source requires the Go toolchain (Go 1.25 or newer). SQLite support is
provided via CGO, so a C compiler must be available.
:::

## Next steps

Once installed, head over to the [usage guide](/guide/usage) to start managing your
binaries.
