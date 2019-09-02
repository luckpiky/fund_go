package util

import (
	"strings"
	"strconv"
	"time"
)


func GetCurMonthFirstDay(curDate string) (firstDate string) {
	d1 := strings.Split(curDate, "-")
	if len(d1) == 0 {
		return ""
	}

	return d1[0] + "-" + d1[1] + "-" + "01" + " 00:00:00"
}

func GetCurMonthLastDay(curDate string) (firstDate string) {
	d1 := strings.Split(curDate, "-")
	if len(d1) == 0 {
		return ""
	}

	day := d1[0] + "-" + d1[1] + "-"

	switch d1[1] {
	case "01", "03", "05", "07", "08", "10", "12", "1", "3", "5", "7", "8":
		return day + "31" + " 00:00:00"
	case "04", "06", "09", "11", "4", "6", "9":
		return day + "30" + " 00:00:00"
	}

	year,_ := strconv.Atoi(d1[0])
	if year % 4 == 0 {
		return day + "29" + " 00:00:00"
	}

	return day + "28" + " 00:00:00"
}

func GetNextMonthFirstDay(curDate string) (firstDate string) {
	d1 := strings.Split(curDate, "-")
	if len(d1) == 0 {
		return ""
	}

	year_str := ""

	d2,_ := strconv.Atoi(d1[1])
	if d2 == 12 {
		d2 = 1
		year,_ := strconv.Atoi(d1[0])
		year = year + 1

		year_str = strconv.Itoa(year)
	} else {
		d2 = d2 + 1
		year_str = d1[0] 
	}

	return year_str + "-" + strconv.Itoa(d2) + "-" + "01" + " 00:00:00"
}

func GetNextMonthLastDay(curDate string) (firstDate string) {
	nextMonthDay := GetNextMonthFirstDay(curDate)
	return GetCurMonthLastDay(nextMonthDay)
}

func TimeStr2Int64(timeStr string) (int64) {
    stamp, _ := time.ParseInLocation("2006-01-02 03:04:05", timeStr, time.Local)
	return stamp.Unix()
}

func TimeStr2Int64_2(timeStr string) (int64) {
    stamp, _ := time.ParseInLocation("2006-1-2 03:04:05", timeStr, time.Local)
	return stamp.Unix()
}

func Int64ToMonthStr(in int64) (out string) {
	tm := time.Unix(in, 0)
	out = tm.Format("2006-01")
	return out
}

func Int64ToDayStr(in int64) (out string) {
	tm := time.Unix(in, 0)
	out = tm.Format("2006-01-01")
	return out
}

func TimeStrToDay(in string) (out string) {
	// 2006-1-2 03:04:05 to 2006-1-2
	s := strings.Split(in, " ")
	if len(s) >= 2 {
		return s[0]
	}

	return in
}