<p align="center">
  <img alt="SKM logo" src="assets/logo.png" height="150" />
  <h3 align="center">SKM</h3>
  <p align="center">A powerful CLI tool for managing FIDO2 security keys.</p>
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/mohammadv184/skm.svg)](https://pkg.go.dev/github.com/mohammadv184/skm)
[![GitHub release](https://img.shields.io/github/release/mohammadv184/skm.svg)](https://github.com/mohammadv184/skm/releases)
[![License](https://img.shields.io/github/license/mohammadv184/skm)](https://github.com/mohammadv184/skm/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/mohammadv184/skm)](https://goreportcard.com/report/github.com/mohammadv184/skm)
---
**SKM** is a powerful, cross-platform CLI tool designed for managing FIDO2 security keys (like YubiKeys, SoloKeys, and others).


## ✨ Features

- 🔍 **Device Discovery**: Quickly list all connected FIDO2 security keys.
- ℹ️ **Detailed Info**: View technical specifications, including AAGUID, supported protocols, and PIN/UV retry counts.
- 🔑 **Credential Management**: List and delete resident (discoverable) credentials stored on your key.
- 🔐 **PIN Management**: Set and change your device PIN.
- ⚙️ **Device Configuration**: Toggle advanced features like *Always UV* and *Enterprise Attestation*.
- 🧹 **Factory Reset**: Completely wipe and reset your security key to factory settings.


## Installation

### Debian-based Linux distributions
```bash
curl -fsSL https://repo.mohammad-abbasi.me/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/skm.gpg
echo "deb https://repo.mohammad-abbasi.me/apt * *" > /etc/apt/sources.list.d/skm.list
sudo apt update && sudo apt install skm
```
### Fedora / RHEL / CentOS
```bash
echo '[skm]
name=Gloader
baseurl=https://repo.mohammad-abbasi.me/yum
enabled=1
gpgcheck=1
gpgkey=https://repo.mohammad-abbasi.me/yum/gpg.key' | sudo tee /etc/yum.repos.d/skm.repo
sudo yum install skm
```

### Homebrew
```bash
brew install --cask mohammadv184/tap/skm
```

### Go install (Not recommended)
```bash
go install github.com/mohammadv184/skm@latest
```

### Binary builds
You can download the binary builds from the [releases](https://github.com/mohammadv184/skm/releases)

### Deb and RPM packages
You can download the deb and rpm packages from the [releases](https://github.com/mohammadv184/skm/releases)



## Usage


```bash
# List all connected keys
skm list

# Get detailed info about a key
skm info

# Manage credentials
skm creds list
```





## Contributing
Contributions are welcome! Please open issues or pull requests for improvements or bug fixes.

## Security
If you discover any security-related issues, please email mohammad.v184@gmail.com instead of using the issue tracker.


## License
The MIT License (MIT). Please see [License File](LICENSE) for more information.