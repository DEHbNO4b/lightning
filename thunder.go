package main

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

// type thunder struct {
// 	id           int
// 	claster      int
// 	polygon      [][]float32
// 	area         float32
// 	countStrikes int
// 	startTime    time.Time
// 	endTime      time.Time
// 	duration     time.Duration
// }

func calcThunderPolygons(db *sql.DB) error {
	var claster int
	rows, err := db.Query(queryCountClasters)
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

		_, err = db.Exec(queryCalcPolygons, claster)
		if err != nil {
			fmt.Println("makepolygons err:")
			return err
		}
	}
	return nil
}
func calcThunderArea(db *sql.DB) error {
	_, err := db.Exec(queryCalcArea)
	if err != nil {
		return err
	}
	return nil
}

func calcThunderCapacity(db *sql.DB) error {
	_, err := db.Exec(queryCalcCapacity)
	if err != nil {
		return err
	}
	return nil
}

func calcTimes(db *sql.DB) error {
	_, err := db.Exec(queryCalcStartTime)
	if err != nil {
		return err
	}
	_, err = db.Exec(queryCalcEndTime)
	if err != nil {
		return err
	}
	return nil
}
