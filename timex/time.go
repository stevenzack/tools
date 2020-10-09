package timex

import "time"

// month
func NextMonthBy(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
}
func NextMonth() time.Time {
	return NextMonthBy(time.Now())
}

func LastMonthBy(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()-1, t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

func ThisMonthBy(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}
func ThisMonth() time.Time {
	return ThisMonthBy(time.Now())
}
func LastMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
}

// week
func LastWeekBy(t time.Time) time.Time {
	left := t.Weekday()
	return time.Date(t.Year(), t.Month(), t.Day()-7-int(left), 0, 0, 0, 0, t.Location())
}
func ThisWeekBy(t time.Time) time.Time {
	left := t.Weekday()
	return time.Date(t.Year(), t.Month(), t.Day()-int(left), 0, 0, 0, 0, t.Location())
}
func NextWeekBy(t time.Time) time.Time {
	left := t.Weekday()
	return time.Date(t.Year(), t.Month(), t.Day()+7-int(left), 0, 0, 0, 0, t.Location())
}
func LastWeek() time.Time {
	return LastWeekBy(time.Now())
}
func ThisWeek() time.Time {
	return ThisWeekBy(time.Now())
}
func NextWeek() time.Time {
	return NextWeekBy(time.Now())
}

// day
func LastDayBy(t time.Time) time.Time {
	return t.Add(-time.Hour * 24)
}

func Yesterday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
}
func Tomorrow() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
}
