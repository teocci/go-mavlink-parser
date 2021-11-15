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

	prev         model.FlightRecord
	prevTimeBoot time.Time

	isFirstRecord = true

	isProcessing = false
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

			c.updateFlight(false)
		case <-c.Interrupt:
			log.Println("onRecordMessage-> interrupt")
			c.updateFlight(true)

			// Close file
			select {
			case <-c.Done:
			case <-time.After(2 * time.Second):
			}
			return
		}
	}
}

func (c *DBLogger) insert(r model.FlightRecord) {
	if r.Insert(c.DBMgr) {
		c.Flight.Length++
		c.Flight.Duration += r.Duration
		c.Flight.Distance += r.Distance
		c.Inserts++
	}
}

func (c *DBLogger) updateFlight(close bool) {
	if c.Flight.Length > 0 {
		if !isProcessing {
			c.Flight.Status |= model.FlightStatusProcessed
			isProcessing = true
		}
		if close {
			c.Flight.Status |= model.FlightStatusCompleted
		}
		ok := c.Flight.Update(c.DBMgr)

		if close && ok {
			log.Printf("flight: %#v\n", c.Flight)
		}
	}
}

func (c *DBLogger) Close() {
	c.Close()
}

func insertRecord(rtt *datamgr.RTT) {
	record := model.FlightRecord{}
	record.Parse(*rtt)

	rttTimeBoot := timemgr.UInt32ToUnixTime(rtt.TimeBootMs)

	var duration int64 = 0
	var distance float32 = 0
	var speed float32 = 0

	if !isFirstRecord {
		duration = rttTimeBoot.Sub(prevTimeBoot).Milliseconds()

		orig := gcs.SCS{Lat: float64(prev.Latitude), Lon: float64(prev.Longitude)}
		dest := gcs.SCS{Lat: float64(record.Latitude), Lon: float64(record.Longitude)}

		distance = float32(orig.MetersTo(dest))

		if duration > 0 {
			speed = distance / float32(duration)
		}
	}

	record.Duration = duration
	record.Distance = distance
	record.Speed = speed

	dbl.Insert <- record

	if isFirstRecord {
		prev = record
		isFirstRecord = false
	}

	prevTimeBoot = rttTimeBoot
}
