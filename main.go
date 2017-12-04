package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kr/pretty"
)

type Response struct {
	Status            string               `json:"Response"`
	Type              int                  `json:"Type"`
	Aggregated        bool                 `json:"Aggregated"`
	Data              []Tick               `json:"Data"`
	TimeTo            int                  `json:"TimeTo"`
	TimeFrom          int                  `json:"TimeFrom"`
	FirstValueInArray bool                 `json:"FirstValueInArray"`
	ConversionType    ConversionTypeStruct `json:"ConversionType"`
}

type Tick struct {
	Time       int     `json:"time"`
	Close      float64 `json:"close"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Open       float64 `json:"open"`
	VolumeFrom float64 `json:"volumefrom"`
	VolumeTo   float64 `json:"volumeto"`
}

type ConversionTypeStruct struct {
	Type             string `json:"type"`
	ConversionSymbol string `json:"conversionSymbol"`
}

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

		res, err := http.Get(url.String())
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var resp Response
		err = json.Unmarshal(body, &resp)
		if err != nil {
			log.Fatal(err)
		}
		pretty.Logln(resp)

		t += h * 60 * 60
	}
}
