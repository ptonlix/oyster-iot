package utils

import (
	"strconv"
	"time"
)

// 获取月份的每天的时间戳范围
func GetMonthDayTimeInterval(myYear string, myMonth string) ([]map[int64]int64, error) {
	// 数字月份必须前置补零
	if len(myMonth) == 1 {
		myMonth = "0" + myMonth
	}
	yInt, err := strconv.Atoi(myYear)
	if err != nil {
		return nil, err
	}
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(timeLayout, myYear+"-"+myMonth+"-01 00:00:00", loc)
	if err != nil {
		return nil, err
	}
	newMonth := theTime.Month()

	// 获取当前月份的天数
	days := count(yInt, int(newMonth))
	result := make([]map[int64]int64, days)
	for day := 1; day <= days; day++ {
		dayMap := make(map[int64]int64)
		t1 := time.Date(yInt, newMonth, day, 0, 0, 0, 0, time.Local).Unix()
		if day == days {
			t2 := time.Date(yInt, newMonth+1, 0, 0, 0, 0, 0, time.Local).Unix()
			dayMap[t1] = t2
		} else {
			t2 := time.Date(yInt, newMonth, day+1, 0, 0, 0, 0, time.Local).Unix()
			dayMap[t1] = t2
		}
		result[day-1] = dayMap
	}
	return result, nil
}

func count(year int, month int) (days int) {
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30

		} else {
			days = 31
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	}
	return
}
