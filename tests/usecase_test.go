package tests

import (
	"context"
	"testing"
	"time"

	"github.com/hanif-afianto/jira-timesheet-generator/internal/domain/entity"
	"github.com/hanif-afianto/jira-timesheet-generator/internal/usecase/timesheet"
)

// MockJiraRepo implements repository.JiraRepository
type MockJiraRepo struct {
	Issues []entity.Issue
}

func (m *MockJiraRepo) FetchIssues(ctx context.Context, jql string) ([]entity.Issue, error) {
	return m.Issues, nil
}

func (m *MockJiraRepo) FetchWorklogs(ctx context.Context, issueKey string) ([]entity.Worklog, error) {
	return nil, nil
}

func TestGenerateTimesheetUsecase_Execute(t *testing.T) {
	mockRepo := &MockJiraRepo{
		Issues: []entity.Issue{
			{
				Key: "GQA-1",
				Worklogs: []entity.Worklog{
					{
						AuthorAccountID: "user123",
						Started:         time.Date(2025, 8, 20, 10, 0, 0, 0, time.UTC),
						Comment:         "Coding task",
					},
				},
			},
		},
	}

	usecase := timesheet.NewGenerateTimesheetUsecase(mockRepo)
	ctx := context.Background()

	// Test case: Valid period
	// Period 09-2025 means 16th Aug to 15th Sep
	t.Run("Valid Period 09-2025", func(t *testing.T) {
		ts, err := usecase.Execute(ctx, "hanif", "user123", "09-2025")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if ts.Actor != "hanif" {
			t.Errorf("Expected actor hanif, got %s", ts.Actor)
		}

		// Check dates
		expectedStart := "2025-08-16"
		expectedEnd := "2025-09-15"
		if ts.StartDate.Format("2006-01-02") != expectedStart {
			t.Errorf("Expected start %s, got %s", expectedStart, ts.StartDate.Format("2006-01-02"))
		}
		if ts.EndDate.Format("2006-01-02") != expectedEnd {
			t.Errorf("Expected end %s, got %s", expectedEnd, ts.EndDate.Format("2006-01-02"))
		}

		// Check if worklog was captured
		found := false
		for _, row := range ts.Rows {
			if row.Date.Format("2006-01-02") == "2025-08-20" {
				if len(row.Tasks) > 0 && row.Tasks[0] == "> GQA-1 - Coding task" {
					found = true
				}
			}
		}
		if !found {
			t.Error("Worklog for 2025-08-20 not found in timesheet")
		}
	})

	t.Run("Invalid Period", func(t *testing.T) {
		_, err := usecase.Execute(ctx, "hanif", "user123", "invalid")
		if err == nil {
			t.Error("Expected error for invalid period, got nil")
		}
	})
}
