package timex

import (
	"time"
)

// day
func BetweenYesterday() (time.Time, time.Time) {
	yesterday := time.Now().Add(-time.Hour * 24)
	return DayBetween(yesterday)
}

func BetweenToday() (time.Time, time.Time) {
	return DayBetween(time.Now())
}

func DayBetween(t time.Time) (time.Time, time.Time) {
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
	return start, end
}

// week
func BetweenThisWeek()(time.Time,time.Time){
	return WeekBetween(time.Now())
}
func WeekBetween(t time.Time) (time.Time, time.Time) {
	return ThisWeekBy(t), NextWeekBy(t)
}
func BetweenLastWeek() (time.Time, time.Time) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day()-7, 0, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return start, end
}

// month
func BetweenLastMonth() (time.Time, time.Time) {
	lastMonth := LastMonth()
	return MonthBetween(lastMonth)
}
func BetweenThisMonth() (time.Time, time.Time) {
	return MonthBetween(time.Now())
}
func MonthBetween(t time.Time) (time.Time, time.Time) {
	return ThisMonthBy(t), NextMonthBy(t)
}
