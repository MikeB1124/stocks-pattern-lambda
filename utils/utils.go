package utils

import "time"

func GetCurrentTime() *time.Time {
	timeNow := time.Now().UTC()
	return &timeNow
}
