// Package cmd
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-27
package cmd

import (
	"errors"
	"fmt"
)

const (
	errFileDoesNotExist = "%s file does not exist"
	errCanNotLoadLogger = "cannot load logger -> %s"
	errInitDataIsNil = "initialization data is nil"
)

func ErrInitDataIsNil()  error {
	return errors.New(errInitDataIsNil)
}

func ErrFileDoesNotExist(f string) error {
	return errors.New(fmt.Sprintf(errFileDoesNotExist, f))
}

func ErrCanNotLoadLogger(e error) error {
	return errors.New(fmt.Sprintf(errCanNotLoadLogger, e))
}
