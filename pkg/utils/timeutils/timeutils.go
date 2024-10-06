package timeutils

import "time"

func GetDatetimeString() string {
	current_time := time.Now()
	return current_time.Format("2006-01-02-15-04-05")
}
