package api

import (
	"database/sql"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

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

func Init() {

	database, err := sql.Open("sqlite3", "/Users/atec/EthToBtc.db")
	if err != nil {
		log.Fatal("[FATAL] sqlite3:", err)
	}
	db = database

	b, err := ioutil.ReadFile("/Users/atec/go/src/github.com/crypt-data/coinmarket-scraper/create_table.sql")
	if err != nil {
		log.Fatal("[FATAL] unix:", err)
	}

	if _, err := db.Exec(string(b)); err != nil {
		log.Fatal("[FATAL] sqlite3:", err)
	}
}

func (tick *Tick) Put() {

	// lazily load db
	if db == nil {
		Init()
	}

	for {
		if err := tick.put(); err != nil {
			log.Println("[ERROR] failed to put tick")
			log.Printf("[ERROR] sqlite3: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		return
	}
}

func (tick *Tick) put() error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert or replace into EthToBtc (time, close, high, low, open, volumefrom, volumeto) values (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tick.Time, tick.Close, tick.High, tick.Low, tick.Open, tick.VolumeFrom, tick.VolumeTo)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("[INFO] successfully inserted tick (%d, %f, %f, %f, %f, %f, %f)", tick.Time, tick.Close, tick.High, tick.Low, tick.Open, tick.VolumeFrom, tick.VolumeTo)
	return nil
}

type ConversionTypeStruct struct {
	Type             string `json:"type"`
	ConversionSymbol string `json:"conversionSymbol"`
}
