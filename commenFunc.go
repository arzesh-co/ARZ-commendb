package CommenDb

import (
	"strconv"
	"time"
)

func ConvertorToInt64(n any) int64 {
	switch n := n.(type) {
	case int:
		return int64(n)
	case int8:
		return int64(n)
	case int16:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return int64(n)
	case string:
		num, err := strconv.ParseInt(n, 10, 64)
		if err == nil {
			return int64(0)
		}
		return num
	}
	return int64(0)
}
func today(days any) any {
	year, month, day := time.Now().Date()
	theTime := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return theTime.Unix()
}
func nDayBefore(days any) any {
	intDay := ConvertorToInt64(days)
	return time.Now().Add((time.Duration(intDay) * -24) * time.Hour).Unix()
}
func nDayAfter(days any) any {
	intDay := ConvertorToInt64(days)
	return time.Now().Add((time.Duration(intDay) * -24) * time.Hour).Unix()
}
func BeginOfThisYear(days any) any {
	year, _, _ := time.Now().Date()
	theTime := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	return theTime.Unix()
}
func BeginOfThisMonth(days any) any {
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

func SetCommenFunc() map[string]func(days any) any {
	funcs := make(map[string]func(days any) any)
	funcs["today"] = today
	funcs["nday_before"] = nDayBefore
	funcs["nday_after"] = nDayAfter
	funcs["begin_of_this_year"] = BeginOfThisYear
	funcs["begin_of_this_month"] = BeginOfThisMonth
	return funcs
}
