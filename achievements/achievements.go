package achievements

import (
	"goreporter/report"
	"goreporter/utils"
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
		ConditionDescription: "No Time to Stop: Tracked >= 20 hours a day",
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
	//log.Printf("Checking time Turner, len of events is %d", len(report.RawTimeEntries))

	for i, report_item := range report.RawTimeEntries {
		start_timeA := report_item.Start
		end_timeA := report_item.Stop
		for j, report_item_inner := range report.RawTimeEntries {
			if j < i+1 {
				continue
			}
			start_timeB := report_item_inner.Start
			end_timeB := report_item_inner.Stop

			if maxCollision, ok := utils.CheckMaxTimeCollision(start_timeA, end_timeA, start_timeB, end_timeB); ok && maxCollision >= 5*time.Minute {
				return true
			} else if ok {
				//log.Printf("Not achieved - max collision for events is %v", maxCollision)
			}
		}
	}
	return false
}

func checkTimeWizard(report report.Report) bool {
	return report.TotalDuration >= 20*time.Hour
}
