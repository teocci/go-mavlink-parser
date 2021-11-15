// Package wsnet
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-30
package wsnet

import (
	"fmt"
	"log"
	"net"
)

var (
	serverIP      string
	localIP       net.IP
	localAddress  string
	remoteAddress string
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

func WebsocketServerIP(ip string) {
	serverIP = ip

	localIP = net.ParseIP(serverIP)

	localAddress = fmt.Sprintf("localhost:%d", wsPort)
	remoteAddress = fmt.Sprintf("%s:%d", serverIP, wsPort)
}
