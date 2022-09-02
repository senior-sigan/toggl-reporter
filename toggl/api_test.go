package toggl

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"
)

var token string
var startDate time.Time
var endDate time.Time

type FakeHandler struct {
	Executor func(method string, query url.Values, response interface{}) error
}

func (h *FakeHandler) Execute(method string, query url.Values, response interface{}) error {
	return h.Executor(method, query, response)
}

func TestMain(m *testing.M) {
	token = os.Getenv("TOGGLE_TOKEN")
	startDate = time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC)
	endDate = time.Date(2022, 7, 2, 0, 0, 0, 0, time.UTC)

	runTests := m.Run()
	os.Exit(runTests)
}

func TestToggle_GetMe(t *testing.T) {
	if len(token) == 0 {
		t.Skip("Token is not set")
	}
	toggle := NewToggl(token)
	me, err := toggle.GetMe()
	if err != nil {
		t.Errorf("Me error %v", err)
	}
	fmt.Println(me)
}

func TestToggle_GetWorkspaces(t *testing.T) {
	if len(token) == 0 {
		t.Skip("Token is not set")
	}
	toggle := NewToggl(token)
	workspaces, err := toggle.GetWorkspaces()
	if err != nil {
		t.Errorf("Me error %v", err)
	}
	fmt.Println(workspaces)
}

func TestToggle_GetTimeEntries(t *testing.T) {
	if len(token) == 0 {
		t.Skip("Token is not set")
	}
	toggle := NewToggl(token)
	entries, err := toggle.GetTimeEntries(startDate, endDate)
	if err != nil {
		t.Errorf("Me error %v", err)
	}
	fmt.Println(entries)
}

func TestToggle_GetTimeEntriesForWorkspace_EmptyLen(t *testing.T) {
	if len(token) == 0 {
		t.Skip("Token is not set")
	}
	toggle := NewToggl(token)
	entries, err := toggle.GetTimeEntriesForWorkspace(startDate, endDate, -1)
	if err != nil {
		t.Errorf("Me error %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("Expected 0 but got %d", len(entries))
	}
	fmt.Println(entries)
}

func TestToggle_GetTimeEntriesForWorkspace(t *testing.T) {
	if len(token) == 0 {
		t.Skip("Token is not set")
	}
	toggle := NewToggl(token)
	entries, err := toggle.GetTimeEntriesForWorkspace(startDate, endDate, 2200006)
	if err != nil {
		t.Errorf("Me error %v", err)
	}
	expected := 6
	if len(entries) != expected {
		t.Errorf("Expected %d but got %d", expected, len(entries))
	}
	fmt.Println(entries)
}

func TestToggl_GetProjects(t *testing.T) {
	if len(token) == 0 {
		t.Skip("Token is not set")
	}
	toggl := NewToggl(token)
	projects, err := toggl.GetProjects(2200006)
	if err != nil {
		t.Errorf("GetProjects error %v", err)
	}
	fmt.Println(projects)
}

func TestTogglData_GetTimeEntriesForWorkspaceV2(t *testing.T) {
	if len(token) == 0 {
		t.Skip("Token is not set")
	}
	toggl := NewToggl(token)
	startDate = time.Date(2022, 8, 5, 0, 0, 0, 0, time.UTC)
	endDate = time.Date(2022, 8, 6, 0, 0, 0, 0, time.UTC)
	entries, err := toggl.GetTimeEntriesForWorkspaceV2(startDate, endDate, 2200006)
	if err != nil {
		t.Errorf("GetTimeEntriesForWorkspaceV2 error %v", err)
	}
	fmt.Printf("%v", entries)
}
