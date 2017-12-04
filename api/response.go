package api

import (
	"database/sql"
	"io/ioutil"
	"log"

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
	database, err := sql.Open("sqlite3", "/keybase/team/crypt_data/EthToBtc.db")
	if err != nil {
		log.Fatal("failed to open db", err)
	}

	b, err := ioutil.ReadFile("create_table.sql")
	if err != nil {
		log.Fatal("failed to read create_table.sql", err)
	}

	db = database

	if _, err := db.Exec(string(b)); err != nil {
		log.Fatal("failed to create table", err)
	}
}

func (tick *Tick) Put() {

	// lazily load db
	if db == nil {
		Init()
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
