// Package websocket
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-30
package websocket

import (
	"fmt"
	"testing"
)

func TestGetOutboundIP(t *testing.T) {
	fmt.Println(GetOutboundIP())
}
