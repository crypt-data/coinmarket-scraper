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

type response struct {
	Status            string               `json:"Response"`
	Type              int                  `json:"Type"`
	Aggregated        bool                 `json:"Aggregated"`
	Data              []Tick               `json:"Data"`
	TimeTo            int                  `json:"TimeTo"`
	TimeFrom          int                  `json:"TimeFrom"`
	FirstValueInArray bool                 `json:"FirstValueInArray"`
	ConversionType    ConversionTypeStruct `json:"ConversionType"`
}

func get(u *url.URL, from, to string, limit, t int) *response {

	for {
		q := u.Query()
		for k, v := range map[string]string{
			"fsym":  from,
			"tsym":  to,
			"limit": strconv.Itoa(limit),
			"toTs":  strconv.Itoa(t),
		} {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()

		url := u.String()
		log.Println("[INFO] getting", url, "...")
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var resp response
		err = json.Unmarshal(body, &resp)
		if err != nil {
			log.Fatal(err)
		}

		// TODO check error
		if len(resp.Data) == 0 {
			pretty.Logln("[INFO] ", resp)
			time.Sleep(1 * time.Minute)
			continue
		}

		return &resp
	}
}

type ConversionTypeStruct struct {
	Type             string `json:"type"`
	ConversionSymbol string `json:"conversionSymbol"`
}
