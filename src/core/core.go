// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-27
package core

import (
	"fmt"
	"github.com/aler9/gomavlib"
	"github.com/aler9/gomavlib/pkg/dialects/ardupilotmega"
)

func Start() {
	// create a node which
	// - communicates with a TCP endpoint in client mode
	// - understands ardupilotmega dialect
	// - writes messages with given system id
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointTCPClient{Address: "106.244.179.242:20102"},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print every message we receive
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			fmt.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
		}
	}
}
