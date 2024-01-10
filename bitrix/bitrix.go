package bitrix

import (
	"fmt"
	"goreporter/report"
)

type BitrixGenerator struct {
	URL      string
	Projects map[string]string
}

func (gen *BitrixGenerator) buildUrl(project report.Project) string {
	bitrixID, ok := gen.Projects[project.Name]
	if ok {
		return fmt.Sprintf("%sworkgroups/group/%s/tasks/", gen.URL, bitrixID)
	}
	return ""
}

func (gen *BitrixGenerator) BuildForms(report report.Report) map[int]string {
	urls := make(map[int]string)
	for projectID := range report.Projects {
		project := report.Projects[projectID]
		urls[projectID] = gen.buildUrl(project)
	}
	return urls
}
