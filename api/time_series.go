package api

import (
	"database/sql"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type TimeSeries struct {
	Name string

	db *sql.DB
}

func (series *TimeSeries) init() {

	database, err := sql.Open("sqlite3", "/Users/atec/"+series.Name+".db")
	if err != nil {
		log.Fatal("[FATAL] sqlite3:", err)
	}
	series.db = database

	b, err := ioutil.ReadFile("/Users/atec/go/src/github.com/crypt-data/coinmarket-scraper/create_" + series.Name + ".sql")
	if err != nil {
		log.Fatal("[FATAL] unix:", err)
	}

	if _, err := series.db.Exec(string(b)); err != nil {
		log.Fatal("[FATAL] sqlite3:", err)
	}
}

func (series *TimeSeries) Put(tick *Tick) {

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
