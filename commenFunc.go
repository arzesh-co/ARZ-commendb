package CommenDb

import "time"

func today() int64 {
	return time.Now().Unix()
}
func nDayBefore(days int64) int64 {
	return time.Now().Add((time.Duration(days) * -24) * time.Hour).Unix()
}
func nDayAfter(days int64) int64 {
	return time.Now().Add((time.Duration(days) * -24) * time.Hour).Unix()
}
func BeginOfThisYear() int64 {
	year, _, _ := time.Now().Date()
	theTime := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	return theTime.Unix()
}
func BeginOfThisMonth() int64 {
	year, month, _ := time.Now().Date()
	theTime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	return theTime.Unix()
}

//func beginOfThisWeek(local int8) int64 {
//	year, month, day := time.Now().Date()
//	weekday := time.Now().Weekday()
//	switch weekday.String() {
//	case "Monday":
//		if local == 1 {
//			day = day - 2
//		} else if local == 3 {
//
//		}
//	}
//	theTime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
//	return theTime.Unix()
//}
