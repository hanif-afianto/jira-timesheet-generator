package system

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type PathManager struct{}

func NewPathManager() *PathManager {
	return &PathManager{}
}

func (m *PathManager) AddToPath() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	binDir := filepath.Dir(exePath)

	switch runtime.GOOS {
	case "windows":
		return m.addToPathWindows(binDir)
	case "darwin", "linux":
		return m.addToPathUnix(binDir)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func (m *PathManager) addToPathWindows(binDir string) error {
	// Use PowerShell to add to User PATH
	// First check if already in PATH
	checkCmd := exec.Command("powershell", "-Command", "[Environment]::GetEnvironmentVariable('Path', 'User')")
	output, err := checkCmd.Output()
	if err == nil && strings.Contains(string(output), binDir) {
		fmt.Println("Path already exists in Windows User PATH")
		return nil
	}

	cmd := exec.Command("powershell", "-Command", 
		fmt.Sprintf("[Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path', 'User') + ';%s', 'User')", binDir))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update Windows PATH: %w", err)
	}
	fmt.Println("Successfully added to Windows User PATH. Please restart your shell.")
	return nil
}

func (m *PathManager) addToPathUnix(binDir string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	shell := os.Getenv("SHELL")
	var configFile string

	if strings.Contains(shell, "zsh") {
		configFile = filepath.Join(homeDir, ".zshrc")
	} else {
		configFile = filepath.Join(homeDir, ".bashrc")
	}

	// Check if already in PATH in config file
	content, _ := os.ReadFile(configFile)
	if strings.Contains(string(content), binDir) {
		fmt.Printf("Path already exists in %s\n", configFile)
		return nil
	}

	exportLine := fmt.Sprintf("\nexport PATH=\"$PATH:%s\"\n", binDir)
	f, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file %s: %w", configFile, err)
	}
	defer f.Close()

	if _, err := f.WriteString(exportLine); err != nil {
		return fmt.Errorf("failed to write to config file: %w", err)
	}
	fmt.Printf("Added to %s. Please restart your terminal or run: source %s\n", configFile, configFile)

	return nil
}
const defaultEnvContent = `JIRA_BASE_URL=https://your-domain.atlassian.net
JIRA_EMAIL=your-email@domain.com
JIRA_API_TOKEN=your-api-token

# User Mapping (Add your team's account IDs)
USER_ID_HANIF=your-jira-account-id
`

func (m *PathManager) SetupConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	var configDir string
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			configDir = filepath.Join(appData, "jtg")
		} else {
			configDir = filepath.Join(homeDir, ".jtg")
		}
	} else {
		configDir = filepath.Join(homeDir, ".jtg")
	}

	// 1. Create directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	targetEnv := filepath.Join(configDir, ".env")
	if _, err := os.Stat(targetEnv); err == nil {
		fmt.Printf("Config file already exists at %s\n", targetEnv)
		return nil
	}

	// 2. Initialize .env using embedded template
	if err := os.WriteFile(targetEnv, []byte(defaultEnvContent), 0644); err != nil {
		return fmt.Errorf("failed to initialize %s: %w", targetEnv, err)
	}

	fmt.Printf("Successfully created configuration at %s\n", targetEnv)
	fmt.Println("Please edit this file with your Jira credentials.")

	return nil
}
