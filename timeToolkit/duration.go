package timeToolkit

import "time"

const (
	DAY   = time.Hour * 24
	WEEK  = DAY * 7
	MONTH = DAY * 30
	YEAR  = MONTH * 12
)
