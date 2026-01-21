package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hanif-afianto/jira-timesheet-generator/internal/infrastructure/excel"
	"github.com/hanif-afianto/jira-timesheet-generator/internal/infrastructure/system"
	"github.com/hanif-afianto/jira-timesheet-generator/internal/usecase/timesheet"
)

type CLIHandler struct {
	usecase       *timesheet.GenerateTimesheetUsecase
	excelExporter *excel.ExcelExporter
	pathManager   *system.PathManager
}

func NewCLIHandler(u *timesheet.GenerateTimesheetUsecase, e *excel.ExcelExporter, p *system.PathManager) *CLIHandler {
	return &CLIHandler{usecase: u, excelExporter: e, pathManager: p}
}

func (h *CLIHandler) Handle(args []string) {
	if len(args) > 0 && args[0] == "install" {
		if err := h.pathManager.AddToPath(); err != nil {
			fmt.Printf("Error during installation: %v\n", err)
		}
		return
	}

	fs := flag.NewFlagSet("jtg", flag.ExitOnError)
	actor := fs.String("a", "", "Actor name")
	period := fs.String("p", "", "Period in format MM-YYYY")
	fs.Parse(args)

	if *actor == "" || *period == "" {
		fmt.Println("Usage:")
		fmt.Println("  jtg -a <actor> -p <period>  (Generate timesheet)")
		fmt.Println("  jtg install                 (Add jtg to system PATH)")
		return
	}

	envKey := fmt.Sprintf("USER_ID_%s", strings.ToUpper(*actor))
	userID := os.Getenv(envKey)
	if userID == "" {
		fmt.Printf("Error: missing %s in environment variables\n", envKey)
		return
	}

	fmt.Printf("Generating timesheet for %s, period %s...\n", *actor, *period)

	ctx := context.Background()
	ts, err := h.usecase.Execute(ctx, *actor, userID, *period)
	if err != nil {
		fmt.Printf("Error generating timesheet: %v\n", err)
		return
	}

	filename := fmt.Sprintf("timesheet_%s_%s.xlsx", *actor, *period)
	downloadDir := os.Getenv("HOME") + "/Downloads"
	outputPath := fmt.Sprintf("%s/%s", downloadDir, filename)

	if err := h.excelExporter.Export(ts, outputPath); err != nil {
		fmt.Printf("Error exporting to Excel: %v\n", err)
		return
	}

	fmt.Printf("✅ Timesheet exported successfully: %s\n", outputPath)
}
