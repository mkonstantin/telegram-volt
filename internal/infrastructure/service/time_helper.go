package service

import (
	"time"
)

const eveningHour = 18

func EveningTimeWithTimeZone(timeZone string) (time.Time, error) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return TodayZeroTimeUTC(), err
	}
	date := time.Now()
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), eveningHour, 0, 0, 0, loc)
	return bookDate, nil
}

func CurrentTimeWithTimeZone(timeZone string) (time.Time, error) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return TodayZeroTimeUTC(), err
	}
	date := time.Now().In(loc)
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), 0, 0, loc)
	return bookDate, nil
}

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
