// Package timemgr
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-30
package timemgr

import (
	"math"
	"time"
)

const (
	TimeLayoutUriaio         = "2006-01-02 15:04:05.000"
	TimeLayoutSmati          = "2006-01-02T15:04:05.999"

	msInSecond       int64   = 1e3
	nsInMillisecond  int64   = 1e6

	sInNanosecond    float64 = 1e9
)

func ParseTime(layout, date string) time.Time {
	t, _ := time.Parse(layout, date)
	return t
}

func Int64ToUnixTime(ms int64) time.Time {
	return time.Unix(ms/msInSecond, (ms%msInSecond)*nsInMillisecond)
}

func UInt32ToUnixTime(ms uint32) time.Time {
	return Int64ToUnixTime(int64(ms))
}

func Float32ToUnixTime(gpsTime float32) time.Time {
	sec, dec := math.Modf(float64(gpsTime))
	return time.Unix(int64(sec), int64(dec*sInNanosecond))
}

func GenBaseDate(day int) time.Time {
	t := time.Now()
	return time.Date(2021, 8, day, 13, 0, 0, 0, t.Location())
}
