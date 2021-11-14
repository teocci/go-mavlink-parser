// Package core
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Nov-14
package core

import (
	"bytes"
	"fmt"
	"github.com/jszwec/csvutil"
	"github.com/teocci/go-mavlink-parser/src/datamgr"
	"io"
	"log"
)

func appendRecord(rtt *datamgr.RTT) {
	rtts := []datamgr.RTT{*rtt}
	b, err := csvutil.Marshal(rtts)
	if err != nil {
		log.Println("error:", err)
	}

	buf := bytes.NewBuffer(b)

	header, err := buf.ReadBytes('\n')
	if err != nil && err != io.EOF {
		log.Println("error:", err)
	}

	line, err := buf.ReadBytes('\n')
	if err != nil && err != io.EOF {
		log.Println("error:", err)
	}

	h := string(header)
	fmt.Println(h)

	s := string(line)
	fmt.Println(s)

	if !headerSent {
		csvl.Append <- header
		headerSent = true
	}
	csvl.Append <- line
}
