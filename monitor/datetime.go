package monitor

import (
	"time"
)

const (
	// 08/05/2017 11:45:22
	DATE_LAYOUT = "01/02/2006 15:04:05"
)

func dateTime(orig string) time.Time {
	t, _ := time.Parse(DATE_LAYOUT, orig)
	return t
}
