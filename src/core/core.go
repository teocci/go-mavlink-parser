// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-27
package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aler9/gomavlib"
	"github.com/aler9/gomavlib/pkg/dialects/ardupilotmega"
	"github.com/teocci/go-mavlink-parser/src/data"
	"github.com/teocci/go-mavlink-parser/src/model"
	"github.com/teocci/go-mavlink-parser/src/wsnet"
)

type InitConf struct {
	Host     string
	Port     int64
	ClientID int64
	DroneID  int64
	FlightID int64
}

var (
	initConf InitConf
	rtt      *data.RTT
	ws       *wsnet.Client
)

func Start(c InitConf) error {
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

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// init ws
	ws = wsnet.NewClient(interrupt)

	var trigger = 0
	for event := range node.Events() {
		if frm, ok := event.(*gomavlib.EventFrame); ok {
			//fmt.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
			if trigger == 0 {
				rtt = &data.RTT{
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
				rtt.Alt = msg.Alt
				rtt.LastUpdate = time.Now()

				trigger |= 2
			}

			if trigger&2 == 2 {
				process(rtt)
				trigger = 0
			}
		}
	}

	return nil
}

func process(rtt *data.RTT) {
	req := &data.ReqUpdate{
		CMD:     wsnet.CMDUpdateTelemetry,
		ToID:    initConf.ClientID,
		DroneID: initConf.DroneID,
		Record:  *rtt,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("process(): %v", err)
	}

	ws.Send <- jsonData
}
