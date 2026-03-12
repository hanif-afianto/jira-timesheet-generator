package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	JiraBaseURL  string
	JiraEmail    string
	JiraAPIToken string
}

func LoadConfig() (*Config, error) {
	// Search for .env in multiple locations
	homeDir, _ := os.UserHomeDir()
	configDirs := []string{
		".", // Current directory
	}

	// Add OS-specific config directories
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			configDirs = append(configDirs, filepath.Join(appData, "jtg"))
		}
	}
	configDirs = append(configDirs, filepath.Join(homeDir, ".jtg"))

	for _, dir := range configDirs {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			_ = godotenv.Load(envPath)
			break
		}
	}

	return &Config{
		JiraBaseURL:  os.Getenv("JIRA_BASE_URL"),
		JiraEmail:    os.Getenv("JIRA_EMAIL"),
		JiraAPIToken: os.Getenv("JIRA_API_TOKEN"),
	}, nil
}
