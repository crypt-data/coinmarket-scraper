package main

import (
	"net/url"

	"github.com/crypt-data/coinmarket-scraper/api"
)

const (
	btcStart = 1281960000
	ethStart = 1439164800
)

func main() {

	series := &api.TimeSeries{
		Name: "usd_to_btc",
	}

	var u = &url.URL{
		Scheme: "https",
		Host:   "min-api.cryptocompare.com",
		Path:   "data/histohour",
	}

	now := 1512849201
	for t := btcStart; t < now; t += 60 * 60 * 60 {

		resp := api.Get(u, "USD", "BTC", t)

		for _, tick := range resp.Data {
			series.Put(&tick)
		}

	}
}
