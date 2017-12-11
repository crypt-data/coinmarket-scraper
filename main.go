package main

import (
	"flag"
	"net/url"
	"strings"
	"time"

	"github.com/crypt-data/coinmarket-scraper/api"
)

const (
	btcStart int64 = 1281960000
	ethStart int64 = 1439164800
)

var (
	to   = flag.String("to", "eth", "to")
	from = flag.String("from", "btc", "from")
)

func main() {

	flag.Parse()

	// TODO stay dry by encapsulating gets
	series := &api.TimeSeries{
		Name: *to + "_to_" + *from,
	}

	var u = &url.URL{
		Scheme: "https",
		Host:   "min-api.cryptocompare.com",
		Path:   "data/histohour",
	}

	var start int64
	if *to == "eth" {
		start = ethStart
	} else if *to == "usd" {
		start = btcStart
	}

	// TODO paramaterize term with flags
	for t := start; t < time.Now().Unix(); t += 60 * 60 * 60 {

		// TODO stay dry by encapsulating gets
		resp := api.Get(u, strings.ToUpper(*to), strings.ToUpper(*from), int(t))

		for _, tick := range resp.Data {
			series.Put(&tick)
		}

	}
}
