// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-27
package core

import (
	"fmt"
	"github.com/teocci/go-mavlink-parser/src/model"
	"time"

	"github.com/aler9/gomavlib"
	"github.com/aler9/gomavlib/pkg/dialects/ardupilotmega"
	"github.com/teocci/go-mavlink-parser/src/data"
)

type InitData struct {
	Host     string
	Port     int64
	DroneID  int64
	FlightID int64
}

var rtt *data.RTT

func Start(initData InitData) error {
	address := fmt.Sprintf("%s:%d", initData.Host, initData.Port)
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

	var trigger = 0
	for event := range node.Events() {
		if frm, ok := event.(*gomavlib.EventFrame); ok {
			fmt.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())

			if trigger == 0 {
				rtt = &data.RTT{
					DroneID:  initData.DroneID,
					FlightID: initData.FlightID,
				}
			}

			switch msg := frm.Message().(type) {
			case *ardupilotmega.MessageHeartbeat:
				fmt.Printf("received heartbeat (type %d)\n", msg.Type)
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

}
