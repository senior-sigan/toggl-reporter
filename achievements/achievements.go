package achievements

import (
	"goreporter/report"
	"time"
)

type UserAchievement struct {
	Name                 string
	IsUnlocked           bool
	CheckCommand         func(report report.Report) bool
	ConditionDescription string
	ImagePath            string
}

var TimeWizardAchievement string = "TimeWizard"
var TimeTurnerAchievement = "Time-Turner"

var AchievementsList = map[string]UserAchievement{
	TimeWizardAchievement: UserAchievement{
		Name:                 TimeWizardAchievement,
		IsUnlocked:           false,
		ConditionDescription: "No Time to Stop: Tracked >= 20 hours",
		CheckCommand:         checkTimeWizard,
		ImagePath:            "/static/images/timeWizard.png",
	},
	TimeTurnerAchievement: UserAchievement{
		Name:                 TimeTurnerAchievement,
		IsUnlocked:           false,
		ConditionDescription: "Time-Turner: Tracked 2 events with colliding time of 5+ minutes",
		CheckCommand:         checkTimeTurner,
		ImagePath:            "/static/images/timeTurner.jpg",
	},
}

func checkTimeTurner(report report.Report) bool {
	return true
}

func checkTimeWizard(report report.Report) bool {
	return report.TotalDuration >= 20*time.Hour
	//return true
}
