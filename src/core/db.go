// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-28
package core

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	gopg "github.com/go-pg/pg/v10"
	"github.com/teocci/go-mavlink-parser/src/datamgr"
	"github.com/teocci/go-mavlink-parser/src/gcs"
	"github.com/teocci/go-mavlink-parser/src/model"
	"github.com/teocci/go-mavlink-parser/src/timemgr"
)

type DBLogger struct {
	DBMgr   *gopg.DB
	Flight  model.Flight
	Inserts int
	// Buffered channel of outbound messages.
	Insert    chan model.FlightRecord
	Done      chan struct{}
	Interrupt chan os.Signal
}

var (
	db *gopg.DB

	prev *datamgr.RTT

	isFirstRecord = true
)

func NewDBLogger(c datamgr.InitConf) *DBLogger {
	flight := model.Flight{
		ID: c.FlightID,
	}

	flight.Select(db)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	dbLogger := &DBLogger{
		DBMgr:     db,
		Flight:    flight,
		Inserts:   0,
		Insert:    make(chan model.FlightRecord, 256),
		Done:      make(chan struct{}),
		Interrupt: interrupt,
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go dbLogger.onRecordMessage()

	return dbLogger
}

func (c *DBLogger) onRecordMessage() {
	for {
		select {
		case <-c.Done:
			return
		case r, ok := <-c.Insert:
			if !ok {
				return
			}

			c.insert(r)

			// Add queued chat messages to the current websocket r.
			n := len(c.Insert)
			for i := 0; i < n; i++ {
				c.insert(<-c.Insert)
			}
		case <-c.Interrupt:
			log.Println("onRecordMessage-> interrupt")
			c.updateFlight()

			// Close file
			select {
			case <-c.Done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (c *DBLogger) insert(r model.FlightRecord) {
	if r.Insert(c.DBMgr) {
		c.Flight.Duration += r.Duration
		c.Flight.Distance += r.Distance
		c.Inserts++
	}
}

func (c *DBLogger) updateFlight() {
	log.Printf("recv: %#v", c.Flight)
	if c.Inserts > 0 && c.Flight.Distance > 0 && c.Flight.Duration > 0 {
		c.Flight.Status |= model.FlightStatusCompleted | model.FlightStatusProcessed
		c.Flight.Update(c.DBMgr)
	}
}

func (c *DBLogger) Close() {
	c.Close()
}

func insertRecord(rtt *datamgr.RTT) {
	rttTimeBoot := timemgr.UInt32ToUnixTime(rtt.TimeBootMs)

	var prevTimeBoot time.Time
	var duration int64
	var distance float32
	var speed float32

	if rtt.Seq > 0 {
		prevTimeBoot = timemgr.UInt32ToUnixTime(prev.TimeBootMs)
		duration = rttTimeBoot.Sub(prevTimeBoot).Milliseconds()

		orig := gcs.SCS{Lat: float64(prev.Lat), Lon: float64(prev.Lon)}
		dest := gcs.SCS{Lat: float64(rtt.Lat), Lon: float64(rtt.Lon)}

		distance = float32(orig.MetersTo(dest))

		if duration > 0 {
			speed = distance / float32(duration)
		}
	}

	record := model.FlightRecord{}
	record.Parse(*rtt)
	record.Duration = duration
	record.Distance = distance
	record.Speed = speed

	if isFirstRecord {
		prev = rtt
		isFirstRecord = false
	}
}
