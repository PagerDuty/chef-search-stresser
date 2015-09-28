// Copyright 2015 PagerDuty, Inc. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package stresser

import (
	"errors"

	. "gopkg.in/check.v1"
)

func (t *TestSuite) Test_parseChannel(c *C) {
	var res *StressResults

	ch := make(chan *StressResult, 4)
	tErr := errors.New("TEST ERROR")

	ch <- &StressResult{
		ID:         0,
		Duration:   1000,
		Err:        nil,
		NumResults: 10,
		Query:      "false_attribute:*",
	}

	ch <- &StressResult{
		ID:         1,
		Duration:   42,
		Err:        tErr,
		NumResults: 0,
		Query:      "recipes:bad::syntax",
	}

	ch <- nil

	ch <- &StressResult{
		ID:         2,
		Duration:   2400,
		Err:        nil,
		NumResults: 84,
		Query:      "roles:*",
	}

	close(ch)

	res = parseChannel(ch)
	c.Assert(res, Not(IsNil))
	c.Assert(res.Fastest, Not(IsNil))
	c.Assert(res.Slowest, Not(IsNil))
	c.Assert(len(res.Results), Equals, 3)

	c.Check(res.Average, Equals, 1700.0)

	c.Check(res.Fastest.ID, Equals, 0)
	c.Check(res.Fastest.Duration, Equals, int64(1000))
	c.Check(res.Fastest.Err, IsNil)
	c.Check(res.Fastest.NumResults, Equals, 10)
	c.Check(res.Fastest.Query, Equals, "false_attribute:*")

	c.Check(res.Slowest.ID, Equals, 2)
	c.Check(res.Slowest.Duration, Equals, int64(2400))
	c.Check(res.Slowest.Err, IsNil)
	c.Check(res.Slowest.NumResults, Equals, 84)
	c.Check(res.Slowest.Query, Equals, "roles:*")

	c.Check(res.Results[0].ID, Equals, 0)
	c.Check(res.Results[0].Duration, Equals, int64(1000))
	c.Check(res.Results[0].Err, IsNil)
	c.Check(res.Results[0].NumResults, Equals, 10)
	c.Check(res.Results[0].Query, Equals, "false_attribute:*")

	c.Check(res.Results[1].ID, Equals, 1)
	c.Check(res.Results[1].Duration, Equals, int64(42))
	c.Check(res.Results[1].Err, Equals, tErr)
	c.Check(res.Results[1].NumResults, Equals, 0)
	c.Check(res.Results[1].Query, Equals, "recipes:bad::syntax")

	c.Check(res.Results[2].ID, Equals, 2)
	c.Check(res.Results[2].Duration, Equals, int64(2400))
	c.Check(res.Results[2].Err, IsNil)
	c.Check(res.Results[2].NumResults, Equals, 84)
	c.Check(res.Results[2].Query, Equals, "roles:*")
}

func (t *TestSuite) Test_fromNsToMs(c *C) {
	var ms int64
	ms = fromNsToMs(20400000)
	c.Check(ms, Equals, int64(20))
}
