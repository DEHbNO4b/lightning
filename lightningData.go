package main

import (
	"strconv"
	"strings"
	"time"
)

var dsn string = "postgres://postgres:917836@localhost:5432/lightning?"

var queryMakeTab string = ` DROP TABLE IF EXISTS strikes;
					CREATE TABLE IF NOT EXISTS strikes(
						id serial primary key,
						time timestamptz,
						latitude numeric(6,4),
						longitude numeric(6,4),
						geog GEOGRAPHY(Point),
						signal smallint,
						cloud boolean,
						cluster integer
					);`
var queryMakeIndex string = `CREATE INDEX ON strikes USING GIST(geog);`

var queryInsert string = `INSERT INTO strikes (time,longitude,latitude,geog,signal,cloud) 
							VALUES($1,$2,$3,ST_MakePoint($4, $5)::GEOGRAPHY,$6,$7)
							RETURNING ID;`

func parseStroke(s string) stroke {
	l := stroke{}
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
