package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JiraBaseURL  string
	JiraEmail    string
	JiraAPIToken string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	return &Config{
		JiraBaseURL:  os.Getenv("JIRA_BASE_URL"),
		JiraEmail:    os.Getenv("JIRA_EMAIL"),
		JiraAPIToken: os.Getenv("JIRA_API_TOKEN"),
	}, nil
}
