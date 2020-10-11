package timeToolkit

import "time"

func RangeBackByDate(end time.Time, each func(t time.Time, date string) error) error {
	for cursor := LastDayBy(time.Now()); !cursor.Before(end); cursor = LastDayBy(cursor) {
		e := each(cursor, cursor.Format(LAYOUT_DATE))
		if e != nil {
			return e
		}
	}
	return nil
}

func RangeBackByWeek(end time.Time, each func(t time.Time, date string) error) error {
	for cursor := LastWeek(); !cursor.Before(end); cursor = LastWeekBy(cursor) {
		e := each(cursor, cursor.Format(LAYOUT_DATE))
		if e != nil {
			return e
		}
	}
	return nil
}

func RangeBackByMonth(end time.Time, each func(t time.Time, month string) error) error {
	for cursor := LastMonth(); !cursor.Before(end); cursor = LastMonthBy(cursor) {
		e := each(cursor, cursor.Format(LAYOUT_MONTH))
		if e != nil {
			return e
		}
	}
	return nil
}
