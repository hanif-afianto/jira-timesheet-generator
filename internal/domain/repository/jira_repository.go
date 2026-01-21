package repository

import (
	"context"

	"github.com/hanif-afianto/jira-timesheet-generator/internal/domain/entity"
)

type JiraRepository interface {
	FetchIssues(ctx context.Context, jql string) ([]entity.Issue, error)
	FetchWorklogs(ctx context.Context, issueKey string) ([]entity.Worklog, error)
}
