package main

import (
	"database/sql"
	"time"
)

var queryNeighbors string = `SELECT id FROM strikes WHERE geog<->st_setSRID(st_makePoint($1,$2),4326)::GEOGRAPHY < $3;`

type stroke struct {
	time      time.Time
	latitude  float32
	longitude float32
	signal    int
	cloud     bool
	err       error
	cluster   int
	id        int
}

func (s *stroke) neighbours(db *sql.DB, eps int) ([]int, error) {
	var ans []int
	var id int
	rows, err := db.Query(queryNeighbors, s.longitude, s.latitude, eps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}

		ans = append(ans, id)
	}
	return ans, nil
}
