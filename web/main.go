package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"goreporter/achievements"
	"goreporter/db"
	"goreporter/forms"
	"goreporter/redmine"
	"goreporter/report"
	"goreporter/toggl"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static/*
var assetsFS embed.FS

//go:embed templates/*
var templatesFS embed.FS

var renderer *Renderer
var formGenerator forms.GoogleFormGenerator
var redmineGenerator redmine.ReportGenerator
var config Config
var storage *db.DataBase

const (
	CookieToken     = "togglToken"
	CookieWorkspace = "togglWorkspaceID"
	CookieMaxAge    = 2592000
)

type ContextKey int

const userContextKey ContextKey = 1

type User struct {
	Toggl       toggl.Toggl
	WorkspaceId int
}

type WorkspacesPage struct {
	User       toggl.Me
	Workspaces []toggl.Workspace
}

type AchievementsPage struct {
	ReportJSON      string
	At              time.Time
	User            toggl.Me
	AchievementsMap map[string]achievements.UserAchievement
}

type ReportPage struct {
	Report          report.Report
	ReportJSON      string
	FormData        map[string]string
	RedmineData     map[int]map[string]string
	AchievementsMap map[string]achievements.UserAchievement
}

func main() {
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		panic(err)
	}
	storage = db.NewDatabase()

	renderer = NewRenderer(&templatesFS)
	renderer.Register("login", "templates/login.tmpl")
	renderer.Register("index", "templates/index.tmpl")
	renderer.Register("workspaces", "templates/workspaces.tmpl")

	formGenerator = forms.GoogleFormGenerator{
		FormURL:             config.Forms.Google.Params.Url,
		Mapping:             config.Forms.Google.Mapping,
		InternalProjectName: config.Forms.Google.Params.Internal,
		Formatter:           forms.NewFormFormatter(),
	}

	redmineGenerator = redmine.ReportGenerator{
		UrlMap: config.Forms.Redmine.Urls,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	fs := http.FileServer(http.FS(assetsFS))
	r.Handle("/static/*", fs)

	r.Get("/login", ShowLogin)
	r.Post("/login", LoginUser)

	r.Mount("/workspace", func() http.Handler {
		r := chi.NewRouter()
		r.Use(UserOnly)
		r.Get("/", ShowWorkspaces)
		r.Post("/", SaveWorkspace)
		return r
	}())

	r.With(UserOnly).With(UserWithWorkspaceOnly).With(MustHaveDateParam).Get("/", ShowIndex)

	fmt.Printf("Listening to http://%v\n", config.Addr)
	err = http.ListenAndServe(config.Addr, r)
	if err != nil {
		log.Fatal(err)
	}
}

// here is source of problem with 307 redirect
func MustHaveDateParam(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dateStr := r.URL.Query().Get("date")
		if dateStr != "" {
			_, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				next.ServeHTTP(w, r)
			} else {
				log.Printf("[ERROR] cannot parse date str '%v': %v", dateStr, err)
			}
		} else {
			date := time.Now()
			u := fmt.Sprintf("/?date=%s", date.Format("2006-01-02"))
			http.Redirect(w, r, u, http.StatusTemporaryRedirect)
		}
	})
}

func UserWithWorkspaceOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, ok := ctx.Value(userContextKey).(*User)
		if !ok {
			log.Printf("[EROR] cannot marshall User from context")
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if user.WorkspaceId == -1 {
			log.Printf("[DEBUG] user must select workspace")
			http.Redirect(w, r, "/workspace", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetStrictCookie(r, CookieToken)
		if err != nil {
			log.Printf("Cannot get cookie 'togglToken': %v", err)
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		workspaceID := -1
		workspaceIDs, err := GetStrictCookie(r, CookieWorkspace)
		if err != nil {
			log.Printf("Cannot get cookie 'togglWorkspace': %v", err)
		} else {
			workspaceID, err = strconv.Atoi(workspaceIDs)
			if err != nil {
				log.Printf("Cannot parse togglWorkspace as int %v", err)
				workspaceID = -1
			}
		}

		storage.CreateUser(token, workspaceID)

		ctx := context.WithValue(r.Context(), userContextKey, &User{
			Toggl:       toggl.NewToggl(token),
			WorkspaceId: workspaceID,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ShowWorkspaces(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(*User)
	if !ok {
		log.Println("[ERROR] cannot get user from context")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	workspaces, err := user.Toggl.GetWorkspaces()
	if err != nil {
		log.Printf("[ERROR] cannot load workspaces from toggle: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	me, err := user.Toggl.GetMe()
	if err != nil {
		log.Printf("[ERROR] cannot load user info from toggle: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderer.RenderHTML(w, "workspaces", WorkspacesPage{
		User:       me,
		Workspaces: workspaces,
	})
}

func SaveWorkspace(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("[ERR] parse form %v", err)
		return
	}
	workspaceID := r.FormValue("workspace_id")
	http.SetCookie(w, &http.Cookie{
		Name:   CookieWorkspace,
		Value:  workspaceID,
		MaxAge: CookieMaxAge,
		Secure: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ShowIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		renderer.RenderHTML(w, "login", map[string]string{
			"Instructions": config.Instructions,
		})
		return
	}

	dateStr := r.URL.Query().Get("date")
	startDate := time.Now()
	if dateStr != "" {
		dt, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("[ERROR] cannot parse date str '%v': %v", dateStr, err)
		} else {
			startDate = dt
		}
	}

	reporter := report.Reporter{
		TogglClient:          user.Toggl,
		ProjectTimePrecision: time.Duration(config.Reporter.ProjectTimePrecision) * time.Second,
		TaskTimePrecision:    time.Duration(config.Reporter.TaskTimePrecision) * time.Second,
	}
	dailyReport, err := reporter.BuildDailyReport(user.WorkspaceId, startDate)
	if err != nil {
		log.Printf("[ERROR] %v", err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	weeklyReport, err := reporter.BuildReport(user.WorkspaceId, startDate.AddDate(0, 0, -7), startDate)
	if err != nil {
		log.Printf("[ERROR] %v", err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	reportJson, err := json.Marshal(dailyReport)
	if err != nil {
		log.Printf("[ERROR] %v", err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	achievementsList := make(map[string]achievements.UserAchievement)
	for k, v := range achievements.AchievementsList {
		achievementsList[k] = v
	}

	for _, achievement := range achievementsList {
		if achievement.CheckCommand(dailyReport) {
			achievement.IsUnlocked = true
			achievementsList[achievement.Name] = achievement
		} else if achievement.CheckWeeklyCommand(weeklyReport) {
			achievement.IsUnlocked = true
			achievementsList[achievement.Name] = achievement
		}
	}

	pageData := ReportPage{
		Report:          dailyReport,
		FormData:        formGenerator.ConvertReportToForms(dailyReport),
		RedmineData:     redmineGenerator.BuildRedmineReportForms(dailyReport),
		AchievementsMap: achievementsList,
		ReportJSON:      string(reportJson),
	}

	renderer.RenderHTML(w, "index", pageData)
}

func ShowLogin(w http.ResponseWriter, r *http.Request) {
	renderer.RenderHTML(w, "login", map[string]string{
		"Instructions": config.Instructions,
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("[ERR] parse form %v", err)
		return
	}
	token := r.FormValue("token")
	pwd := r.FormValue("password")

	if pwd != config.Password {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   CookieToken,
		Value:  token,
		MaxAge: CookieMaxAge,
		Secure: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
