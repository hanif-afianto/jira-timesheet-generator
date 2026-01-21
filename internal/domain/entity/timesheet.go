package entity

import "time"

type Worklog struct {
	AuthorAccountID string
	TimeSpent       string
	Comment         string
	Started         time.Time
}

type Issue struct {
	Key      string
	Summary  string
	Worklogs []Worklog
}

type TimesheetRow struct {
	Date  time.Time
	Tasks []string
}

type Timesheet struct {
	Actor     string
	Period    string
	StartDate time.Time
	EndDate   time.Time
	Rows      []TimesheetRow
}
