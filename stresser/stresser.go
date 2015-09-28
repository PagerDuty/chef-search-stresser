// Copyright 2015 PagerDuty, Inc. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package stresser

// DefaultQueries is a default set of queries to run when stress testing.
var DefaultQueries = []string{"roles:*", "name:*"}

// DefaultAttributes is the default set of attributes to pull from the
// ChefServer for stress testing. This is just a random assortment
// of node data.
var DefaultAttributes = map[string]interface{}{
	"name":                            []string{"name"},
	"role":                            []string{"role"},
	"roles":                           []string{"roles"},
	"environment":                     []string{"chef_environment"},
	"ipaddress":                       []string{"ipaddress"},
	"ec2.public_ipv4":                 []string{"ec2", "public_ipv4"},
	"ec2.local_ipv4":                  []string{"ec2", "local_ipv4"},
	"ec2.placement_availability_zone": []string{"ec2", "placement_availability_zone"},
	"cloud.private_ips":               []string{"cloud", "private_ips"},
	"cloud.public_ips":                []string{"cloud", "public_ips"},
	"cloud.provider":                  []string{"cloud", "provider"},
	"cloud.public_ip":                 []string{"cloud", "public_ip"},
	"cloud.private_ip":                []string{"cloud", "private_ip"},
}

// StressResult is a struct that contains the information for a single
// Search result (Partial or Full).
//
// The ID filed is the auto-incremented ID of which goroutine this was,
// this allows you to know whether it was an earlier or later goroutine.
type StressResult struct {
	ID         int
	Duration   int64
	Err        error
	NumResults int
	Query      string
}

// StressResults is the entire resultset for the stressin gof the Chef
// server. The field names should be pretty self-explanatory.
type StressResults struct {
	Fastest *StressResult
	Slowest *StressResult
	Average float64
	Results []*StressResult
}

// ChefConfig is a struct that contains all the items we need to build
// a *chef.Client.
type ChefConfig struct {
	PemFile    string
	ChefServer string
	ClientName string
}

// Stresser is a struct for easily running different stress tests against the
// Chef server. The Queries field is a slice of queries to run against the Chef
// server, with one being picked randomly for each search. The Attributes
// field is the field required if you want to use PartialSearch. Each item
// in the map is usually a []string.
type Stresser struct {
	Queries    []string
	Attributes map[string]interface{}
	Config     ChefConfig
}

// PartialSearch is a function to execute a partial search using the Attributes field. The provided
// parameter is how many concurrent queries you want. This function wraps PartialSearch().
func (s *Stresser) PartialSearch(concurrency int) *StressResults {
	return PartialSearch(s.Config, concurrency, s.Queries, s.Attributes)
}

// Search is a function to execute a full search. The provided parameter is how
// many concurrent queries you want. This function wraps Search().
func (s *Stresser) Search(concurrency int) *StressResults {
	return Search(s.Config, concurrency, s.Queries)
}
