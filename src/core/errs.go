// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-27
package core

import (
	"errors"
)

const (
	errInitDataIsNil = "initialization data is nil"
)

func ErrInitDataIsNil()  error {
	return errors.New(errInitDataIsNil)
}
