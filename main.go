package main

import (
	"net/url"
	"time"

	"github.com/crypt-data/coinmarket-scraper/api"
)

const (
	btcStart int64 = 1281960000
	ethStart int64 = 1439164800
)

func main() {

	// TODO paramaterize to and from with flags
	series := &api.TimeSeries{
		Name: "usd_to_btc",
	}

	var u = &url.URL{
		Scheme: "https",
		Host:   "min-api.cryptocompare.com",
		Path:   "data/histohour",
	}

	now := time.Now().Unix()

	// TODO paramaterize term with flags
	for t := btcStart; t < now; t += 60 * 60 * 60 {

		// TODO paramaterize to and from with flags
		resp := api.Get(u, "USD", "BTC", int(t))

		for _, tick := range resp.Data {
			series.Put(&tick)
		}

	}
}
