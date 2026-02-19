<p align="center">
  <img alt="SKM logo" src="assets/logo.png" height="150" />
  <h3 align="center">SKM</h3>
  <p align="center">A powerful CLI tool for managing FIDO2 security keys with ease and style.</p>
</p>

[![Go Version](https://img.shields.io/github/go-mod/go-version/mohammadv184/skm)](https://go.dev/)
[![License](https://img.shields.io/github/license/mohammadv184/skm)](LICENSE)
---
**SKM** is a powerful, cross-platform CLI tool designed for managing FIDO2 security keys (like YubiKeys, SoloKeys, and others). Built with a focus on both user experience and automation, it provides a rich Terminal User Interface (TUI) for interactive use and comprehensive flags for scripting.


## ✨ Features

- 🔍 **Device Discovery**: Quickly list all connected FIDO2 security keys.
- ℹ️ **Detailed Info**: View technical specifications, including AAGUID, supported protocols, and PIN/UV retry counts.
- 🔑 **Credential Management**: List and delete resident (discoverable) credentials stored on your key.
- 🔐 **PIN Management**: Set, change, and monitor PIN retries.
- ⚙️ **Device Configuration**: Toggle advanced features like *Always UV* and *Enterprise Attestation*.
- 🧹 **Factory Reset**: Completely wipe and reset your security key to factory settings.
- 🎨 **Beautiful UI**: Powered by [Charmbracelet](https://charm.sh/) (`bubbletea`, `lipgloss`) for a modern terminal experience.
- 🤖 **Automation Friendly**: Every command supports non-interactive execution via flags.

---

## 🚀 Installation

### From Source

Ensure you have [Go](https://go.dev/doc/install) installed (version 1.21 or later).

```bash
git clone https://github.com/mohammadv184/skm.git
cd skm
make install
```

---

## 🛠️ Usage

SKM is designed to be intuitive. Simply running a command without arguments will usually trigger an interactive prompt.

### Quick Start

```bash
# List all connected keys
skm list

# Get detailed info about a key
skm info

# Manage credentials
skm creds list
```

### Command Overview

| Command | Alias | Description |
| :--- | :--- | :--- |
| `list` | `ls` | List connected security keys |
| `info` | - | Show detailed device information |
| `creds` | `c` | Manage resident credentials (list/delete) |
| `pin` | `p` | Manage device PIN (set/change/retries) |
| `config` | `cfg` | Manage device configuration (always-uv, etc.) |
| `reset` | - | Factory reset the security key |

---

## 🤖 Automation & Scripting

For non-interactive use, SKM provides a consistent set of flags across all commands.

### Examples

**View info for a specific device:**
```bash
skm info --device-path /dev/hidraw0
```

**List credentials without a prompt:**
```bash
skm creds list --device-path /dev/hidraw0 --pin 123456
```

**Delete a specific credential:**
```bash
skm creds delete --device-path /dev/hidraw0 --pin 123456 --credential-id <BASE64_ID>
```

**Reset a device silently:**
```bash
skm reset --device-path /dev/hidraw0 --yes
```

---

## 🛠️ Development

### Building
```bash
make build
```

### Linting
```bash
make lint
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">Made with ❤️ for the Go & Security community.</p>
