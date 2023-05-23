package report

import (
	"goreporter/toggl"
	"goreporter/utils"
	"sort"
	"strconv"
	"time"
)

type Report struct {
	At             time.Time         `json:"at"`
	TotalDuration  time.Duration     `json:"totalDuration"`
	Projects       map[int]Project   `json:"projects"`
	WorkspaceId    int               `json:"workspaceId"`
	RawTimeEntries []toggl.TimeEntry `json:"rawTimeEntries"`
}

type Project struct {
	Name          string        `json:"name"`
	NonPaid       TasksBlock    `json:"nonPaid"`
	Paid          TasksBlock    `json:"paid"`
	TotalDuration time.Duration `json:"totalDuration"`
}

type TaskEntry struct {
	At       time.Time     `json:"at"`
	Duration time.Duration `json:"duration"`
	Text     string        `json:"text"`
}

type TasksBlock struct {
	Duration time.Duration `json:"duration"`
	Tasks    []TaskEntry   `json:"tasks"`
	tasksMap map[string]TaskEntry
}

type Reporter struct {
	TogglClient          toggl.Toggl
	ProjectTimePrecision time.Duration
	TaskTimePrecision    time.Duration
}

func newProject(name string) Project {
	return Project{
		Name: name,
		NonPaid: TasksBlock{
			Duration: 0,
			tasksMap: make(map[string]TaskEntry),
			Tasks:    make([]TaskEntry, 0),
		},
		Paid: TasksBlock{
			Duration: 0,
			tasksMap: make(map[string]TaskEntry),
			Tasks:    make([]TaskEntry, 0),
		},
	}
}

func newReport(at time.Time, workspaceId int) Report {
	return Report{
		At:            at,
		TotalDuration: 0,
		Projects:      make(map[int]Project),
		WorkspaceId:   workspaceId,
	}
}

func (reporter *Reporter) BuildDailyReport(workspaceId int, startDate time.Time) (Report, error) {

	// Only place where endDate used - to get details from API; here it not used, since needs to be the same as startDate
	// Change logic here to be able to grab weekly reports as well
	//endDate := startDate.Add(time.Duration(24) * time.Hour)

	endDate := startDate
	return reporter.BuildReport(workspaceId, startDate, endDate)
}

func (reporter *Reporter) BuildReport(workspaceId int, startDate time.Time, endDate time.Time) (Report, error) {
	report := newReport(startDate, workspaceId)

	if err := reporter.groupByProjectTag(&report, startDate, endDate); err != nil {
		return report, err
	}
	if err := reporter.injectProjectNames(&report); err != nil {
		return report, err
	}
	flattenTaskEntries(&report)
	sumTime(&report, reporter.ProjectTimePrecision)

	return report, nil
}

func (reporter *Reporter) groupByProjectTag(report *Report, startDate time.Time, endDate time.Time) error {
	timeEntries, err := reporter.TogglClient.GetTimeEntriesForWorkspaceV2(startDate, endDate, report.WorkspaceId)
	if err != nil {
		return err
	}

	for _, entry := range timeEntries {
		project, ok := report.Projects[entry.ProjectId]
		if !ok {
			project = newProject(strconv.Itoa(entry.ProjectId))
		}

		if _, isNonPaid := entry.Tags["non-paid"]; isNonPaid {
			project.NonPaid = upsertTask(project.NonPaid, entry, reporter.TaskTimePrecision)
		} else {
			project.Paid = upsertTask(project.Paid, entry, reporter.TaskTimePrecision)
		}

		report.Projects[entry.ProjectId] = project
	}

	report.RawTimeEntries = timeEntries

	return nil
}

func (reporter *Reporter) injectProjectNames(report *Report) error {
	projects, err := reporter.TogglClient.GetProjects(report.WorkspaceId)
	if err != nil {
		return err
	}

	for projectId := range report.Projects {
		if proj, ok := report.Projects[projectId]; ok {
			proj.Name = projects[projectId].Name
			report.Projects[projectId] = proj
		}
	}
	return nil
}

func upsertTask(block TasksBlock, entry toggl.TimeEntry, precision time.Duration) TasksBlock {
	duration := utils.RoundTime(entry.Duration, precision)

	if taskEntry, ok := block.tasksMap[entry.Description]; ok {
		taskEntry.Duration += duration
		taskEntry.At = entry.Start
		block.tasksMap[entry.Description] = taskEntry
	} else {
		block.tasksMap[entry.Description] = TaskEntry{
			Text:     entry.Description,
			Duration: duration,
			At:       entry.Start,
		}
	}

	return block
}

func sumTime(report *Report, precision time.Duration) {
	for projectId, project := range report.Projects {
		project.Paid.Duration = sumProjectTime(project.Paid, precision)
		report.TotalDuration += project.Paid.Duration

		project.NonPaid.Duration = sumProjectTime(project.NonPaid, precision)
		report.TotalDuration += project.NonPaid.Duration

		project.TotalDuration = project.Paid.Duration + project.NonPaid.Duration

		report.Projects[projectId] = project
	}
}

func sumProjectTime(block TasksBlock, precision time.Duration) time.Duration {
	duration := time.Duration(0)
	for _, task := range block.Tasks {
		duration += task.Duration
	}
	return utils.RoundTime(duration, precision)
}

func flattenTaskEntries(report *Report) {
	for projectId, project := range report.Projects {
		for _, entry := range project.Paid.tasksMap {
			project.Paid.Tasks = append(project.Paid.Tasks, entry)
		}
		sort.Slice(project.Paid.Tasks, func(i, j int) bool {
			return project.Paid.Tasks[i].At.Before(project.Paid.Tasks[j].At)
		})
		project.Paid.tasksMap = nil

		for _, entry := range project.NonPaid.tasksMap {
			project.NonPaid.Tasks = append(project.NonPaid.Tasks, entry)
		}
		sort.Slice(project.NonPaid.Tasks, func(i, j int) bool {
			return project.NonPaid.Tasks[i].At.Before(project.NonPaid.Tasks[j].At)
		})
		project.NonPaid.tasksMap = nil

		report.Projects[projectId] = project
	}
}
