// Copyright 2015 PagerDuty, Inc. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/PagerDuty/chef-search-stresser/stresser"
	"github.com/jessevdk/go-flags"
)

type binArgs struct {
	Queries     []string `short:"q" long:"query" description:"set some queries to run for testing. Defaults to ['roles:*' 'name:*']. One flag for each query"`
	Concurrency uint     `short:"c" long:"concurrency" default:"1" description:"set the number of concurrent searches to run"`
	Partial     bool     `short:"p" long:"partial" default:"false" description:"run partial search"`
	Full        bool     `short:"f" long:"full" default:"false" description:"run full search"`
	PemFile     string   `short:"k" long:"client-key" description:"path to the client.pem file" default:"/etc/chef/client.pem"`
	ChefServer  string   `short:"S" long:"server" description:"specify which Chef server to use"`
	ClientName  string   `short:"N" long:"name" description:"name of the client being used"`
	PrintErr    bool     `short:"E" long:"print-errors" default:"false" description:"print out errors that occured"`
	DisableSeed bool     `long:"disable-seed" default:"false" description:"by default, we seed the random number generator to help with better randomization of the queries ran"`
}

func (a *binArgs) parse(args []string) (string, error) {
	if args == nil {
		args = os.Args
	}

	p := flags.NewParser(a, flags.HelpFlag|flags.PassDoubleDash)

	_, err := p.ParseArgs(args[1:])

	// determine if there was a parsing error
	// unfortunately, help message is returned as an error
	if err != nil {
		// determine whether this was a help message by doing a type
		// assertion of err to *flags.Error and check the error type
		// if it was a help message, do not return an error
		if errType, ok := err.(*flags.Error); ok {
			if errType.Type == flags.ErrHelp {
				return err.Error(), nil
			}
		}

		return "", err
	}

	if a.Partial && a.Full {
		return "", fmt.Errorf("unable to run both full and partial search tests at the same time")
	}

	if !a.Partial && !a.Full {
		return "", fmt.Errorf("need to choose one mode of operation: --partial or --full; see --help for more details")
	}

	if a.ClientName == "" {
		curUser, err := user.Current()

		if err != nil {
			return "", fmt.Errorf("error getting user: %s", err.Error())
		}

		a.ClientName = curUser.Username
	}

	if len(a.Queries) == 0 {
		a.Queries = stresser.DefaultQueries
	}

	return "", err
}
