package timeToolkit

import "time"

func TillTomorrow() time.Duration {
	return Tomorrow().Sub(time.Now())
}
func TillNextWeek() time.Duration {
	return NextWeek().Sub(time.Now())
}

func TillNextMonth() time.Duration {
	return NextMonth().Sub(time.Now())
}
