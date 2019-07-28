package monitor

import (
	"time"
)

const (
	// 08/05/2017 11:45:22
	cadDateLayout      = "01/02/2006 15:04:05"
	cadDateLayoutShort = "01/02/2006"
)

func CadDateTime(orig string) time.Time {
	return dateTime(orig)
}

func CadDateTimeShort(orig string) time.Time {
	return dateTimeShort(orig)
}

func dateTime(orig string) time.Time {
	t, _ := time.Parse(cadDateLayout, orig)
	return t
}

func dateTimeShort(orig string) time.Time {
	t, _ := time.Parse(cadDateLayoutShort, orig)
	return t
}
