package utils

import (
	"fmt"
	"time"
)

func FormatDuration(duration time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%02d", Hours(duration), Minutes(duration), Seconds(duration))
}

func Hours(duration time.Duration) int {
	return int(duration / time.Hour)
}

func Minutes(duration time.Duration) int {
	return int(duration/time.Minute) % 60
}

func Seconds(duration time.Duration) int {
	return int(duration/time.Second) % 60
}
