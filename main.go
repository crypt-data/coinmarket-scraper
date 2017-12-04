package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"github.com/crypt-data/coinmarket-scraper/api"
	"github.com/kr/pretty"
)

func get(u url.URL) *api.Response {

	res, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var resp api.Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatal(err)
	}
	pretty.Logln(resp)

	return &resp
}

func main() {

	var u = url.URL{
		Scheme: "https",
		Host:   "min-api.cryptocompare.com",
		Path:   "data/histohour",
	}

	t := 1439521878
	for h := 0; h < 24*5; h++ {

		q := u.Query()
		for k, v := range map[string]string{
			"fsym":      "ETH",
			"tsym":      "BTC",
			"limit":     "60",
			"aggregate": "1",
			"toTs":      strconv.Itoa(t),
		} {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()

		pretty.Logln(u)

		resp := get(u)

		for _, tick := range resp.Data {
			tick.Put()
		}

		t += h * 60 * 60
	}
}
