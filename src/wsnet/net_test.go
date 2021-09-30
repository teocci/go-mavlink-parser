// Package wsnet
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-30
package wsnet

import (
	"fmt"
	"testing"
)

func TestGetOutboundIP(t *testing.T) {
	fmt.Println(GetOutboundIP())
}
