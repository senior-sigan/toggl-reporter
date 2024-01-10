package forms

import (
	"goreporter/report"
	"goreporter/utils"
	"net/url"
	"strconv"
)

type GoogleFormFieldsMapping struct {
	ProjectName           string `toml:"project_name"`
	ProjectTasks          string `toml:"project_tasks"`
	ProjectHours          string `toml:"project_hours"`
	ProjectMinutes        string `toml:"project_minutes"`
	ProjectSeconds        string `toml:"project_seconds"`
	NonpaidProjectTasks   string `toml:"nonpaid_project_tasks"`
	NonpaidProjectHours   string `toml:"nonpaid_project_hours"`
	NonpaidProjectMinutes string `toml:"nonpaid_project_minutes"`
	NonpaidProjectSeconds string `toml:"nonpaid_project_seconds"`
	ReportYear            string `toml:"report_year"`
	ReportMonth           string `toml:"report_month"`
	ReportDay             string `toml:"report_day"`
}

type GoogleFormData struct {
	ProjectName           string
	ProjectTasks          string
	ProjectHours          int
	ProjectMinutes        int
	ProjectSeconds        int
	NonpaidProjectTasks   string
	NonpaidProjectHours   int
	NonpaidProjectMinutes int
	NonpaidProjectSeconds int
	ReportYear            int
	ReportMonth           int
	ReportDay             int
}

type GoogleFormGenerator struct {
	FormURL             string
	Mapping             GoogleFormFieldsMapping
	InternalProjectName string
	Formatter           *FormFormatter
}

func (gen *GoogleFormGenerator) ConvertReportToFormsData(report report.Report) map[string]GoogleFormData {
	formData := make(map[string]GoogleFormData)

	for _, project := range report.Projects {
		formData[project.Name] = GoogleFormData{
			ProjectName:           project.Name,
			ProjectTasks:          gen.Formatter.Format(project.Paid.Tasks),
			ProjectHours:          utils.Hours(project.Paid.Duration),
			ProjectMinutes:        utils.Minutes(project.Paid.Duration),
			ProjectSeconds:        utils.Seconds(project.Paid.Duration),
			NonpaidProjectTasks:   gen.Formatter.Format(project.NonPaid.Tasks),
			NonpaidProjectHours:   utils.Hours(project.NonPaid.Duration),
			NonpaidProjectMinutes: utils.Minutes(project.NonPaid.Duration),
			NonpaidProjectSeconds: utils.Seconds(project.NonPaid.Duration),
			ReportYear:            report.At.Year(),
			ReportMonth:           int(report.At.Month()),
			ReportDay:             report.At.Day(),
		}
	}

	return formData
}

func (gen *GoogleFormGenerator) ConvertReportToForms(report report.Report) map[string]string {
	u, err := url.Parse(gen.FormURL)
	if err != nil {
		panic(err)
	}

	forms := gen.ConvertReportToFormsData(report)

	formUrls := make(map[string]string)
	for _, formData := range forms {
		u.RawQuery = gen.encode(formData)
		formUrls[formData.ProjectName] = u.String()
	}

	return formUrls
}

func (gen *GoogleFormGenerator) encode(form GoogleFormData) string {
	query := url.Values{}
	query.Set(gen.Mapping.ProjectName, form.ProjectName)

	if form.ProjectTasks != "" {
		query.Set(gen.Mapping.ProjectTasks, form.ProjectTasks)
	} else {
		query.Set(gen.Mapping.ProjectTasks, "-")
	}
	query.Set(gen.Mapping.ProjectHours, strconv.Itoa(form.ProjectHours))
	query.Set(gen.Mapping.ProjectMinutes, strconv.Itoa(form.ProjectMinutes))
	query.Set(gen.Mapping.ProjectSeconds, strconv.Itoa(form.ProjectSeconds))

	query.Set(gen.Mapping.NonpaidProjectTasks, form.NonpaidProjectTasks)
	query.Set(gen.Mapping.NonpaidProjectHours, strconv.Itoa(form.NonpaidProjectHours))
	query.Set(gen.Mapping.NonpaidProjectMinutes, strconv.Itoa(form.NonpaidProjectMinutes))
	query.Set(gen.Mapping.NonpaidProjectSeconds, strconv.Itoa(form.NonpaidProjectSeconds))

	query.Set(gen.Mapping.ReportYear, strconv.Itoa(form.ReportYear))
	query.Set(gen.Mapping.ReportMonth, strconv.Itoa(form.ReportMonth))
	query.Set(gen.Mapping.ReportDay, strconv.Itoa(form.ReportDay))

	return query.Encode()
}
