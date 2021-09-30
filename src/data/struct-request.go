// Package data
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-30
package data

type Register struct {
	CMD       string `json:"cmd"`
	ConnID    int64  `json:"connection_id"`
	ModuleTag string `json:"module_tag"`
	DroneID   int64  `json:"drone_id"`
}

type UpdateTelemetry struct {
	CMD       string `json:"cmd"`
	ToConnID  int64  `json:"to_connection_id"`
	ModuleTag string `json:"module_tag"`
	DroneID   int64  `json:"drone_id"`
	Record    RTT    `json:"record"`
}
