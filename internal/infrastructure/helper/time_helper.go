package helper

import (
	"time"
)

type TimeStage int

const (
	Morning     TimeStage = 9
	OpenBooking           = 14
	Evening               = 18
)

func EveningTimeWithTimeZone(timeZone string) (time.Time, error) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return TodayZeroTimeUTC(), err
	}
	date := time.Now()
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), Evening, 0, 0, 0, loc)
	return bookDate, nil
}

func TimeWithTimeZone(hour TimeStage, timeZone string) (time.Time, error) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return TodayZeroTimeUTC(), err
	}
	date := time.Now().In(loc)
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), int(hour), 0, 0, 0, loc)
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

func TodayPlusUTC(days int) time.Time {
	date := time.Now()
	zeroDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	return zeroDate.AddDate(0, 0, days)
}

func PlusDaysUTC(startDate time.Time, days int) time.Time {
	zeroDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	return zeroDate.AddDate(0, 0, days)
}

func TomorrowZeroTimeUTC() time.Time {
	currentTime := time.Now()
	date := currentTime.AddDate(0, 0, 1)
	bookDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	return bookDate
}
