// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-27
package core

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aler9/gomavlib"
	"github.com/aler9/gomavlib/pkg/dialects/ardupilotmega"
	"github.com/teocci/go-mavlink-parser/src/csvmgr"
	"github.com/teocci/go-mavlink-parser/src/datamgr"
	"github.com/teocci/go-mavlink-parser/src/model"
	"github.com/teocci/go-mavlink-parser/src/wsnet"
)

var (
	initConf datamgr.InitConf
	rtt      *datamgr.RTT
	ws       *wsnet.Client
	csvl     *csvmgr.CSVLogger

	headerSent = false
)

func Start(c datamgr.InitConf) error {
	initConf = c
	address := fmt.Sprintf("%s:%d", initConf.Host, initConf.Port)
	// create a node which
	// - communicates with a TCP endpoint in client mode
	// - understands ardupilotmega dialect
	// - writes messages with given system id
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointTCPClient{Address: address},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// init db
	db = model.Setup()
	defer db.Close()

	drone := &model.Drone{ID: initConf.DroneID}
	ok := drone.Select(db)
	if ok {
		initConf.CompanyID = drone.CompanyID
	}

	// init ws
	ws = wsnet.NewClient(initConf)

	// init csvlogger
	csvl = csvmgr.NewCSVLogger(initConf)

	var seq int64 = 0
	var trigger = 0
	for event := range node.Events() {
		if frm, ok := event.(*gomavlib.EventFrame); ok {
			//fmt.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
			if trigger == 0 {
				rtt = &datamgr.RTT{
					DroneID:  initConf.DroneID,
					FlightID: initConf.FlightID,
				}
			}

			switch msg := frm.Message().(type) {
			case *ardupilotmega.MessageHeartbeat:
				//fmt.Printf("received heartbeat (type %d)\n", msg.Type)
			case *ardupilotmega.MessageAttitude:
				rtt.Roll = msg.Roll
				rtt.Pitch = msg.Pitch
				rtt.Yaw = msg.Yaw

				trigger |= 1
			case *ardupilotmega.MessageGlobalPositionInt:
				rtt.TimeBootMs = msg.TimeBootMs
				rtt.Lat = msg.Lat
				rtt.Lon = msg.Lon
				rtt.Alt = msg.Alt // convert to meters
				rtt.LastUpdate = time.Now()

				trigger |= 2
			}

			if trigger&2 == 2 {
				rtt.Seq = seq
				process(rtt)
				trigger = 0
				seq++
			}
		}
	}

	return nil
}

func process(rtt *datamgr.RTT) {
	req := &datamgr.UpdateTelemetry{
		CMD:       wsnet.CMDUpdateTelemetry,
		ToConnID:  initConf.ConnID,
		ModuleTag: initConf.ModuleTag,
		DroneID:   initConf.DroneID,
		Record:    *rtt,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Printf("process(): %v", err)
	}

	ws.Send <- jsonData

	appendRecord(rtt)
	insertRecord(rtt)
}
