// Copyright 2015 PagerDuty, Inc. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package main

import (
	"os/user"
	"testing"

	. "gopkg.in/check.v1"
)

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (t *TestSuite) TestbinArgs_parse(c *C) {
	const arg0 = "/usr/local/bin/chef-search-stresser"

	var output string
	var err error

	//
	// test that the short-flags work
	//
	args := &binArgs{}
	cli := []string{
		arg0,
		"-q", "role:*",
		"-q", "thing:*",
		"-c", "42",
		"-p",
		"-k", "/tmp/file",
		"-S", "http://localhost:4000/",
		"-N", "fauxhai",
		"-E",
	}

	output, err = args.parse(cli)
	c.Assert(err, IsNil)
	c.Check(output, Equals, "")
	c.Assert(len(args.Queries), Equals, 2)
	c.Check(args.Queries[0], Equals, "role:*")
	c.Check(args.Queries[1], Equals, "thing:*")
	c.Check(args.Concurrency, Equals, uint(42))
	c.Check(args.Partial, Equals, true)
	c.Check(args.Full, Equals, false)
	c.Check(args.PemFile, Equals, "/tmp/file")
	c.Check(args.ChefServer, Equals, "http://localhost:4000/")
	c.Check(args.ClientName, Equals, "fauxhai")
	c.Check(args.PrintErr, Equals, true)
	c.Check(args.DisableSeed, Equals, false)

	args = &binArgs{}
	cli = []string{
		arg0,
		"-q", "role:*",
		"-q", "thing:*",
		"-c", "42",
		"-f",
		"-k", "/tmp/file",
		"-S", "http://localhost:4000/",
		"-N", "fauxhai",
		"-E",
	}

	output, err = args.parse(cli)
	c.Assert(err, IsNil)
	c.Check(output, Equals, "")
	c.Assert(len(args.Queries), Equals, 2)
	c.Check(args.Queries[0], Equals, "role:*")
	c.Check(args.Queries[1], Equals, "thing:*")
	c.Check(args.Concurrency, Equals, uint(42))
	c.Check(args.Partial, Equals, false)
	c.Check(args.Full, Equals, true)
	c.Check(args.PemFile, Equals, "/tmp/file")
	c.Check(args.ChefServer, Equals, "http://localhost:4000/")
	c.Check(args.ClientName, Equals, "fauxhai")
	c.Check(args.PrintErr, Equals, true)
	c.Check(args.DisableSeed, Equals, false)

	//
	// test that long flags work
	//
	args = &binArgs{}
	cli = []string{
		arg0,
		"--query", "role:*",
		"--query", "thing:*",
		"--concurrency", "42",
		"--partial",
		"--client-key", "/tmp/file",
		"--server", "http://localhost:4000/",
		"--name", "fauxhai",
		"--print-errors",
		"--disable-seed",
	}

	output, err = args.parse(cli)
	c.Assert(err, IsNil)
	c.Check(output, Equals, "")
	c.Assert(len(args.Queries), Equals, 2)
	c.Check(args.Queries[0], Equals, "role:*")
	c.Check(args.Queries[1], Equals, "thing:*")
	c.Check(args.Concurrency, Equals, uint(42))
	c.Check(args.Partial, Equals, true)
	c.Check(args.Full, Equals, false)
	c.Check(args.PemFile, Equals, "/tmp/file")
	c.Check(args.ChefServer, Equals, "http://localhost:4000/")
	c.Check(args.ClientName, Equals, "fauxhai")
	c.Check(args.PrintErr, Equals, true)
	c.Check(args.DisableSeed, Equals, true)

	args = &binArgs{}
	cli = []string{
		arg0,
		"--query", "role:*",
		"--query", "thing:*",
		"--concurrency", "42",
		"--full",
		"--client-key", "/tmp/file",
		"--server", "http://localhost:4000/",
		"--name", "fauxhai",
		"--print-errors",
		"--disable-seed",
	}

	output, err = args.parse(cli)
	c.Assert(err, IsNil)
	c.Check(output, Equals, "")
	c.Assert(len(args.Queries), Equals, 2)
	c.Check(args.Queries[0], Equals, "role:*")
	c.Check(args.Queries[1], Equals, "thing:*")
	c.Check(args.Concurrency, Equals, uint(42))
	c.Check(args.Partial, Equals, false)
	c.Check(args.Full, Equals, true)
	c.Check(args.PemFile, Equals, "/tmp/file")
	c.Check(args.ChefServer, Equals, "http://localhost:4000/")
	c.Check(args.ClientName, Equals, "fauxhai")
	c.Check(args.PrintErr, Equals, true)
	c.Check(args.DisableSeed, Equals, true)

	//
	// test that the Queries field is automatically populated
	//
	args = &binArgs{}
	cli = []string{
		arg0,
		"--full",
	}

	output, err = args.parse(cli)
	c.Assert(err, IsNil)
	c.Assert(len(args.Queries), Equals, 2)
	c.Check(args.Queries[0], Equals, "roles:*")
	c.Check(args.Queries[1], Equals, "name:*")

	//
	// test that the ClientName field is automatically populated
	//
	args = &binArgs{}
	cli = []string{
		arg0,
		"--full",
	}

	output, err = args.parse(cli)
	c.Assert(err, IsNil)

	curUser, err := user.Current()
	c.Assert(err, IsNil)
	c.Check(args.ClientName, Equals, curUser.Username)

	//
	// test that setting both --full and --partiail causes an error
	//
	args = &binArgs{}
	cli = []string{
		arg0,
		"--full",
		"--partial",
	}

	output, err = args.parse(cli)
	c.Assert(err, Not(IsNil))
	c.Check(err.Error(), Equals, "unable to run both full and partial search tests at the same time")

	//
	// test that setting neither --full and --partiai causes an error
	//
	args = &binArgs{}
	cli = []string{arg0}

	output, err = args.parse(cli)
	c.Assert(err, Not(IsNil))
	c.Check(err.Error(), Equals, "need to choose one mode of operation: --partial or --full; see --help for more details")
}
