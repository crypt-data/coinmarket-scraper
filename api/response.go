package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

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

func Get(u *url.URL, from, to string, t int) *Response {

	q := u.Query()
	for k, v := range map[string]string{
		"fsym":      from,
		"tsym":      to,
		"limit":     "60",
		"aggregate": "1",
		"toTs":      strconv.Itoa(t),
	} {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	url := u.String()
	log.Println("[INFO] getting...", url)
	res, err := http.Get(url)
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

	if len(resp.Data) == 0 {
		pretty.Logln("[INFO] ", resp)
		time.Sleep(1 * time.Minute)
		resp = *Get(u, from, to, t)
	}

	return &resp
}

type ConversionTypeStruct struct {
	Type             string `json:"type"`
	ConversionSymbol string `json:"conversionSymbol"`
}
