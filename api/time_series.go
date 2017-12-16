package api

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	btcStart int64 = 1281960000
	ethStart int64 = 1439164800
)

type TimeSeries struct {
	Name string
	From string
	To   string

	db *sql.DB
}

func (series *TimeSeries) Run() {

	series.init()

	var start int64
	start, err := series.getLatestTick()
	if err != nil {
		if series.From == "ETH" {
			start = ethStart
		} else if series.From == "USD" {
			start = btcStart
		}
	}

	var u = &url.URL{
		Scheme: "https",
		Host:   "min-api.cryptocompare.com",
		Path:   "data/histohour",
	}

	// TODO paramaterize term with flags
	for t := start; t < time.Now().Unix(); t += 60 * 60 * 60 {

		resp := get(u, series.From, series.To, int(t))

		for _, tick := range resp.Data {
			series.put(&tick)
		}
	}
}

func (series *TimeSeries) getLatestTick() (int64, error) {

	log.Println("[INFO] getting latest tick...")

	var latest int
	err := series.db.QueryRow("select time from " + series.Name + " order by time desc limit 1").Scan(&latest)
	if err != nil || &latest == nil {
		log.Println("[ERROR] sqlite3:", err)
		return -1, err
	}

	return int64(latest), nil
}

func (series *TimeSeries) init() {

	dbPath := os.Getenv("HOME") + "/" + series.Name + ".db"
	log.Println("[INFO] opening db at", dbPath)
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("[FATAL] sqlite3: ", err)
	}
	series.db = database

	tablePath := os.Getenv("GOPATH") + "/src/github.com/crypt-data/coinmarket-scraper/create_" + series.Name + ".sql"
	log.Println("[INFO] creating table from", tablePath)
	b, err := ioutil.ReadFile(tablePath)
	if err != nil {
		log.Fatal("[FATAL] unix: ", err)
	}

	if _, err := series.db.Exec(string(b)); err != nil {
		log.Fatal("[FATAL] sqlite3: ", err)
	}
}

func (series *TimeSeries) put(tick *Tick) {

	// lazily load db
	if series.db == nil {
		series.init()
	}

	for {
		if err := tick.Put(series.Name, series.db); err != nil {
			log.Println("[ERROR] failed to put tick")
			log.Printf("[ERROR] sqlite3: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		return
	}
}
