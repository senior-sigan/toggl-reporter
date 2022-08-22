package achievements

import "goreporter/toggl"

type UserAchievement struct {
	Name                 string
	IsUnlocked           bool
	Condition            string
	ConditionDescription string
}

//var hardWorkerName string = "HardWorker"
var FullTimeAchievement string = "Full-Time"

var DedicatedWorkerAchievement string = "Dedicated worker"

var AchievementsList = map[string]UserAchievement{
	// hardWorkerName: UserAchievement{
	// 	Name:                 hardWorkerName,
	// 	IsUnlocked:           false,
	// 	Condition:            "",
	// 	ConditionDescription: "working time >= 10 hours for 7 days",
	// },
	FullTimeAchievement: UserAchievement{
		Name:                 FullTimeAchievement,
		IsUnlocked:           false,
		Condition:            "",
		ConditionDescription: "working time >= 8 hours for today",
	},
	DedicatedWorkerAchievement: UserAchievement{
		Name:                 DedicatedWorkerAchievement,
		IsUnlocked:           false,
		Condition:            "",
		ConditionDescription: "working at project >= 6 hours for today",
	},
}

func (ua *UserAchievement) CheckIfUnlocked(toggl toggl.Toggl, workspaceId int) bool {
	return false
}
