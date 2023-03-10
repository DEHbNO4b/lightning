package maindb

import (
	"database/sql"
	"fmt"
)

var queryCalcStartTime string = `update thunders as t set startTime = (select min(time) from strikes as s where t.cluster = s.cluster );`
var queryCalcEndTime string = `update thunders as t set endTime = (select max(time) from strikes as s where t.cluster = s.cluster );`
var queryCountClasters string = `select DISTINCT cluster from strikes order by cluster;`
var queryCalcArea string = `update thunders set area = (select st_area(geog) / 1000000);`
var queryCalcCapacity string = `update thunders as t set capacity = (select count(*) from strikes as s where t.cluster = s.cluster);`
var queryCalcPolygons string = `insert into thunders (cluster,geog) 
								values ($1,(select st_transform(st_convexHull(st_collect(geog::geometry)),4326) 
								from strikes where cluster = $1));`

type PgThunderStorage struct {
	db *sql.DB
}

func NewPgThunderStorage(db *sql.DB) PgThunderStorage {
	return PgThunderStorage{db: db}
}

func (pgts PgThunderStorage) CalcThundersPolygon() error {
	var claster int
	rows, err := pgts.db.Query(queryCountClasters)
	if err != nil {
		fmt.Println("query count claster errr:", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&claster)
		if claster == -1 {
			continue
		}

		_, err = pgts.db.Exec(queryCalcPolygons, claster)
		if err != nil {
			fmt.Println("makepolygons err:", err)
			return err
		}
	}
	return nil
}
func (pgts PgThunderStorage) CalcThundersArea() error {
	_, err := pgts.db.Exec(queryCalcArea)
	if err != nil {
		return err
	}
	return nil
}

func (pgts PgThunderStorage) CalcThundersCapacity() error {
	_, err := pgts.db.Exec(queryCalcCapacity)
	if err != nil {
		return err
	}
	return nil
}

func (pgts PgThunderStorage) CalcTimes() error {
	_, err := pgts.db.Exec(queryCalcStartTime)
	if err != nil {
		return err
	}
	_, err = pgts.db.Exec(queryCalcEndTime)
	if err != nil {
		return err
	}
	return nil
}
