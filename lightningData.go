package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var dsn string = "postgres://postgres:917836@localhost:5432/lightning?"

var queryMakeTab = ` drop table if exists strikes;
					create table if not exists strikes(
						id serial primary key,
						time timestamptz,
						latitude numeric(6,4),
						longitude numeric(6,4),
						signal smallint,
						cloud boolean,
						cluster integer
					);`

var queryInsert string = `INSERT INTO strikes (time,latitude,longitude,signal,cloud) 
							values($1,$2,$3,$4,$5)`

// структура разрядов молний
type lightning struct {
	time      time.Time
	latitude  float32
	longitude float32
	signal    int
	cloud     bool
	err       error
}

type lightningData struct {
	data []lightning
	db   *sql.DB
}

func NewLightningData() (*lightningData, error) {
	ld := lightningData{}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec(queryMakeTab)
	if err != nil {
		return nil, err
	}
	ld.db = db
	return &ld, nil
}
func (ld *lightningData) readFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stroke := parseLightning(line)
		ld.data = append(ld.data, stroke)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
func (ld *lightningData) loadDataToDb() error {
	for _, el := range ld.data {
		_, err := ld.db.Exec(queryInsert, el.time, el.latitude, el.longitude, el.signal, el.cloud)
		if err != nil {
			return err
		}
	}
	return nil
}
func parseLightning(s string) lightning {
	l := lightning{}
	data := strings.Split(s, "\t")
	//parse cloud
	cloud, err := strconv.ParseBool(data[21])
	if err != nil {
		l.err = err
	}
	l.cloud = cloud
	//parse signal
	sig, err := strconv.Atoi(data[10])
	if err != nil {
		l.err = err
	}
	l.signal = sig
	//parse longitude
	if long, err := strconv.ParseFloat(data[9], 32); err == nil {
		l.longitude = float32(long)
	}
	//parse latitude
	if lat, err := strconv.ParseFloat(data[8], 32); err == nil {
		l.latitude = float32(lat)
	}

	//parse time
	year, err := strconv.Atoi(data[1])
	if err != nil {
		l.err = err
	}
	month, err := strconv.Atoi(data[2])
	if err != nil {
		l.err = err
	}
	day, err := strconv.Atoi(data[3])
	if err != nil {
		l.err = err
	}
	hour, err := strconv.Atoi(data[4])
	if err != nil {
		l.err = err
	}
	min, err := strconv.Atoi(data[5])
	if err != nil {
		l.err = err
	}
	sec, err := strconv.Atoi(data[6])
	if err != nil {
		l.err = err
	}
	nano, err := strconv.Atoi(data[7])
	if err != nil {
		l.err = err
	}

	t := time.Date(year, time.Month(month), day, hour, min, sec, nano, time.UTC)
	l.time = t
	return l
}
