package report

import (
	"encoding/json"
	"fmt"
	"goreporter/toggl"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestReporter_BuildDailyReport(t *testing.T) {
	reporter := Reporter{
		TogglClient: &toggl.TogglData{
			Handler: &FakeHandler{},
		},
		ProjectTimePrecision: time.Duration(1) * time.Minute,
		TaskTimePrecision:    time.Duration(1) * time.Minute,
	}

	date := time.Date(2022, 8, 4, 0, 0, 0, 0, time.UTC)
	report, err := reporter.BuildDailyReport(2200006, date)
	if err != nil {
		t.Fatalf("BuildDailyReport err: %v", err)
	}

	if report.At != date {
		t.Errorf("Expected .At to be %v, got %v", date, report.At)
	}

	totalDuration := time.Duration(8) * time.Hour
	if report.TotalDuration != totalDuration {
		t.Errorf("Expected .TotalDuration to be %v, got %v", totalDuration, report.TotalDuration)
	}

	expected := map[int]time.Duration{
		158083225: time.Duration(30) * time.Minute,
		161955744: time.Duration(15) * time.Minute,
		162347134: time.Duration(30) * time.Minute,
		176616296: time.Duration(6)*time.Hour + time.Duration(45)*time.Minute,
	}

	for projectId, duration := range expected {
		got := report.Projects[projectId].TotalDuration
		if got != duration {
			t.Errorf("Expected [%v].TotalDuration to be %v, got %v", projectId, duration, got)
		}
	}

	gotPaid := report.Projects[176616296].Paid.Duration
	expectedPaid := time.Duration(5)*time.Hour + time.Duration(5)*time.Minute
	if gotPaid != expectedPaid {
		t.Errorf("Expected [%v].Paid.Duration to be %v, got %v", 176616296, expectedPaid, gotPaid)
	}

	expectedNonPaid := time.Duration(100) * time.Minute
	gotNonPaid := report.Projects[176616296].NonPaid.Duration
	if gotNonPaid != expectedNonPaid {
		t.Errorf("Expected [%v].NonPaid.Duration to be %v, got %v", 176616296, expectedNonPaid, gotNonPaid)
	}

	fmt.Println(report)
}

type FakeHandler struct {
}

func (h *FakeHandler) Execute(method string, query url.Values, response interface{}) error {
	if method == "/api/v9/workspaces/2200006/projects" {
		return json.NewDecoder(strings.NewReader(jsonProjects)).Decode(&response)
	}

	if method == "/reports/api/v2/details" {
		return json.NewDecoder(strings.NewReader(jsonDetails)).Decode(&response)
	}

	return fmt.Errorf("not implemented %v", method)
}

const jsonProjects = `[
    {
        "id": 161955744,
        "name": "Pm",
        "workspace_id": 2200006
    },
    {
        "id": 158083225,
        "name": "Other",
        "workspace_id": 2200006
    },
    {
        "id": 162347134,
        "name": "Invest",
        "workspace_id": 2200006
    },
    {
        "id": 176616296,
        "name": "Project",
        "workspace_id": 2200006
    }
]`

const jsonDetails = `{
	"total_billable": null,
	"total_count": 11,
    "per_page": 50,
	"data": [{"id":2596515000,"pid":176616296,"project":"name_176616296","start":"2022-08-04T18:25:30+06:00","end":"2022-08-04T19:00:30+06:00","dur":2100000,"description":"description","tags":["non-paid"]},{"id":2596391465,"pid":176616296,"project":"name_176616296","start":"2022-08-04T16:45:29+06:00","end":"2022-08-04T18:25:29+06:00","dur":6000000,"description":"description","tags":[]},{"id":2596357707,"pid":161955744,"project":"name_161955744","start":"2022-08-04T16:30:34+06:00","end":"2022-08-04T16:45:34+06:00","dur":900000,"description":"description","tags":[]},{"id":2596323642,"pid":162347134,"project":"name_162347134","start":"2022-08-04T16:00:14+06:00","end":"2022-08-04T16:30:14+06:00","dur":1800000,"description":"description","tags":[]},{"id":2596166690,"pid":176616296,"project":"name_176616296","start":"2022-08-04T14:00:04+06:00","end":"2022-08-04T16:00:04+06:00","dur":7200000,"description":"description","tags":[]},{"id":2596038940,"pid":176616296,"project":"name_176616296","start":"2022-08-04T12:30:46+06:00","end":"2022-08-04T13:05:46+06:00","dur":2100000,"description":"description","tags":[]},{"id":2596029530,"pid":158083225,"project":"name_158083225","start":"2022-08-04T12:15:12+06:00","end":"2022-08-04T12:30:12+06:00","dur":900000,"description":"description","tags":[]},{"id":2595993449,"pid":176616296,"project":"name_176616296","start":"2022-08-04T11:45:12+06:00","end":"2022-08-04T12:15:12+06:00","dur":1800000,"description":"description","tags":["non-paid"]},{"id":2595956872,"pid":176616296,"project":"name_176616296","start":"2022-08-04T10:40:14+06:00","end":"2022-08-04T11:30:14+06:00","dur":3000000,"description":"description","tags":[]},{"id":2595909078,"pid":176616296,"project":"name_176616296","start":"2022-08-04T09:25:57+06:00","end":"2022-08-04T10:00:57+06:00","dur":2100000,"description":"description","tags":["non-paid"]},{"id":2595927555,"pid":158083225,"project":"name_158083225","start":"2022-08-04T09:10:27+06:00","end":"2022-08-04T09:25:27+06:00","dur":900000,"description":"description","tags":[]}]
}`
