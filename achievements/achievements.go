package achievements

import (
	"goreporter/report"
	"goreporter/utils"
	"time"
)

//need to refactor to have more than one type of check command - daily, weekly, 14 days etc
// To send appropriate report from outside without checking, what type of command have specific achievement
// And maybe to have different/random text description depending on day/just random

type UserAchievement struct {
	Name                 string
	IsUnlocked           bool
	CheckCommand         func(report report.Report) bool
	CheckWeeklyCommand   func(report report.Report) bool
	ConditionDescription string
	ImagePath            string
}

var TimeWizardAchievement string = "TimeWizard"
var TimeTurnerAchievement = "TimeTurner"
var LongTimeNoSeeAchievement = "LongTimeNoSee"
var IsItSalaryTimeAlreadyAchievement = "IsItSalaryTimeAlready"
var TypeYourTextHereAchievement = "TypeYourTextHere"
var LongRestAchievement = "LongRest"
var HollowMemeAchievement = "HollowMeme"
var YearAchievement = "YearAchievement"
var FilleanAchievement = "FilleanAchievement"

var AchievementsList = map[string]UserAchievement{
	TimeWizardAchievement: UserAchievement{
		Name:                 TimeWizardAchievement,
		IsUnlocked:           false,
		ConditionDescription: "Space Worker - Tracked >= 20 hours today",
		CheckCommand:         checkTimeWizard,
		CheckWeeklyCommand:   returnFalseCheck,
		ImagePath:            "/static/images/timeWizard.svg",
	},
	TimeTurnerAchievement: UserAchievement{
		Name:                 TimeTurnerAchievement,
		IsUnlocked:           false,
		ConditionDescription: "Time-Turner - Tracked 2 events with colliding time of 5+ minutes today",
		CheckCommand:         checkTimeTurner,
		CheckWeeklyCommand:   returnFalseCheck,
		ImagePath:            "/static/images/timeTurner.svg",
	},
	LongTimeNoSeeAchievement: UserAchievement{
		Name:                 LongTimeNoSeeAchievement,
		IsUnlocked:           false,
		ConditionDescription: "Where did I put my report? - Opened report page at previous week",
		CheckCommand:         checkLongTimeNoSee,
		CheckWeeklyCommand:   returnFalseCheck,
		ImagePath:            "/static/images/longTimeNoSee.svg",
	},
	IsItSalaryTimeAlreadyAchievement: UserAchievement{
		Name:                 IsItSalaryTimeAlreadyAchievement,
		IsUnlocked:           false,
		ConditionDescription: "Is It Salary Time Already? - Opened report page at two weeks before",
		CheckCommand:         checkIsItSalaryTimeAlready,
		CheckWeeklyCommand:   returnFalseCheck,
		ImagePath:            "/static/images/isItSalaryTimeAlready.svg",
	},
	TypeYourTextHereAchievement: UserAchievement{
		Name:                 TypeYourTextHereAchievement,
		IsUnlocked:           false,
		ConditionDescription: "Type Your Text Here - Have 1+ event tracked with no message for today",
		CheckCommand:         checkTypeYourTextHere,
		CheckWeeklyCommand:   returnFalseCheck,
		ImagePath:            "/static/images/typeYourTextHere.svg",
	},
	LongRestAchievement: UserAchievement{
		Name:                 LongRestAchievement,
		IsUnlocked:           false,
		ConditionDescription: "Ghost in a Tracker - Have no events tracked for 7 days (excluding today)",
		CheckCommand:         returnFalseCheck,
		CheckWeeklyCommand:   checkLongRestWeekly,
		ImagePath:            "/static/images/longRest.svg",
	},
	HollowMemeAchievement: UserAchievement{
		Name:       HollowMemeAchievement,
		IsUnlocked: false,
		//Hidden condition: Enter tracker from 31 October to 2 November after 19:00
		ConditionDescription: "May I have some treat :-)? Or you prefer empty report sheet?",
		CheckCommand:         checkHollowMeme,
		CheckWeeklyCommand:   returnFalseCheck,
		ImagePath:            "/static/images/hollowMeme.svg",
	},
	YearAchievement: UserAchievement{
		Name:       YearAchievement,
		IsUnlocked: false,
		//Hidden condition: Enter tracker from 31 December after 19:00
		ConditionDescription: "Let's go celebrate and have some presents!",
		CheckCommand:         checkYear,
		CheckWeeklyCommand:   returnFalseCheck,
		ImagePath:            "/static/images/year.svg",
	},
	FilleanAchievement: UserAchievement{
		Name:       FilleanAchievement,
		IsUnlocked: false,
		//Hidden condition: Tracked on Fillean project today
		ConditionDescription: "O'rly?",
		CheckCommand:         checkFillean,
		CheckWeeklyCommand:   returnFalseCheck,
		ImagePath:            "/static/images/fillean.svg",
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
			}
		}
	}
	return false
}

