package api

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db     *sql.DB
	logger *syslog.Writer
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

func Init() {

	syslogger, err := syslog.Dial("", "", syslog.LOG_USER, "")
	if err != nil {
		log.Fatal("syslog:", err)
	}
	logger = syslogger

	database, err := sql.Open("sqlite3", "/keybase/team/crypt_data/EthToBtc.db")
	if err != nil {
		logger.Emerg("failed to open db")
		log.Fatal("sqlite3:", err)
	}
	db = database

	b, err := ioutil.ReadFile("/Users/atec/go/src/github.com/crypt-data/coinmarket-scraper/create_table.sql")
	if err != nil {
		logger.Emerg("failed to read create_table.sql")
		log.Fatal("unix:", err)
	}

	if _, err := db.Exec(string(b)); err != nil {
		logger.Emerg("failed to create table")
		log.Fatal("sqlite3:", err)
	}
}

func (tick *Tick) Put() {

	// lazily load db
	if db == nil {
		Init()
	}

	for {
		if err := tick.put(); err != nil {
			logger.Err("failed to put tick")
			logger.Err(fmt.Sprintf("sqlite3: %v", err))
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
	logger.Info(fmt.Sprintf("successfully inserted tick (%d, %f, %f, %f, %f, %f, %f)", tick.Time, tick.Close, tick.High, tick.Low, tick.Open, tick.VolumeFrom, tick.VolumeTo))
	return nil
}

type ConversionTypeStruct struct {
	Type             string `json:"type"`
	ConversionSymbol string `json:"conversionSymbol"`
}
