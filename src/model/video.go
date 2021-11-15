// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2021-Oct-22
package model
import (
	"fmt"
	"log"
	"time"

	gopg "github.com/go-pg/pg/v10"
)

type Video struct {
	ID          int64     `json:"id" csv:"id" pg:"id,pk,unique"`
	FlightID    int64     `json:"flight_id" csv:"flight_id" pg:"flight_id,notnull"`
	Name        string    `json:"name" csv:"name" pg:"name,unique,notnull"`
	Path        string    `json:"path" csv:"path" pg:"path,unique,notnull"`
	Status      int       `json:"status" csv:"status" pg:"status"`
	LastUpdate  time.Time `json:"last_update" csv:"last_update" pg:"last_update"`
}

func (v *Video) Insert(db *gopg.DB) bool {
	res, err := db.Model(v).OnConflict("DO NOTHING").Insert()
	if err != nil {
		log.Println(err)
		return false
	}

	if res.RowsAffected() > 0 {
		fmt.Printf("Video[%d] inserted.\n", v.ID)
		return true
	} else {
		err = db.Model(v).Where("flight_id = ?", v.FlightID).Select()
		if err != nil {
			return false
		}

		if v.ID > 0 {
			fmt.Printf("Video[%d] exits.\n", v.ID)
			return true
		}
	}

	return false
}

func (v *Video) Update(db *gopg.DB) bool {
	res, err := db.Model(v).WherePK().Update()
	if err != nil {
		panic(err)
	}
	if res.RowsAffected() > 0 {
		fmt.Printf("Video[%d]  updated.\n", v.ID)
		return true
	}

	return false
}