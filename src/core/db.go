// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-28
package core

import (
	gopg "github.com/go-pg/pg/v10"
	"github.com/teocci/go-mavlink-parser/src/datamgr"
	"github.com/teocci/go-mavlink-parser/src/gcs"
	"github.com/teocci/go-mavlink-parser/src/model"
	"github.com/teocci/go-mavlink-parser/src/timemgr"
	"time"
)

var (
	db       *gopg.DB
	previous model.FlightRecord
	record   model.FlightRecord

	prev *datamgr.RTT

	firstRecord = true
)

func insertRecord(rtt *datamgr.RTT) {
	currFCCTime := timemgr.UInt32ToUnixTime(rtt.TimeBootMs)

	var prevFCCTime time.Time
	var duration int64
	var distance float32
	var speed float32

	if rtt.Seq > 0 {
		prevFCCTime = timemgr.UInt32ToUnixTime(prev.TimeBootMs)
		duration = currFCCTime.Sub(prevFCCTime).Milliseconds()

		orig := gcs.SCS{Lat: float64(prev.Lat), Lon: float64(prev.Lon)}
		dest := gcs.SCS{Lat: float64(rtt.Lat), Lon: float64(rtt.Lon)}

		distance = float32(orig.MetersTo(dest))

		if duration > 0 {
			speed = distance / float32(duration)
		}
	}

	record = model.FlightRecord{}
	record.Parse(*rtt)
	record.Duration = duration
	record.Distance = distance
	record.Speed = speed

	record.Insert(db)

	if firstRecord {
		prev = rtt
		firstRecord = false
	}
}
