package utils

import "time"

func GetCurrentTime() string {
	timeZone, _ := time.LoadLocation("America/Los_Angeles")
	return time.Now().UTC().In(timeZone).Format("2006-01-02T15:04:05Z")
}
