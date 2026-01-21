package timesheet

import (
	"context"
	"fmt"
	"time"

	"github.com/hanif-afianto/jira-timesheet-generator/internal/domain/entity"
	"github.com/hanif-afianto/jira-timesheet-generator/internal/domain/repository"
)

type GenerateTimesheetUsecase struct {
	jiraRepo repository.JiraRepository
}

func NewGenerateTimesheetUsecase(repo repository.JiraRepository) *GenerateTimesheetUsecase {
	return &GenerateTimesheetUsecase{jiraRepo: repo}
}

func (u *GenerateTimesheetUsecase) Execute(ctx context.Context, actor, userID, period string) (*entity.Timesheet, error) {
	startDate, endDate, err := u.getCutoffDates(period)
	if err != nil {
		return nil, err
	}

	jql := fmt.Sprintf(`project = GQA AND worklogAuthor = %s AND worklogDate >= "%s" AND worklogDate <= "%s"`,
		userID, startDate.Format("2006/01/02"), endDate.Format("2006/01/02"))

	issues, err := u.jiraRepo.FetchIssues(ctx, jql)
	if err != nil {
		return nil, err
	}

	// Group worklogs by date
	grouped := make(map[string][]string)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		grouped[d.Format("2006-01-02")] = []string{}
	}

	for _, issue := range issues {
		worklogs := issue.Worklogs

		for _, wl := range worklogs {
			if wl.AuthorAccountID != userID {
				continue
			}
			dateStr := wl.Started.Format("2006-01-02")
			if _, ok := grouped[dateStr]; ok {
				line := fmt.Sprintf("> %s - %s", issue.Key, wl.Comment)
				grouped[dateStr] = append(grouped[dateStr], line)
			}
		}
	}

	// Build rows
	var rows []entity.TimesheetRow
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		rows = append(rows, entity.TimesheetRow{
			Date:  d,
			Tasks: grouped[dateStr],
		})
	}

	return &entity.Timesheet{
		Actor:     actor,
		Period:    period,
		StartDate: startDate,
		EndDate:   endDate,
		Rows:      rows,
	}, nil
}

func (u *GenerateTimesheetUsecase) getCutoffDates(period string) (time.Time, time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	t, err := time.ParseInLocation("01-2006", period, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid period format")
	}

	// Start = 16th of previous month
	startMonth := t.AddDate(0, -1, 0)
	start := time.Date(startMonth.Year(), startMonth.Month(), 16, 0, 0, 0, 0, loc)

	// End = 15th of current month
	end := time.Date(t.Year(), t.Month(), 15, 23, 59, 59, 0, loc)

	return start, end, nil
}
