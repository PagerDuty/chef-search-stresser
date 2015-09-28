// Copyright 2015 PagerDuty, Inc. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package stresser

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/go-chef/chef"
)

// Search is a function to execute a full search where all node attributes are returned. The
// provided concurrency param is how many concurrent queries you want. The queries parameter is a
// slice of possible queries to run. This function does not loop and only does  one invocation of n
// queries.
//
// If multiple queries are provided a random query is used per invocation.
func Search(config ChefConfig, concurrency int, queries []string) *StressResults {
	wg := new(sync.WaitGroup)
	ch := make(chan *StressResult, concurrency)

	// add the number of goroutines we'll be spinning off to the WaitGroup
	wg.Add(concurrency)

	// loop n times and create a new search goroutine each time
	for i := 0; i < concurrency; i++ {
		// a new Chef client is created to avoid any sort
		// of mutex within the client causing issues running
		// searches in parallel
		client, err := getChef(config)

		if err != nil {
			log.Fatalf("failed to build chef client: %s", err.Error())
		}

		// spin off a search goroutine
		go search(
			client,
			queries[rand.Intn(len(queries))], // pick a random query
			i, ch, wg,
		)
	}

	// once all goroutines have returned
	// close the channel and parse results
	wg.Wait()
	close(ch)

	return parseChannel(ch)
}

func search(c *chef.Client, query string, id int, ch chan<- *StressResult, wg *sync.WaitGroup) {
	start := time.Now()
	res, err := c.Search.Exec("node", query)
	dur := fromNsToMs(time.Since(start).Nanoseconds())

	ch <- &StressResult{
		ID:         id,
		Duration:   dur,
		Err:        err,
		NumResults: len(res.Rows),
		Query:      query,
	}

	wg.Done()
}
