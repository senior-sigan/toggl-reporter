package toggl

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

var token string
var startDate time.Time
var endDate time.Time

func TestMain(m *testing.M) {
	token = os.Getenv("TOGGLE_TOKEN")
	if token == "" {
		log.Fatal("TOGGLE_TOKEN must be set")
	}
	startDate = time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC)
	endDate = time.Date(2022, 7, 2, 0, 0, 0, 0, time.UTC)

	runTests := m.Run()
	os.Exit(runTests)
}

func TestToggle_GetMe(t *testing.T) {
	toggle := NewToggl(token)
	me, err := toggle.GetMe()
	if err != nil {
		t.Errorf("Me error %v", err)
	}
	fmt.Println(me)
}

func TestToggle_GetWorkspaces(t *testing.T) {
	toggle := NewToggl(token)
	workspaces, err := toggle.GetWorkspaces()
	if err != nil {
		t.Errorf("Me error %v", err)
	}
	fmt.Println(workspaces)
}

func TestToggle_GetTimeEntries(t *testing.T) {
	toggle := NewToggl(token)
	entries, err := toggle.GetTimeEntries(startDate, endDate)
	if err != nil {
		t.Errorf("Me error %v", err)
	}
	fmt.Println(entries)
}

func TestToggle_GetTimeEntriesForWorkspace_EmptyLen(t *testing.T) {
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
	toggl := NewToggl(token)
	projects, err := toggl.GetProjects(2200006)
	if err != nil {
		t.Errorf("GetProjects error %v", err)
	}
	fmt.Println(projects)
}
