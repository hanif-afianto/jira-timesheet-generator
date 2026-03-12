# 🛠️ JTG (Jira Timesheet Generator)

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey.svg)](#-installation)

**JTG** is a high-performance CLI tool designed to simplify timesheet management. It fetches worklogs directly from the Jira API and transforms them into professional Excel reports, tailored for business submissions.

---

## 📖 Table of Contents
- [✨ Features](#-features)
- [📋 Prerequisites](#-prerequisites)
- [🚀 Installation](#-installation)
  - [Method 1: Using Global Binaries (No Go Required)](#method-1-using-global-binaries-no-go-required)
  - [Method 2: Building from Source](#method-2-building-from-source)
- [⚙️ Configuration](#-configuration)
- [💡 Usage](#-usage)

---

## ✨ Features

- 🔄 **Automated Sync**: Seamlessly retrieves worklogs from Jira API.
- 📊 **Smart Excel Export**: Generates beautifully formatted reports with daily summaries.
- 👤 **Multi-User Support**: Easily handle multiple actors via simple environment aliases.
- 📦 **Cross-Platform**: Native support for macOS, Linux, and Windows.
- ⚡ **High Performance**: Built with Go for speed and minimal footprint.

---

## 📋 Prerequisites

- **Jira API Token**: Required for authentication. [Generate one here](https://id.atlassian.com/manage-profile/security/api-tokens).
- **Go 1.24+**: Only required if building from source.

---

## 🚀 Installation

Choose the method that best fits your environment:

### Method 1: Using Global Binaries (No Go Required)
Ideal for users who want to get started quickly without installing development tools.

1. **Download** the binary for your platform from the [`bin/`](./bin) directory.
2. **Setup**:
   - **macOS / Linux**:
     ```bash
     mv jtg-darwin-arm64 jtg # Use your platform suffix
     chmod +x jtg
     ./jtg setup-config
     ./jtg install
     ```
   - **Windows**:
     ```powershell
     rename jtg-windows-amd64.exe jtg.exe
     .\jtg.exe setup-config
     .\jtg.exe install
     ```

### Method 2: Building from Source
For developers or users with the Go toolchain installed.

- **Option A (With Make)**:
  ```bash
  make build
  ```
- **Option B (Standard Go)**:
  ```bash
  go build -o jtg cmd/jtg/main.go
  ./jtg install
  ```

> [!NOTE]
> After using the `install` command, remember to **restart your terminal** or run `source ~/.zshrc` (macOS/Linux) to activate the `jtg` command globally.

---

## ⚙️ Configuration

Set up your credentials once and run from anywhere.

### 1. Initialize Configuration
```bash
# Automates directory (~/.jtg) and .env creation
jtg setup-config 
```

### 2. Configure Credentials
Edit your newly created `.env` file:
- **macOS/Linux**: `~/.jtg/.env`
- **Windows**: `%AppData%\jtg\.env`

```env
JIRA_BASE_URL=https://your-domain.atlassian.net
JIRA_EMAIL=your@email.com
JIRA_API_TOKEN=your_token_here

# User Mapping (Add your team's account IDs)
USER_ID_HANIF=6365...
```

---

## 💡 Usage

Generate a report by specifying your alias and the month/year.

```bash
jtg -a hanif -p 01-2026
```

> [!TIP]
> Reports are automatically saved to your **Downloads** folder for instant access and sharing.

---

<p align="center">Made with ❤️ for faster workflows.</p>