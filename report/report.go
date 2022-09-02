package report

import (
	"goreporter/toggl"
	"goreporter/utils"
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
	Name          string
	NonPaid       TasksBlock
	Paid          TasksBlock
	TotalDuration time.Duration
}

type TasksBlock struct {
	Duration time.Duration
	Tasks    map[string]time.Duration
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
			Tasks:    make(map[string]time.Duration),
		},
		Paid: TasksBlock{
			Duration: 0,
			Tasks:    make(map[string]time.Duration),
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
	endDate := startDate.Add(time.Duration(24) * time.Hour)
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
	report = sumTime(report, reporter.ProjectTimePrecision)

	return report, nil
}

func (reporter *Reporter) groupByProjectTag(report *Report, startDate time.Time, endDate time.Time) error {
	timeEntries, err := reporter.TogglClient.GetTimeEntriesForWorkspaceV2(startDate, endDate, report.WorkspaceId)
	if err != nil {
		return err
	}

	for _, entry := range timeEntries {
		if _, ok := report.Projects[entry.ProjectId]; !ok {
			report.Projects[entry.ProjectId] = newProject(strconv.Itoa(entry.ProjectId))
		}

		if _, isNonPaid := entry.Tags["non-paid"]; isNonPaid {
			addTime(report.Projects[entry.ProjectId].NonPaid, entry, reporter.TaskTimePrecision)
		} else {
			addTime(report.Projects[entry.ProjectId].Paid, entry, reporter.TaskTimePrecision)
		}
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

func addTime(block TasksBlock, entry toggl.TimeEntry, precision time.Duration) TasksBlock {
	block.Tasks[entry.Description] += utils.RoundTime(entry.Duration, precision)
	return block
}

func sumTime(report Report, precision time.Duration) Report {
	for projectId := range report.Projects {
		if project, ok := report.Projects[projectId]; ok {
			project.Paid.Duration = sumProjectTime(project.Paid, precision)
			report.TotalDuration += project.Paid.Duration

			project.NonPaid.Duration = sumProjectTime(project.NonPaid, precision)
			report.TotalDuration += project.NonPaid.Duration

			project.TotalDuration = project.Paid.Duration + project.NonPaid.Duration

			report.Projects[projectId] = project
		}
	}

	return report
}

func sumProjectTime(block TasksBlock, precision time.Duration) time.Duration {
	duration := time.Duration(0)
	for _, taskDuration := range block.Tasks {
		duration += taskDuration
	}
	return utils.RoundTime(duration, precision)
}
