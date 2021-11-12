// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2021-Nov-11
package model

import (
	gopg "github.com/go-pg/pg/v10"
	"time"
)

type Company struct {
	ID         int64     `json:"id" csv:"id" pg:"id,pk,unique"`
	Code       string    `json:"code" csv:"code" pg:"code,notnull"`
	Name       string    `json:"name" csv:"name" pg:"name,notnull"`
	LastUpdate time.Time `json:"last_update" csv:"last_update" pg:"last_update,notnull"`
}

func (c *Company) Select(db *gopg.DB) bool {
	err := db.Model(c).WherePK().Select()
	if err != nil {
		panic(err)
	}

	return c.ID > 0
}

func (c *Company) ByCode(db *gopg.DB) bool {
	err := db.Model(c).Where("code = ?", c.Code).Limit(1).Select()
	if err != nil {
		panic(err)
	}

	return c.ID > 0
}
