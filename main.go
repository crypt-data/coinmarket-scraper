package main

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kr/pretty"
)

func main() {

	var url = url.URL{
		Scheme: "https",
		Host:   "min-api.cryptocompare.com",
		Path:   "data/histohour",
	}

	q := url.Query()
	for k, v := range map[string]string{
		"fsym":      "ETH",
		"tsym":      "BTC",
		"limit":     "60",
		"aggregate": "1",
		"toTs":      "1439521878",
	} {
		q.Set(k, v)
	}
	url.RawQuery = q.Encode()

	res, _ := http.Get(url.String())
	body, _ := ioutil.ReadAll(res.Body)
	pretty.Logln(string(body))

}
