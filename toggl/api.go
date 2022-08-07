package toggl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	DateFormat = "2006-01-02"
	BaseURL    = "https://api.track.toggl.com"
	AppName    = "senior_sigan_reporter"
	UserAgent  = "toggl/ (+https://github.com/senior-sigan/toggl-reporter)"
)

type Params map[string]interface{}

type TogglData struct {
	Handler RequestHandler
}

type Toggl interface {
	GetMe() (response Me, err error)
	GetWorkspaces() (response []Workspace, err error)
	GetTimeEntries(startDate time.Time, endDate time.Time) (response []TimeEntry, err error)
	GetTimeEntriesForWorkspace(startDate time.Time, endDate time.Time, workspaceId int) ([]TimeEntry, error)
	GetTimeEntriesForWorkspaceV2(startDate time.Time, endDate time.Time, workspaceId int) ([]TimeEntry, error)
	GetProjects(workspaceId int) (map[int]Project, error)
}
type TimeEntry struct {
	Id          int64         `json:"id"`
	WorkspaceId int           `json:"workspace_id"`
	ProjectId   int           `json:"project_id"`
	Start       time.Time     `json:"start"`
	Stop        time.Time     `json:"stop"`
	Duration    time.Duration `json:"duration"`
	Description string        `json:"description"`
	TagsArray   []string      `json:"tags"`
	Tags        map[string]bool
}

type Me struct {
	Id                 int    `json:"id"`
	Email              string `json:"email"`
	Fullname           string `json:"fullname"`
	Timezone           string `json:"timezone"`
	DefaultWorkspaceId int    `json:"default_workspace_id"`
	ImageURL           string `json:"image_url"`
}

type Workspace struct {
	Id             int    `json:"id"`
	OrganizationId int    `json:"organization_id"`
	Name           string `json:"name"`
	Profile        int    `json:"profile"`
}

type Project struct {
	Id          int    `json:"id"`
	WorkspaceId int    `json:"workspace_id"`
	Name        string `json:"name"`
}

func NewToggl(token string) *TogglData {
	return &TogglData{
		Handler: &DefaultHandler{
			token:     token,
			UserAgent: UserAgent,
			Client:    http.DefaultClient,
			BaseURL:   BaseURL,
		},
	}
}

type DefaultHandler struct {
	token     string
	UserAgent string
	Client    *http.Client
	BaseURL   string
}

type RequestHandler interface {
	Execute(method string, query url.Values, response interface{}) error
}

func (h *DefaultHandler) Execute(method string, query url.Values, response interface{}) error {
	ctx := context.Background()
	u := h.BaseURL + method

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	req.URL.RawQuery = query.Encode()

	if err != nil {
		return fmt.Errorf("%v %v", req.URL.String(), err)
	}

	req.Header.Set("User-Agent", h.UserAgent)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(h.token, "api_token")

	resp, err := h.Client.Do(req)
	if err != nil {
		return fmt.Errorf("%v %v", req.URL.String(), err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("%v %v", req.URL.String(), resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	return nil
}

func (toggl *TogglData) GetMe() (response Me, err error) {
	err = toggl.Handler.Execute("/api/v9/me", url.Values{}, &response)
	return response, err
}

func (toggl *TogglData) GetWorkspaces() (response []Workspace, err error) {
	err = toggl.Handler.Execute("/api/v9/me/workspaces", url.Values{}, &response)
	return response, err
}

func (toggl *TogglData) GetTimeEntries(startDate time.Time, endDate time.Time) (response []TimeEntry, err error) {
	query := url.Values{}
	query.Set("start_date", startDate.Format(DateFormat))
	query.Add("end_date", endDate.Format(DateFormat))
	err = toggl.Handler.Execute("/api/v9/me/time_entries", query, &response)

	for i := range response {
		response[i].Duration *= time.Second // convert raw int into time.Duration seconds
		response[i].Tags = make(map[string]bool)

		for _, tag := range response[i].TagsArray { // convert into set
			response[i].Tags[tag] = true
		}
	}

	return
}

func (toggl *TogglData) GetTimeEntriesForWorkspace(startDate time.Time, endDate time.Time, workspaceId int) ([]TimeEntry, error) {
	// TODO: replace this method with the API request once Toggl adds it

	entries, err := toggl.GetTimeEntries(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var filteredEntries []TimeEntry
	for _, timeEntry := range entries {
		if timeEntry.WorkspaceId == workspaceId {
			filteredEntries = append(filteredEntries, timeEntry)
		}
	}

	return filteredEntries, nil
}

func (toggl *TogglData) GetTimeEntriesForWorkspaceV2(startDate time.Time, endDate time.Time, workspaceId int) ([]TimeEntry, error) {
	query := url.Values{}
	query.Set("since", startDate.Format(DateFormat))
	query.Add("until", startDate.Format(DateFormat)) // TODO: Sadly this API need start and end date to be the same to generate a daily report
	query.Add("workspace_id", strconv.Itoa(workspaceId))
	query.Add("user_agent", AppName)

	var response struct {
		Data []struct {
			Id          int64         `json:"id"`
			ProjectId   int           `json:"pid"`
			ProjectName string        `json:"project"`
			Start       time.Time     `json:"start"`
			Stop        time.Time     `json:"end"`
			Duration    time.Duration `json:"dur"`
			Description string        `json:"description"`
			TagsArray   []string      `json:"tags"`
		} `json:"data"`
	}
	if err := toggl.Handler.Execute("/reports/api/v2/details", query, &response); err != nil {
		return nil, err
	}

	entries := make([]TimeEntry, len(response.Data))
	for i, entry := range response.Data {
		tags := make(map[string]bool)

		for _, tag := range entry.TagsArray { // convert into set
			tags[tag] = true
		}

		entries[i] = TimeEntry{
			Id:          entry.Id,
			WorkspaceId: workspaceId,
			ProjectId:   entry.ProjectId,
			Start:       entry.Start,
			Stop:        entry.Stop,
			Duration:    entry.Duration * time.Millisecond,
			Description: entry.Description,
			TagsArray:   entry.TagsArray,
			Tags:        tags,
		}
	}

	return entries, nil
}

func (toggl *TogglData) GetProjects(workspaceId int) (map[int]Project, error) {
	var response []Project
	path := fmt.Sprintf("/api/v9/workspaces/%d/projects", workspaceId)
	err := toggl.Handler.Execute(path, url.Values{}, &response)
	if err != nil {
		return nil, err
	}

	projects := make(map[int]Project)
	for _, project := range response {
		projects[project.Id] = project
	}
	return projects, nil
}
