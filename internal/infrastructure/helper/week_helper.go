package helper

import (
	"time"
)

func WeekRange(year, weekNumber int) []time.Time {
	var days []time.Time
	start := weekStart(year, weekNumber)

	for i := 0; i <= 6; i++ {
		end := start.AddDate(0, 0, i)
		days = append(days, end)
	}
	return days
}

func weekStart(year, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}
