// Copyright 2015 PagerDuty, Inc. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/PagerDuty/chef-search-stresser/stresser"
)

type runResult struct {
	res      *stresser.StressResults
	elapsed  float64
	c        int
	printErr bool
}

func readClientKey(path string) (string, error) {
	// read the entire PEM file in
	key, err := ioutil.ReadFile(path)

	return string(key), err
}

// processFlags proceses the command line arguments.
// if parsing fails, or a special mode is requested (--version),
// this may exit the program
func processFlags() *binArgs {
	args := &binArgs{}
	out, err := args.parse(os.Args)

	if err != nil {
		log.Fatalf(err.Error())
	}

	if out != "" {
		fmt.Print(out)
		os.Exit(0)
	}

	return args
}

func printResults(rnresult *runResult) {
	res := rnresult.res

	tookStr := strconv.FormatFloat(rnresult.elapsed, 'f', -1, 64)
	avgStr := strconv.FormatFloat(res.Average, 'f', -1, 64)

	log.Printf("the %d search(es) took %s seconds to complete", rnresult.c, tookStr)

	if res.Fastest == nil {
		log.Printf("[warn] it's possible that none of the searches returned?")
	} else {
		log.Printf("the fastest took %dms and the slowest took %dms", res.Fastest.Duration, res.Slowest.Duration)
	}

	log.Printf("the average was %sms", avgStr)

	counts := make(map[int]int)

	for _, result := range res.Results {
		if result.Err != nil {
			if rnresult.printErr {
				log.Printf("[error] [%d] %s", result.ID, result.Err)
			}

			// counts[-1] is our error counter
			counts[-1]++
		}

		counts[result.NumResults]++
	}

	for key, value := range counts {
		if key < 0 {
			continue
		}

		log.Printf("%d search(es) had %d results", value, key)
	}

	if counts[-1] > 0 {
		log.Printf("%d search(es) had errors", counts[-1])
	}

	fmt.Print("\n")
}

func main() {
	// process command line flags
	// this funcion can may abort execution
	args := processFlags()

	if !args.DisableSeed {
		// seed the random number generator so when setting multiple queries
		// the order they are ran in appears more random
		rand.Seed(time.Now().UnixNano())
	}

	key, err := readClientKey(args.PemFile)

	if err != nil {
		log.Fatalf(err.Error())
	}

	str := &stresser.Stresser{
		Queries:    args.Queries,
		Attributes: stresser.DefaultAttributes,
		Config: stresser.ChefConfig{
			PemFile:    key,
			ChefServer: args.ChefServer,
			ClientName: args.ClientName,
		},
	}

	log.Printf("starting stresser against %s concurrency %d", args.ChefServer, args.Concurrency)

	result := &runResult{
		c:        int(args.Concurrency),
		printErr: args.PrintErr,
	}

	start := time.Now()

	if args.Partial {
		result.res = str.PartialSearch(int(args.Concurrency))
	} else {
		result.res = str.Search(int(args.Concurrency))
	}

	result.elapsed = time.Since(start).Seconds()

	printResults(result)
}
