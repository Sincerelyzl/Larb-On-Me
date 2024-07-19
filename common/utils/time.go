package utils

import "time"

func GetNowUTCTime() time.Time {
	nowTime := time.Now()
	// location bangkok
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return nowTime.In(loc).UTC()
}
