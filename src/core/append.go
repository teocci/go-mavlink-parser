// Package core
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Nov-14
package core

import (
	"bytes"
	"github.com/jszwec/csvutil"
	"github.com/teocci/go-mavlink-parser/src/datamgr"
	"io"
	"log"
)

func appendRecord(record *datamgr.RTT) {
	recordBundle := []datamgr.RTT{*record}
	b, err := csvutil.Marshal(recordBundle)
	HasError(err, false)

	buf := bytes.NewBuffer(b)

	header, err := buf.ReadBytes('\n')
	if err != nil && err != io.EOF {
		log.Println("error:", err)
	}

	line, err := buf.ReadBytes('\n')
	if err != nil && err != io.EOF {
		log.Println("error:", err)
	}

	if !headerSent {
		csvl.Append <- header
		headerSent = true
	}
	csvl.Append <- line
}
