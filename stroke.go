package main

import (
	"database/sql"
	"strconv"
	"time"
)

var queryNeighbors string = `SELECT id,longitude,latitude FROM strikes WHERE geog<->st_setSRID(st_makePoint($1,$2),4326)::GEOGRAPHY < $3;`

type stroke struct {
	time      time.Time
	latitude  float32
	longitude float32
	signal    int
	cloud     bool
	err       error
	claster   int
	id        int
}

func (s *stroke) neighbours(db *sql.DB, eps int) (map[string]stroke, error) {
	var ans = make(map[string]stroke)
	var lat, long float32
	var id int
	rows, err := db.Query(queryNeighbors, s.longitude, s.latitude, eps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&id, &long, &lat); err != nil {
			return nil, err
		}
		s := stroke{id: id, longitude: long, latitude: lat}
		ans[strconv.Itoa(id)] = s
	}
	return ans, nil
}
