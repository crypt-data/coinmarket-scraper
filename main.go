package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kr/pretty"
)

func main() {

	var url = url.URL{
		Scheme: "https",
		Host:   "min-api.cryptocompare.com",
		Path:   "data/histohour",
	}

	t := 1439521878
	for h := 0; h < 24*5; h++ {

		q := url.Query()
		for k, v := range map[string]string{
			"fsym":      "ETH",
			"tsym":      "BTC",
			"limit":     "60",
			"aggregate": "1",
			"toTs":      strconv.Itoa(t),
		} {
			q.Set(k, v)
		}
		url.RawQuery = q.Encode()

		pretty.Logln(url)

		res, _ := http.Get(url.String())
		body, _ := ioutil.ReadAll(res.Body)
		pretty.Logln(string(body))

		t += h * 60 * 60
	}
}
