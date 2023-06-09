package utils

import "time"

func GetEuropeTime() time.Time {
	return time.Now().UTC().Add(time.Hour * 3)
}

func GetTodayAndYesterday() (time.Time, time.Time) {
	t := GetEuropeTime()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	yesterday := time.Date(t.Year(), t.Month(), t.Day()-1, 0, 0, 0, 0, t.Location())
	return yesterday, today
}
