package utils

import (
	"time"
)

// Check max collision !
func CheckMaxTimeCollision(start_timeA, end_timeA, start_timeB, end_timeB time.Time) (time.Duration, bool) {
	// check who's first
	if start_timeB.Before(start_timeA) {
		start_timeA, start_timeB = start_timeB, start_timeA
		end_timeA, end_timeB = end_timeB, end_timeA
	}

	//partial equality
	if start_timeA.Equal(start_timeB) {
		//log.Printf("Partially equal on start")
		if end_timeA.Before(end_timeB) {
			//log.Printf("Whole A is colliding")
			return end_timeA.Sub(start_timeA), true
		} else {
			//log.Printf("Whole B is colliding")
			return end_timeB.Sub(start_timeB), true
		}
	}
	if end_timeA.Equal(end_timeB) {
		//log.Printf("Partially equal on end")
		if start_timeA.Before(start_timeB) {
			//log.Printf("Whole B is colliding")
			return end_timeB.Sub(start_timeB), true
		} else {
			//log.Printf("Whole A is colliding")
			return end_timeA.Sub(start_timeA), true
		}
	}

	if end_timeB.After(start_timeA) {
		if end_timeB.Before(end_timeA) {
			//log.Printf("Whole B is colliding")
			return end_timeB.Sub(start_timeB), true
		} else {
			//log.Printf("Start of B is colliding")
			return end_timeA.Sub(start_timeB), true
		}
	}

	return 0 * time.Minute, false
}
