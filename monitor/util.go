package monitor

import "strconv"

func mustAtoi(s string) int {
	x, _ := strconv.Atoi(s)
	return x
}
