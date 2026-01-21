package main

import (
	"log"
	"os"

	"github.com/hanif-afianto/jira-timesheet-generator/internal/infrastructure/config"
	"github.com/hanif-afianto/jira-timesheet-generator/internal/infrastructure/excel"
	"github.com/hanif-afianto/jira-timesheet-generator/internal/infrastructure/jira"
	"github.com/hanif-afianto/jira-timesheet-generator/internal/infrastructure/system"
	"github.com/hanif-afianto/jira-timesheet-generator/internal/interface/cli"
	"github.com/hanif-afianto/jira-timesheet-generator/internal/usecase/timesheet"
)

func main() {
	// 1. Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Infrastructure
	jiraClient := jira.NewJiraClient(cfg)
	excelExporter := excel.NewExcelExporter()
	pathManager := system.NewPathManager()

	// 3. Initialize Usecase
	genUsecase := timesheet.NewGenerateTimesheetUsecase(jiraClient)

	// 4. Initialize Interface (CLI)
	cliHandler := cli.NewCLIHandler(genUsecase, excelExporter, pathManager)

	// 5. Run CLI
	cliHandler.Handle(os.Args[1:])
}
