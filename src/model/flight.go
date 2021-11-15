// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-01
package model

import (
	"fmt"
	"time"

	gopg "github.com/go-pg/pg/v10"
)

type SessionStatus int
const (
	FlightStatusUnknown          SessionStatus = 0
	FlightStatusGround                         = 1
	FlightStatusAirborne                       = 2
	FlightStatusOnMission                      = 4
	FlightStatusMissionCompleted               = 8
	FlightStatusBackingHome                    = 16
	FlightStatusCreated                        = 256
	FlightStatusActive                         = 512
	FlightStatusCompleted                      = 1024
	FlightStatusProcessed                      = 2048
)

type Flight struct {
	ID          int64     `json:"id" csv:"id" pg:"id,pk,unique"`
	DroneID     int64     `json:"drone_id" csv:"drone_id" pg:"drone_id"`
	Hash        string    `json:"hash" csv:"hash" pg:"hash,unique,notnull"`
	Mission     string    `json:"mission" csv:"mission" pg:"mission,unique,notnull"`
	MissionInfo string    `json:"mission_info" csv:"mission_info" pg:"mission_info"`
	MissionType string    `json:"mission_type" csv:"mission_type" pg:"mission_type"`
	Length      int64     `json:"length" csv:"length" pg:"length"`
	Duration    int64     `json:"duration" csv:"duration" pg:"duration"`
	Distance    float32   `json:"distance" csv:"distance" pg:"distance"`
	Status      int       `json:"status" csv:"status" pg:"status"`
	Date        time.Time `json:"date" csv:"date" pg:"date"`
	LastUpdate  time.Time `json:"last_update" csv:"last_update" pg:"last_update"`
}

func (f *Flight) Select(db *gopg.DB) bool {
	err := db.Model(f).WherePK().Select()
	if err != nil {
		panic(err)
	}

	return f.ID > 0
}

func (f *Flight) ByHash(db *gopg.DB) bool {
	err := db.Model(f).Where("hash = ?", f.Hash).Select()
	if err != nil {
		return false
	}

	return f.ID > 0
}

func (f *Flight) Insert(db *gopg.DB) bool {
	res, err := db.Model(f).OnConflict("DO NOTHING").Insert()
	if err != nil {
		panic(err)
	}

	if res.RowsAffected() > 0 {
		fmt.Printf("Flight[%d] inserted.\n", f.ID)
		return true
	} else {
		err = db.Model(f).Where("hash = ?", f.Hash).Select()
		if err != nil {
			return false
		}

		if f.ID > 0 {
			fmt.Printf("Flight[%d] exits.\n", f.ID)
			return true
		}
	}

	return false
}

func (f *Flight) Update(db *gopg.DB) bool {
	res, err := db.Model(f).WherePK().Update()
	if err != nil {
		panic(err)
	}
	if res.RowsAffected() > 0 {
		fmt.Printf("Flight[%d]  updated.\n", f.ID)
		return true
	}

	return false
}
