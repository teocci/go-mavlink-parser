// Package data
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-30
package data

type ReqUpdate struct {
	CMD string `json:"cmd"`
	Record RTT `json:"record"`
}
