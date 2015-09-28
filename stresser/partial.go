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

// PartialSearch is a function to execute a partial search using the attributes param. The provided
// concurrency is how many concurrent queries you want. The queries parameter is a slice of possible
// queries to run. The attributes parameter is are the attributes for the partial search. Each entry
// in the map is usually a []string. This function does not loop and only does  one invocation of n
// queries.
//
// If multiple queries are provided a random query is used per invocation.
func PartialSearch(config ChefConfig, concurrency int, queries []string, attributes map[string]interface{}) *StressResults {
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

		go partialSearch(
			client,
			queries[rand.Intn(len(queries))],
			i, attributes, ch, wg,
		)
	}

	// once all goroutines have returned
	// close the channel and parse results
	wg.Wait()
	close(ch)

	return parseChannel(ch)
}

func partialSearch(c *chef.Client, query string, id int, attributes map[string]interface{}, ch chan<- *StressResult, wg *sync.WaitGroup) {
	start := time.Now()
	res, err := c.Search.PartialExec("node", query, attributes)
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
