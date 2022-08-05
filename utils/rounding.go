package utils

import "time"

func RoundTime(duration time.Duration, precision time.Duration) time.Duration {
	r := duration % precision
	if r == 0 {
		return duration
	}
	return duration + (precision - r)
}
