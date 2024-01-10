package forms

import (
	"goreporter/report"
	"log"
	"testing"
	"time"
)

func GenTestMapping() GoogleFormFieldsMapping {
	return GoogleFormFieldsMapping{
		ProjectName:           "entry.2004201849",
		ProjectTasks:          "entry.60008525",
		ProjectHours:          "entry.291946000_hour",
		ProjectMinutes:        "entry.291946000_minute",
		ProjectSeconds:        "entry.291946000_second",
		NonpaidProjectTasks:   "entry.714539999",
		NonpaidProjectHours:   "entry.52836812_hour",
		NonpaidProjectMinutes: "entry.52836812_minute",
		NonpaidProjectSeconds: "entry.52836812_second",
		ReportYear:            "entry.1783131858_year",
		ReportMonth:           "entry.1783131858_month",
		ReportDay:             "entry.1783131858_day",
	}
}

func GenTestReport() report.Report {
	return report.Report{
		At:            time.Now(),
		TotalDuration: time.Duration(4) * time.Hour,
		Projects: map[int]report.Project{
			42: {
				Name:    "Project #1",
				NonPaid: report.TasksBlock{},
				Paid: report.TasksBlock{
					Duration: time.Duration(3)*time.Hour + time.Duration(15)*time.Minute,
					Tasks: []report.TaskEntry{
						{Text: "Some task", Duration: time.Duration(2) * time.Hour},
						{Text: "Another task", Duration: time.Duration(29) * time.Minute},
						{Text: "Another task 2", Duration: time.Duration(14) * time.Minute},
						{Text: "Last task", Duration: time.Duration(28) * time.Minute},
					},
				},
			},
			777: {
				Name: "Project #2",
				NonPaid: report.TasksBlock{
					Duration: time.Duration(1) * time.Hour,
					Tasks: []report.TaskEntry{
						{Text: "Nonpaid task 1", Duration: time.Duration(43) * time.Minute},
						{Text: "Nonpaid task 2", Duration: time.Duration(15) * time.Minute},
					},
				},
				Paid: report.TasksBlock{},
			},
		},
		WorkspaceId: 0,
	}
}

func TestGoogleFormGenerator_ConvertReportToFormsData(t *testing.T) {
	gen := GoogleFormGenerator{
		FormURL:             "https://docs.google.com/forms/d/e/KEY/viewform",
		Mapping:             GenTestMapping(),
		InternalProjectName: "internal",
		Formatter:           NewFormFormatter(),
	}
	testReport := GenTestReport()
	forms := gen.ConvertReportToFormsData(testReport)

	log.Println(forms["Project #1"])
	log.Println(forms["Project #2"])

	if forms["Project #1"].ProjectHours != 3 {
		t.Fatalf("Expects 3 got %v", forms["Project #1"].ProjectHours)
	}

	if forms["Project #1"].ProjectMinutes != 15 {
		t.Fatalf("Expects 15 got %v", forms["Project #1"].ProjectMinutes)
	}

	if forms["Project #2"].NonpaidProjectHours != 1 {
		t.Fatalf("Expects 1 got %v", forms["Project #2"].NonpaidProjectHours)
	}

	if forms["Project #2"].NonpaidProjectMinutes != 0 {
		t.Fatalf("Expects 0 got %v", forms["Project #2"].NonpaidProjectMinutes)
	}

}
