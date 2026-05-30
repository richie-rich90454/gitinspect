# Installation

## Quick Install (macOS / Linux)

The fastest way to get started — one command:

```bash
curl -fsSL https://raw.githubusercontent.com/richie-rich90454/gitinspect/main/install.sh | bash
```

Set `INSTALL_DIR` to change the install location (default: `/usr/local/bin`):

```bash
curl -fsSL ... | INSTALL_DIR=~/.local/bin bash
```

## Homebrew

```bash
brew tap richie-rich90454/tap
brew install gitinspect
```

## Scoop (Windows)

```powershell
scoop bucket add gitinspect https://github.com/richie-rich90454/scoop-bucket
scoop install gitinspect
```

## Go Install

Requires Go 1.22+:

```bash
go install github.com/richie-rich90454/gitinspect/cmd/gitinspect@latest
```

## Binary Download

Download the latest binary for your platform from the [Releases page](https://github.com/richie-rich90454/gitinspect/releases):

| Platform | Architecture | File |
|----------|-------------|------|
| Linux | amd64 | `gitinspect_X.Y.Z_linux_amd64.tar.gz` |
| Linux | arm64 | `gitinspect_X.Y.Z_linux_arm64.tar.gz` |
| macOS | amd64 | `gitinspect_X.Y.Z_darwin_amd64.tar.gz` |
| macOS | arm64 | `gitinspect_X.Y.Z_darwin_arm64.tar.gz` |
| Windows | amd64 | `gitinspect_X.Y.Z_windows_amd64.zip` |
| Windows | arm64 | `gitinspect_X.Y.Z_windows_arm64.zip` |

Extract and place the binary in your `PATH`:

```bash
# Linux / macOS
tar xzf gitinspect_*_linux_amd64.tar.gz
chmod +x gitinspect
sudo mv gitinspect /usr/local/bin/
```

## Linux Packages

```bash
# Debian / Ubuntu
sudo dpkg -i gitinspect_*_linux_amd64.deb

# RHEL / Fedora
sudo rpm -i gitinspect_*_linux_amd64.rpm

# Alpine
sudo apk add gitinspect_*_linux_amd64.apk
```

## Verify Installation

```bash
gitinspect help
```
