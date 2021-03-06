// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-01
package model

import (
	"log"
	"time"

	gopg "github.com/go-pg/pg/v10"
	"github.com/teocci/go-mavlink-parser/src/datamgr"
)

type FlightRecord struct {
	ID             int64     `json:"id" csv:"id" pg:"id,pk,unique"`
	DroneID        int64     `json:"drone_id" csv:"drone_id" pg:"drone_id"`
	FlightID       int64     `json:"flight_id" pg:"flight_id" pg:"flight_id"`
	Sequence       int64     `json:"sequence" csv:"seq" pg:"sequence"`
	Duration       int64     `json:"duration" csv:"duration" pg:"duration"`
	Distance       float32   `json:"distance" csv:"distance" pg:"distance"`
	Speed          float32   `json:"speed" csv:"speed" pg:"speed"`
	Latitude       float32   `json:"latitude" csv:"lat" pg:"latitude"`
	Longitude      float32   `json:"longitude" csv:"long" pg:"longitude"`
	Altitude       float32   `json:"altitude" csv:"alt" pg:"altitude"`
	Roll           float32   `json:"roll" csv:"roll" pg:"roll"`
	Pitch          float32   `json:"pitch" csv:"pitch" pg:"pitch"`
	Yaw            float32   `json:"yaw" csv:"yaw" pg:"yaw"`
	BatVoltage     float32   `json:"battery_voltage" csv:"battery_voltage" pg:"battery_voltage"`
	BatCurrent     float32   `json:"battery_current" csv:"battery_current" pg:"battery_current"`
	BatPercent     float32   `json:"battery_percentage" csv:"battery_percentage" pg:"battery_percentage"`
	BatTemperature float32   `json:"battery_temperature" csv:"battery_temperature" pg:"battery_temperature"`
	Temperature    float32   `json:"temperature" csv:"temperature" pg:"temperature"`
	GPSStatus      float32   `json:"gps_status" csv:"gps_status" pg:"gps_status"`
	DroneStatus    float32   `json:"drone_status" csv:"drone_status" pg:"drone_status"`
	LastUpdate     time.Time `json:"last_update" csv:"last_update" pg:"last_update"`
}

func (fsr *FlightRecord) Select(db *gopg.DB) bool {
	err := db.Model(fsr).
		Where("drone_id = ?", fsr.DroneID).
		Where("flight_id = ?", fsr.FlightID).
		Where("sequence = ?", fsr.Sequence).
		Select()
	if err != nil {
		log.Println(err)
		return false
	}

	return fsr.ID > 0
}

func (fsr *FlightRecord) InsertWithCheck(db *gopg.DB) bool {
	res, err := db.Model(fsr).OnConflict("DO NOTHING").Insert()
	if err != nil {
		log.Println(err)
		return false
	}

	if res.RowsAffected() > 0 {
		//fmt.Println("FlightRecord inserted")
		return true
	}

	fr := &FlightRecord{
		DroneID:  fsr.DroneID,
		FlightID: fsr.FlightID,
		Sequence: fsr.Sequence,
	}

	return fr.Select(db)
}

func (fsr *FlightRecord) Insert(db *gopg.DB) bool {
	res, err := db.Model(fsr).OnConflict("DO NOTHING").Insert()
	if err != nil {
		log.Println(err)
		return false
	}

	if res.RowsAffected() > 0 {
		//fmt.Println("FlightRecord inserted")
		return true
	}

	return false
}

func (fsr *FlightRecord) Parse(rtt datamgr.RTT) {
	fsr.DroneID = rtt.DroneID
	fsr.FlightID = rtt.FlightID
	fsr.Sequence = rtt.Seq
	fsr.Latitude = float32(rtt.Lat) / 1e7
	fsr.Longitude = float32(rtt.Lon) / 1e7
	fsr.Altitude = float32(rtt.Alt) / 1e3
	fsr.Roll = rtt.Roll
	fsr.Pitch = rtt.Pitch
	fsr.Yaw = rtt.Yaw
	fsr.LastUpdate = rtt.LastUpdate
}
