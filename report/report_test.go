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
	reporter := NewReporter(&toggl.TogglData{
		Handler: &FakeHandler{},
	})

	date := time.Date(2022, 8, 4, 0, 0, 0, 0, time.UTC)
	report, err := reporter.BuildDailyReport(2200006, date)
	if err != nil {
		t.Fatalf("BuildDailyReport err: %v", err)
	}

	if report.At != date {
		t.Errorf("Expected .At to be %v, got %v", date, report.At)
	}

	totalDuration := time.Duration(7)*time.Hour + time.Duration(35)*time.Minute
	if report.TotalDuration != totalDuration {
		t.Errorf("Expected .TotalDuration to be %v, got %v", totalDuration, report.TotalDuration)
	}

	expected := map[int]time.Duration{
		158083225: time.Duration(35) * time.Minute,
		161955744: time.Duration(15) * time.Minute,
		162347134: time.Duration(30) * time.Minute,
		176616296: time.Duration(6)*time.Hour + time.Duration(15)*time.Minute,
	}

	for projectId, duration := range expected {
		got := report.Projects[projectId].TotalDuration
		if got != duration {
			t.Errorf("Expected [%v].TotalDuration to be %v, got %v", projectId, duration, got)
		}
	}

	gotPaid := report.Projects[176616296].Paid.Duration
	expectedPaid := time.Duration(5)*time.Hour + time.Duration(15)*time.Minute
	if gotPaid != expectedPaid {
		t.Errorf("Expected [%v].Paid.Duration to be %v, got %v", 176616296, expectedPaid, gotPaid)
	}

	expectedNonPaid := time.Duration(1) * time.Hour
	gotNonPaid := report.Projects[176616296].NonPaid.Duration
	if gotNonPaid != expectedNonPaid {
		t.Errorf("Expected [%v].NonPaid.Duration to be %v, got %v", 176616296, expectedNonPaid, gotNonPaid)
	}

	fmt.Println(report)
}

type FakeHandler struct {
}

func (h *FakeHandler) Execute(method string, query url.Values, response interface{}) error {
	if method == "/api/v9/me/time_entries" {
		return json.NewDecoder(strings.NewReader(jsonTimeEntries)).Decode(&response)
	}

	if method == "/api/v9/workspaces/2200006/projects" {
		return json.NewDecoder(strings.NewReader(jsonProjects)).Decode(&response)
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

const jsonTimeEntries = `[
    {
        "id": 2596515000,
        "workspace_id": 2200006,
        "project_id": 176616296,
        "duration": 1731,
        "description": "description",
        "tags": [
            "non-paid"
        ]
    },
    {
        "id": 2596391465,
        "workspace_id": 2200006,
        "project_id": 176616296,
        "duration": 5640,
        "description": "description",
        "tags": null
    },
    {
        "id": 2596357707,
        "workspace_id": 2200006,
        "project_id": 161955744,
        "duration": 761,
        "description": "description",
        "tags": null
    },
    {
        "id": 2596323642,
        "workspace_id": 2200006,
        "project_id": 162347134,
        "duration": 1517,
        "description": "description",
        "tags": null
    },
    {
        "id": 2596166690,
        "workspace_id": 2200006,
        "project_id": 176616296,
        "duration": 6787,
        "description": "description",
        "tags": null
    },
    {
        "id": 2596073008,
        "workspace_id": 527002,
        "project_id": 163370633,
        "duration": 4016,
        "description": "description",
        "tags": null
    },
    {
        "id": 2596038940,
        "workspace_id": 2200006,
        "project_id": 176616296,
        "duration": 1694,
        "description": "description",
        "tags": null
    },
    {
        "id": 2596029530,
        "workspace_id": 2200006,
        "project_id": 158083225,
        "duration": 1112,
        "description": "description",
        "tags": null
    },
    {
        "id": 2595993449,
        "workspace_id": 2200006,
        "project_id": 176616296,
        "duration": 1685,
        "description": "description",
        "tags": null
    },
    {
        "id": 2595980483,
        "workspace_id": 527002,
        "project_id": 163370633,
        "duration": 1023,
        "description": "description",
        "tags": null
    },
    {
        "id": 2595956872,
        "workspace_id": 2200006,
        "project_id": 176616296,
        "duration": 2869,
        "description": "description",
        "tags": null
    },
    {
        "id": 2595943660,
        "workspace_id": 527002,
        "project_id": 163370633,
        "duration": 1304,
        "description": "description",
        "tags": null
    },
    {
        "id": 2595928047,
        "workspace_id": 527002,
        "project_id": 163370633,
        "duration": 1153,
        "description": "description",
        "tags": null
    },
    {
        "id": 2595909078,
        "workspace_id": 2200006,
        "project_id": 176616296,
        "duration": 1782,
        "description": "description",
        "tags": [
            "non-paid"
        ]
    },
    {
        "id": 2595927555,
        "workspace_id": 2200006,
        "project_id": 158083225,
        "duration": 836,
        "description": "description",
        "tags": null
    }
]`
