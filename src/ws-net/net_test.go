// Package ws_net
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-30
package ws_net

import (
	"fmt"
	"testing"
)

func TestGetOutboundIP(t *testing.T) {
	fmt.Println(GetOutboundIP())
}
