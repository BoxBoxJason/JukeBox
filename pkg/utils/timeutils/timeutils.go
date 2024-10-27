package timeutils

import "time"

func GetDatetimeString() string {
	current_time := time.Now()
	return current_time.Format("2006/01/02-15h04:05")
}

func ConvertUnixTimestampToDatetimeString(unix_timestamp int) string {
	timestamp := time.Unix(int64(unix_timestamp), 0)
	return timestamp.Format("2006/01/02-15h04:05")
}
