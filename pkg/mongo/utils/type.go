package utils

import (
	"strconv"
	"time"
)

type DataTimeInt = int

func DataTimeIntFromUnix(unixTimestamp int64) DataTimeInt {
	t := time.Unix(unixTimestamp, 0)
	ts := t.Format("20060102150405")
	v, _ := strconv.Atoi(ts)
	return DataTimeInt(v)
}

func DataTimeIntToUnix(dt DataTimeInt, loc ...*time.Location) int64 {
	l, _ := time.LoadLocation("Asia/Shanghai")
	if len(loc) > 0 {
		l = loc[0]
	}

	t, _ := time.ParseInLocation("20060102150405", strconv.Itoa(int(dt)), l)
	return t.Unix()
}
