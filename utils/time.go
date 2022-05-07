package utils

import "time"

func ParseTime(str string) time.Time {
	date, _ := time.Parse("2006-01-02 15:04:05", str)
	return date
}

func FormatTime(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}
