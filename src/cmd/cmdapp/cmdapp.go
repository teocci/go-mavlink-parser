// Package cmdapp
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-27
package cmdapp

const (
	Name  = "go-mavlink-parser"
	Short = "Simple implementation to parse a tcp Mavlink protocol"
	Long  = `This application is an open-source tool that can power UGVs, UAVs, ground stations, monitoring systems or routers, connected to other Mavlink-capable devices through a serial port, UDP, TCP or custom transport.`

	HName  = "hostname"
	HShort = "H"
	HDesc  = "Hostname to connect."

	PName    = "port"
	PShort   = "p"
	PDesc    = "Port to connect"
	PDefault = 20102

	CName    = "connection-id"
	CShort   = "c"
	CDesc    = "Connection id"

	MName    = "module-tag"
	MShort   = "m"
	MDesc    = "Module tag"

	DName    = "drone-id"
	DShort   = "d"
	DDesc    = "Drone id"

	FName    = "flight-id"
	FShort   = "f"
	FDesc    = "Flight id"

	WName    = "websocket-host"
	WShort   = "w"
	WDesc    = "Hostname of the websocket"

	UName    = "use-tmp"
	UShort   = "u"
	UDesc    = "Use tmp directory to save the logs"
	UDefault = false
)

const (
	VersionTemplate = "%s %s.%s\n"
	Version         = "v1.0"
	Commit          = "0"
)
