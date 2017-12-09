package api

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Tick struct {
	Time       int     `json:"time"`
	Close      float64 `json:"close"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Open       float64 `json:"open"`
	VolumeFrom float64 `json:"volumefrom"`
	VolumeTo   float64 `json:"volumeto"`
}

func (tick *Tick) Put(tableName string, db *sql.DB) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert or replace into " + tableName + " (time, close, high, low, open, volumefrom, volumeto) values (?, ?, ?, ?, ?, ?, ?)")
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
