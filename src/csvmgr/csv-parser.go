// Package csvmgr
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-30
package csvmgr

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/teocci/go-mavlink-parser/src/jobmgr"
	"github.com/teocci/go-mavlink-parser/src/utfmgr"
)

const (
	headerNormalizer = "fcc_time,gps_time,str_01,num_01,lat,lon,alt,roll,pitch"
	regexNormalizer  = `(?P<prefix>.*[^,{2,}])(?P<sufix>,*$)`
)

func explodeHeader() []string {
	return explodeValues(headerNormalizer)
}

func explodeValues(str string) []string {
	return strings.Split(str, ",")
}

func explodeKeyValue(msg string) (dataSlice map[string]string) {
	dataSlice = map[string]string{}

	slice := explodeValues(msg)
	length := len(slice)

	sliceHeader := explodeHeader()
	headerLength := len(sliceHeader)

	if headerLength == length {
		for i := 0; i < length; i++ {
			dataSlice[sliceHeader[i]] = slice[i]
		}
	}

	return dataSlice
}

func NormalizeJob(fn string) (buffer bytes.Buffer) {
	// Add header
	buffer.WriteString(headerNormalizer)

	poolNumber := runtime.NumCPU()
	dispatcher := jobmgr.NewDispatcher(poolNumber).Start(func(id int, job jobmgr.Job) error {
		str := job.Item.(string)
		str = normalizer(str)

		if len(str) > 0 {
			fields := strings.Split(str, ",")
			if len(fields) == 9 {
				buffer.WriteString(str)

				fmt.Printf("%s", str)
			}
		}

		return nil
	})

	f := OpenUTFFile(fn)
	defer utfmgr.CloseFile()(f)

	var index = 0
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		str := fileScanner.Text()
		err := fileScanner.Err()
		if err != nil {
			log.Fatal(err)
		}

		dispatcher.Submit(jobmgr.Job{
			ID:   index,
			Item: str,
		})

		index++
	}

	return buffer
}

func normalizer(job string) string {
	var re = regexp.MustCompile(regexNormalizer)
	if re.MatchString(job) {
		matches := re.FindStringSubmatch(job)
		idIndex := re.SubexpIndex("prefix")

		return matches[idIndex]
	}

	return ""
}


func Normalize(f *os.File) bytes.Buffer {
	fileScanner := bufio.NewScanner(f)
	var buffer bytes.Buffer

	numWorkers := runtime.NumCPU()
	jobs := make(chan string, numWorkers)
	res := make(chan string)

	var wg sync.WaitGroup
	worker := func(jobs <-chan string, results chan<- string) {
		for {
			select {
			case job, ok := <-jobs: // you must check for readable state of the channel.
				if !ok {
					return
				}
				s := normalizer(job)
				if len(s) > 0 {
					results <- normalizer(job)
				}
			}
		}
	}

	// init workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			// this line will exec when chan `res` processed output at line 107 (func worker: line 71)
			defer wg.Done()
			worker(jobs, res)
		}()
	}

	go func() {
		for fileScanner.Scan() {
			str := fileScanner.Text()
			err := fileScanner.Err()
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				break
			}
			jobs <- str
		}
		close(jobs) // close jobs to signal workers that no more job are incoming.
	}()

	go func() {
		wg.Wait()
		close(res) // when you close(res) it breaks the below loop.
	}()

	for r := range res {
		buffer.WriteString(r)
	}

	return buffer
}