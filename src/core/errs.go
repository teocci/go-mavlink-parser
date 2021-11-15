// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-27
package core

import (
	"errors"
	"log"
)

const (
	errInitDataIsNil = "initialization datamgr is nil"
)

func ErrInitDataIsNil()  error {
	return errors.New(errInitDataIsNil)
}

func HasError(e error, fatal bool, format ...string) {
	hasFormat := false
	if len(format) > 0 {
		hasFormat = true
	}
	if e != nil && !fatal && !hasFormat {
		log.Println(e)
		return
	}

	if e != nil && !fatal && hasFormat {
		log.Printf(format[0], e)
		return
	}

	if e != nil && fatal {
		log.Fatal(e)
	}
}
