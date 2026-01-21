package excel

import (
	"fmt"
	"strings"
	"time"

	"github.com/hanif-afianto/jira-timesheet-generator/internal/domain/entity"

	"github.com/xuri/excelize/v2"
)

type ExcelExporter struct{}

func NewExcelExporter() *ExcelExporter {
	return &ExcelExporter{}
}

func (e *ExcelExporter) Export(timesheet *entity.Timesheet, outputPath string) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Timesheet"
	f.SetSheetName("Sheet1", sheet)

	// Headers
	f.SetCellValue(sheet, "A1", "Date")
	f.SetCellValue(sheet, "B1", "Tasks")

	row := 2
	for _, tr := range timesheet.Rows {
		dateStr := tr.Date.Format("2006-01-02")
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), dateStr)
		
		taskStr := strings.Join(tr.Tasks, "\n")
		if len(tr.Tasks) == 0 {
			taskStr = e.getFallbackValue(tr.Date)
		}
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), taskStr)

		// Styling
		style, _ := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{WrapText: true, Vertical: "top"},
		})
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), style)
		
		row++
	}

	f.SetColWidth(sheet, "A", "A", 12)
	f.SetColWidth(sheet, "B", "B", 80)

	return f.SaveAs(outputPath)
}

func (e *ExcelExporter) getFallbackValue(d time.Time) string {
	switch d.Weekday() {
	case time.Saturday:
		return "Sabtu"
	case time.Sunday:
		return "Minggu"
	default:
		return "-"
	}
}
