package main

import (
	"bufio"
	"database/sql"
	"os"
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
					);
					CREATE INDEX ON strikes USING GIST(geog);`

var queryInsert string = `INSERT INTO strikes (time,longitude,latitude,geog,signal,cloud) 
							VALUES($1,$2,$3,ST_MakePoint($4, $5)::GEOGRAPHY,$6,$7)
							RETURNING ID;`
var queryNeighbors string = `SELECT id,longitude,latitude FROM strikes WHERE geog<->st_setSRID(st_makePoint($1,$2),4326)::GEOGRAPHY < $3;`

type neighbours struct {
	db *sql.DB
}

func NewNeighbours(db *sql.DB) neighbours {
	return neighbours{db: db}
}
func (n *neighbours) get(long, lat float32, eps int) (map[string]stroke, error) {
	var ans = make(map[string]stroke)
	var latitude, longitude float32
	var id int
	rows, err := n.db.Query(queryNeighbors, long, lat, eps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&id, &longitude, &latitude); err != nil {
			return nil, err
		}
		s := stroke{id: id, longitude: longitude, latitude: latitude}
		ans[strconv.Itoa(id)] = s
	}
	return ans, nil
}
func readFromFile(filename string) ([]stroke, error) {
	var data []stroke
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stroke := parseStroke(line)
		data = append(data, stroke)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return data, nil
}
func openDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(queryMakeTab)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func loadDataToDb(strokes []stroke, db *sql.DB) (map[string]stroke, error) {
	data := make(map[string]stroke, len(strokes))
	for i, el := range strokes {
		var idInDB int
		err := db.QueryRow(queryInsert, el.time, el.longitude, el.latitude, el.longitude, el.latitude, el.signal, el.cloud).Scan(&idInDB)
		if err != nil {
			return nil, err
		}
		strokes[i].id = idInDB
		data[strconv.Itoa(idInDB)] = strokes[i]
	}

	return data, nil
}

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
