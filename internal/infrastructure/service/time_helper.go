package service

import (
	"time"
)

func TodayZeroTimeUTC() time.Time {
	date := time.Now()
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	return bookDate
}

func TomorrowZeroTimeUTC() time.Time {
	currentTime := time.Now()
	date := currentTime.AddDate(0, 0, 1)
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	return bookDate
}
