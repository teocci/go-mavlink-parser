// Package main
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-18
package main

import (
	"fmt"
	"github.com/teocci/go-mavlink-parser/src/jobmgr"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type Item struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var terms = []int{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}

var c *http.Client

func Benchmark(b *testing.B) {
	c = &http.Client{Timeout: time.Millisecond * 5000}

	fixture := []struct {
		desc string
		pool int
	}{
		{
			desc: "1 worker",
			pool: 1,
		},
		{
			desc: "2 workers",
			pool: 2,
		},
		{
			desc: "4 workers",
			pool: 4,
		},
		{
			desc: "8 workers",
			pool: 8,
		},
		{
			desc: "16 workers",
			pool: 16,
		},
		{
			desc: "32 workers",
			pool: 32,
		},
	}

	tests := []struct {
		desc string
		fn   func(*testing.B, *jobmgr.Dispatcher)
	}{
		{
			desc: "Concurrent",
			fn:   concurrent,
		},
	}

	for _, t := range tests {
		b.Run(t.desc, func(b *testing.B) {
			for _, f := range fixture {
				b.Run(f.desc, func(b *testing.B) {
					dd := jobmgr.NewDispatcher(f.pool).Start(callApi) // start up worker pool
					t.fn(b, dd)
				})
			}
		})
	}
}

func concurrent(b *testing.B, dd *jobmgr.Dispatcher) {
	for n := 0; n < b.N; n++ {
		for i := range terms {
			dd.Submit(jobmgr.Job{
				ID: i,
				Item: Item{
					Name:      fmt.Sprintf("JobID::%d", i),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			})
		}
	}
}

func callApi(id int, job jobmgr.Job) error {
	baseURL := "https://age-of-empires-2-api.herokuapp.com/api/v1/civilization/%d"

	ur := fmt.Sprintf(baseURL, job.ID)
	req, err := http.NewRequest(http.MethodGet, ur, nil)
	if err != nil {
		//log.Printf("error creating a request for term %d :: error is %+v", num, err)
		return err
	}
	res, err := c.Do(req)
	if err != nil {
		//log.Printf("error querying for term %d :: error is %+v", num, err)
		return err
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		//log.Printf("error reading response body :: error is %+v", err)
		return err
	}

	// log.Printf("%d  :: ok", id)
	return nil
}
