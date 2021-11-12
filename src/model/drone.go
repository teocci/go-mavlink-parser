// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2021-Oct-19
package model

import (
	"time"

	gopg "github.com/go-pg/pg/v10"
)

type Drone struct {
	ID               int64     `json:"id" csv:"id" pg:"id,pk,unique"`
	CompanyID        int64     `json:"company_id" csv:"company_id" pg:"company_id,notnull"`
	Code             string    `json:"code" csv:"code" pg:"code,notnull"`
	Name             string    `json:"name" csv:"name" pg:"name"`
	Model            string    `json:"model" csv:"model" pg:"model"`
	TypeInfo         string    `json:"type_info" csv:"type_info" pg:"type_info"`
	Size             string    `json:"size" csv:"size" pg:"size"`
	EngineUnit       string    `json:"engine_unit" csv:"engine_unit" pg:"engine_unit"`
	MaxTakeOffWeight string    `json:"max_take_off_weight" csv:"max_take_off_weight" pg:"max_take_off_weight"`
	OperatingRange   string    `json:"operating_range" csv:"operating_range" pg:"operating_range"`
	MaxOperatingTime string    `json:"max_operating_time" csv:"max_operating_time" pg:"max_operating_time"`
	Restrictions     string    `json:"restrictions" csv:"restrictions" pg:"restrictions"`
	Description      string    `json:"description" csv:"description" pg:"description"`
	TaskInfo         string    `json:"task_info" csv:"task_info" pg:"task_info"`
	Task             int       `json:"task" csv:"task" pg:"task"`
	Autopilot        int       `json:"autopilot" csv:"autopilot" pg:"autopilot"`
	Type             int       `json:"type" csv:"type" pg:"type"`
	BaseMode         int       `json:"base_mode" csv:"base_mode" pg:"base_mode"`
	LastUpdate       time.Time `json:"last_update" csv:"last_update" pg:"last_update"`
}

func (d *Drone) Select(db *gopg.DB) bool {
	err := db.Model(d).WherePK().Select()
	if err != nil {
		panic(err)
	}

	return d.ID > 0
}

func (d *Drone) ByCompanyID(db *gopg.DB) bool {
	err := db.Model(d).Where("company_id = ?", d.CompanyID).Limit(1).Select()
	if err != nil {
		panic(err)
	}

	return d.ID > 0
}