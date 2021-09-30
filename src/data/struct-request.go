// Package data
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-30
package data

type ReqUpdate struct {
	CMD     string `json:"cmd"`
	ToID    int64  `json:"to_id"`
	DroneID int64  `json:"drone_id"`
	Record  RTT    `json:"record"`
}
