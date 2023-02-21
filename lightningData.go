package main

import (
	"database/sql"
	"strconv"
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

type lightningDB struct {
	db *sql.DB
}

func NewLightningDB(db *sql.DB) lightningDB {
	return lightningDB{db: db}
}
func (n *lightningDB) getNeighbourse(long, lat float32, eps int) (map[string]stroke, error) {
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

func (n *lightningDB) openDB() error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	n.db = db
	return nil
}
func (n *lightningDB) makeTab() error {
	_, err := n.db.Exec(queryMakeTab)
	if err != nil {
		return err
	}
	return nil
}
func (n *lightningDB) loadRawToDb(strokes []stroke) (map[string]stroke, error) {
	data := make(map[string]stroke, len(strokes))
	for i, el := range strokes {
		var idInDB int
		err := n.db.QueryRow(queryInsert, el.time, el.longitude, el.latitude, el.longitude, el.latitude, el.signal, el.cloud).Scan(&idInDB)
		if err != nil {
			return nil, err
		}
		strokes[i].id = idInDB
		data[strconv.Itoa(idInDB)] = strokes[i]
	}

	return data, nil
}
