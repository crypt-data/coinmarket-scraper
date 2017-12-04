package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

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

func (tick *Tick) Put() {

	db, err := sql.Open("sqlite3", "/keybase/team/crypt_data/EthToBtc.db")
	if err != nil {
		log.Fatal("failed to open db", err)
	}

	b, err := ioutil.ReadFile("create_table.sql")
	if err != nil {
		log.Fatal("failed to read create_table.sql", err)
	}

	if _, err := db.Exec(string(b)); err != nil {
		log.Fatal("failed to create table", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("failed getting tx", err)
	}

	stmt, err := tx.Prepare("insert or replace into EthToBtc (time, close, high, low, open, volumefrom, volumeto) values (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("failed preparing tick", tick, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(tick.Time, tick.Close, tick.High, tick.Low, tick.Open, tick.VolumeFrom, tick.VolumeTo)
	if err != nil {
		log.Fatal("failed execing tick", tick, err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("failed committing tick", tick, err)
	}

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

		for _, tick := range resp.Data {
			tick.Put()
		}

		t += h * 60 * 60
	}
}
