package redmine

import (
	"fmt"
	"goreporter/report"
	"goreporter/utils"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type ReportGenerator struct {
	UrlMap map[string]string
}

// FindTaskId Finds task ID in the string "Task #42: description" or "42: description".
func FindTaskId(description string) (int, bool) {
	re1 := regexp.MustCompile(`^(\d+):\s+`)   // 42: description
	re2 := regexp.MustCompile(`\s+#(\d+)\s+`) // Task #42 description

	groups1 := re1.FindStringSubmatch(description)
	if len(groups1) >= 2 {
		if taskIdx, err := strconv.Atoi(groups1[1]); err == nil {
			return taskIdx, true
		}
	}

	groups2 := re2.FindStringSubmatch(description)
	if len(groups2) >= 2 {
		if taskIdx, err := strconv.Atoi(groups2[1]); err == nil {
			return taskIdx, true
		}
	}

	return 0, false
}

type TasksBlock struct {
	Tasks map[string]string
}

type RedmineForma struct {
	Date    string
	Hours   string
	Comment string
}

func DurationFormat(duration time.Duration) string {
	return fmt.Sprintf("%02d:%02d", utils.Hours(duration), utils.Minutes(duration))
}

func DateFormat(date time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", date.Year(), int(date.Month()), date.Day())
}

func BuildRedmineUrl(baseUrl string, id int, comment string, duration time.Duration, at time.Time) string {
	query := url.Values{}
	query.Set("time_entry[hours]", DurationFormat(duration))
	query.Set("time_entry[spent_on]", DateFormat(at))
	query.Set("time_entry[comments]", comment)

	rawUrl := fmt.Sprintf("%s/issues/%d/time_entries/new", baseUrl, id)
	u, err := url.Parse(rawUrl)
	if err != nil {
		panic(err)
	}
	u.RawQuery = query.Encode()
	return u.String()
}

func (form *ReportGenerator) GetUrl(project string) string {
	ret := form.UrlMap[project]
	if ret == "" {
		ret = form.UrlMap["default"]
	}
	return ret
}

func (form *ReportGenerator) BuildRedmineReportForms(report report.Report) map[int]map[string]string {
	rreport := make(map[int]map[string]string)
	for projectID, project := range report.Projects {
		tasks := make(map[string]string)
		for _, task := range project.Paid.Tasks {
			if idx, ok := FindTaskId(task.Text); ok {
				tasks[task.Text] = BuildRedmineUrl(form.GetUrl(project.Name), idx, task.Text, task.Duration, report.At)
			}
		}
		rreport[projectID] = tasks
	}
	return rreport
}
