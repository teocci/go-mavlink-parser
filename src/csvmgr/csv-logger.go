// Package csvmgr
// Created by RTT.
// Author: teocci@yandex.com on 2021-Nov-11
package csvmgr

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/teocci/go-mavlink-parser/src/datamgr"
)

const (
	moduleTagTL = "tl"
	moduleTagTR = "tr"
	moduleTagBL = "bl"
	moduleTagBR = "br"
)

const (
	prefixRTTFile        = "rtt"
	formatRTTFilename    = "%s-%d-%s"
	formatRTTFilenameExt = "%s.csv"
	formatRTTDirPath     = "%s/c-%d/d-%d"

	baseLogsPath = "/home/rtt/jinan/logs"
	//baseLogsPath = "./tmp"
)

type CSVLogger struct {
	LogFile   *os.File
	LogWriter *bufio.Writer
	// Buffered channel of outbound messages.
	Append    chan []byte
	Done      chan struct{}
	Interrupt chan os.Signal
}

func (c *CSVLogger) onMavlinkMessage() {
	defer CloseFile()(c.LogFile)

	for {
		select {
		case <-c.Done:
			return
		case buffer, ok := <-c.Append:
			if !ok {
				return
			}

			_, err := c.LogWriter.Write(buffer)
			HasError(err, false)

			// Add queued chat messages to the current websocket buffer.
			n := len(c.Append)
			for i := 0; i < n; i++ {
				_, _ = c.LogWriter.Write(<-c.Append)
			}

			FlushWriter()(c.LogWriter)
		case <-c.Interrupt:
			log.Println("onMavlinkMessage-> interrupt")

			// Close file
			select {
			case <-c.Done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (c *CSVLogger) Close() {
	c.Close()
}

func NewCSVLogger(c datamgr.InitConf) *CSVLogger {
	t := time.Now()
	fn := fmt.Sprintf(formatRTTFilename, prefixRTTFile, c.FlightID, t.Format("20060102-150405"))
	fmt.Println("RTT filename:", fn)

	rttLogPath := fmt.Sprintf(formatRTTDirPath, baseLogsPath, c.CompanyID, c.DroneID)

	fnExt := fmt.Sprintf(formatRTTFilenameExt, fn)
	rttPath := filepath.Join(rttLogPath, fnExt)
	fmt.Println("RTT file path:", rttPath)

	err := os.MkdirAll(filepath.Dir(rttPath), os.ModePerm)
	HasError(err, true)

	file := CreateFile(rttPath)

	writer := bufio.NewWriter(file)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	csvLogger := &CSVLogger{
		LogFile:   file,
		LogWriter: writer,
		Append:    make(chan []byte, 256),
		Done:      make(chan struct{}),
		Interrupt: interrupt,
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go csvLogger.onMavlinkMessage()

	return csvLogger
}