func checkTimeWizard(report report.Report) bool {
	return report.TotalDuration >= 20*time.Hour
}

//Check that current date (At) is during previous week from current day
func checkLongTimeNoSee(report report.Report) bool {
	timeNow := time.Now()
	// need to shift to Sunday of last week, then add 7 days before
	timeday := timeNow.Weekday()
	// 0 - Sunday
	if timeday == 0 {
		timeday = 7
	}
	totalShift := int(timeday) + 7

	timeLastWeekBefore := timeNow.AddDate(0, 0, -1*totalShift)
	timeThisWeekStart := timeNow.AddDate(0, 0, -1*int(timeday)+1)

	return report.At.After(timeLastWeekBefore) && report.At.Before(timeThisWeekStart)
}

//Check that current date (At) is (two weeks) 14 or more days before current day
func checkIsItSalaryTimeAlready(report report.Report) bool {
	timeNow := time.Now()
	// need to shift to start of week, then add 7 days before
	timeday := timeNow.Weekday()
	// 0 - Sunday
	if timeday == 0 {
		timeday = 7
	}
	//from Sunday of previous week to Monday of previous week
	weekAgoShift := int(timeday) + 7

	timeWeekBefore := timeNow.AddDate(0, 0, -1*weekAgoShift)
	timeTwoWeeksBefore := timeWeekBefore.AddDate(0, 0, -7)

	return report.At.After(timeTwoWeeksBefore) && report.At.Before(timeWeekBefore)
}

func checkTypeYourTextHere(report report.Report) bool {
	for _, report_item := range report.RawTimeEntries {
		if report_item.Description == "" {
			return true
		}
	}
	return false
}

func returnFalseCheck(report report.Report) bool {
	return false
}

func checkLongRestWeekly(report report.Report) bool {
	if len(report.RawTimeEntries) == 0 {
		return true
	} else {
		year, month, day := time.Now().Date()
		timeStartOfDay := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
		for _, report_item := range report.RawTimeEntries {
			if report_item.Start.Before(timeStartOfDay) {
				return false
			}
		}
		return true
	}
}

func checkHollowMeme(report report.Report) bool {
	timeChecked := time.Now()
	_, month, day := timeChecked.Date()
	if (month == time.October && day == 31) || (month == time.November && day > 0 && day < 2) {
		hour, _, _ := timeChecked.Clock()
		if hour > 19 {
			return true
		}
	}
	return false
}

func checkYear(report report.Report) bool {
	timeChecked := time.Now()
	_, month, day := timeChecked.Date()
	if month == time.December && day == 31 {
		hour, _, _ := timeChecked.Clock()
		if hour > 19 {
			return true
		}
	}
	return false
}

func checkFillean(report report.Report) bool {
	for _, project := range report.Projects {
		if project.Name == "invest-fillean" {
			return true
		}
	}
	return false
}
