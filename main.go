package main

import (
	"flag"
	"strings"

	"github.com/crypt-data/coinmarket-scraper/api"
)

var (
	from = flag.String("from", "eth", "to")
	to   = flag.String("to", "btc", "from")
)

func main() {

	flag.Parse()

	series := &api.TimeSeries{
		Name: *from + "_to_" + *to,
		From: strings.ToUpper(*from),
		To:   strings.ToUpper(*to),
	}

	series.Run()
}
