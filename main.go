package main

import (
	"flag"
	"strings"

	"github.com/crypt-data/coinmarket-scraper/api"
)

var (
	to   = flag.String("to", "eth", "to")
	from = flag.String("from", "btc", "from")
)

func main() {

	flag.Parse()

	series := &api.TimeSeries{
		Name: *to + "_to_" + *from,
		To:   strings.ToUpper(*to),
		From: strings.ToUpper(*from),
	}

	series.Run()
}
