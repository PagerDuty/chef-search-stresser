# chef-search-stresser
[![TravisCI Build Status](https://img.shields.io/travis/PagerDuty/chef-search-stresser/master.svg?style=flat)](https://travis-ci.org/PagerDuty/chef-search-stresser)

`chef-search-stresser` is a utility used to stress the Chef server with a high
number of concurrent searches. It's especially helpful when you have searches
you want to test that return a large number of results (i.e., all registered nodes).

The primary bits of `chef-search-stresser` is the `stresser` package. It has all
the functions for running searches for stress testing.

## Installation

```
go get -u github.com/PagerDuty/chef-search-stresser

# if you want to use the binary and not just the package
go install github.com/PagerDuty/chef-search-stresser
```

## Usage

This README primarily goes over the command. If you're interested in using the
stress testing code, you should visit the ![GoDoc page]() for docs on the
`stresser` package.

```
Usage:
  chef-search-stresser [OPTIONS]

Application Options:
  -q, --query=         set some queries to run for testing. Defaults to ['roles:*' 'name:*']. One flag for each query
  -c, --concurrency=   set the number of concurrent searches to run (1)
  -p, --partial        run partial search (false)
  -f, --full           run full search (false)
  -k, --client_key=    path to the client.pem file (/etc/chef/client.pem)
  -S, --server=        specify which Chef server to use
  -N, --name=          name of the client being used
  -E, --print-errors   print out errors that occured (false)
      --disalble-seed  by default, we seed the random number generator to help with better randomization of the queries ran (false)

Help Options:
  -h, --help           Show this help message
```

Here is a simple example of using the command:

```
$ /chef-search-stresser --concurrency 10 --partial --server https://chefserver.local/ --client-key ~/.chef/knife/tim-test.pem --query 'name:*'
2015/09/27 18:35:51 starting stresser against https://chefserver.local/ concurrency 10
2015/09/27 18:35:57 the 10 search(es) took 6.252719311 seconds to complete
2015/09/27 18:35:57 the fastest took 6006ms and the slowest took 6248ms
2015/09/27 18:35:57 the average was 6119ms
2015/09/27 18:35:57 10 search(es) had 396 results
```

## License

This project is released under the BSD 3-Clause License. See the `LICENSE` file
for the full contents of the license.
