// Copyright 2015 PagerDuty, Inc. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package stresser

import "github.com/go-chef/chef"

// getChef builds a new chef client
func getChef(config ChefConfig) (*chef.Client, error) {
	return chef.NewClient(&chef.Config{
		Name:    config.ClientName,
		Key:     config.PemFile,
		BaseURL: config.ChefServer,
	})
}

// parseChannel parses the channel we use for the search results
// to build an instance of *StressResults. This does things like
// determine the slowest and fastest goroutine as well as the average.
func parseChannel(ch <-chan *StressResult) *StressResults {
	var totalTime int64
	sr := &StressResults{}

	// by the time we get the channel it should be closed and full
	// of data -- set n to the number of results
	count := len(ch)

	// loop over the closed channel
	for res := range ch {
		// if it's a nil result
		// subtract from the count and continue
		if res == nil {
			count--
			continue
		}

		// append the *StressResult
		sr.Results = append(sr.Results, res)

		// if there was an error
		// remove it from the count as we are not
		// going to include it in the average calculation
		if res.Err != nil {
			count--
			continue
		}

		// if Fastest is nil, then Slowest must be nil too
		// so let's set both
		if sr.Fastest == nil {
			sr.Fastest = res
			sr.Slowest = res
		} else {
			// if this search was faster than the current fastest
			// update the Fastest field
			if res.Duration < sr.Fastest.Duration {
				sr.Fastest = res
			}

			// if this search was slower than the current slowest
			// update the Slowest field
			if res.Duration > sr.Slowest.Duration {
				sr.Slowest = res
			}
		}

		// add the duration to totalTime so we can calulate the average
		totalTime += res.Duration
	}

	sr.Average = float64(totalTime) / float64(count)

	return sr
}

// fromNsToMs is a function to convert Nanoseconds to Milliseconds
func fromNsToMs(nanoseconds int64) int64 {
	return nanoseconds / 1000000
}
