# 🛠️ JTG (Jira Timesheet Generator)

**JTG** is a high-performance CLI tool designed to simplify timesheet management. It fetches worklogs directly from the Jira API and transforms them into beautifully formatted Excel reports, tailored for professional submissions.

---

## ✨ Features

- **Automated Data Retrieval**: Directly integrates with Jira API to fetch accurate worklogs.
- **Smart Formatting**: Generates structured Excel files with daily summaries and professional styling.
- **Multi-Actor Support**: Configure multiple users via environment variables.
- **Seamless Installation**: Built-in cross-platform `install` command for instant `PATH` integration.
- **Clean Architecture**: Built with Go for reliability and maintainability.

---

## 📋 Prerequisites

- **Go**: 1.24.0 or higher.
- **Jira API Token**: Required for authentication.

---

## 🚀 Installation

### 1. Build & Install
Execute the setup command to compile the binary and register it to your system `PATH`:
```bash
make build
```

### 2. Activate Changes
Restart your terminal or reload your shell profile to finalize the installation:
- **macOS/Linux**: `source ~/.zshrc` or `source ~/.bashrc`
- **Windows**: Restart your terminal (Command Prompt or PowerShell).

---

## ⚙️ Configuration

1. Initialize your environment file:
   ```bash
   cp .env.example .env
   ```
2. Configure your credentials in `.env`:
   - `JIRA_BASE_URL`: Your Jira instance URL (e.g., `https://your-domain.atlassian.net`)
   - `JIRA_EMAIL`: Your Atlassian account email.
   - `JIRA_API_TOKEN`: [Generate here](https://id.atlassian.com/manage-profile/security/api-tokens).
   - `USER_ID_<ALIAS>`: The Jira Account ID for the target user.

---

## 💡 Usage

Generate a professional timesheet by specifying the actor alias and period (MM-YYYY):

```bash
jtg -a hanif -p 01-2026
```

> [!TIP]
> Your generated reports are automatically saved to your system's `Downloads` folder for easy access.
 