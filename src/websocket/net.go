// Package websocket
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-30
package websocket

import (
	"log"
	"net"
)

// GetOutboundIP gets preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
