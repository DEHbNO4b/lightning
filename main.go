package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

var dsn string = "postgres://postgres:917836@localhost:5432/lightning?"
var filename string = "./public/Repr.2022.07.03.txt"
var queryInsert string = `INSERT INTO strikes (time,latitude,longitude,signal,cloud) 
values($1,$2,$3,$4,$5)`
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

func main() {
	//получаем слайс молниевых структур
	dr := dataReader{}
	data := dr.readFromFile(filename)
	fmt.Println("datalen:", len(data))

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	resp, err := db.Exec(queryMakeTab)
	if err != nil {
		fmt.Println("query make tab err:", err)
	}
	fmt.Println("query make tab resp:", resp)
	for i, _ := range data {
		_, err := db.Exec(queryInsert, data[i].time, data[i].latitude, data[i].longitude, data[i].signal, data[i].cloud)
		if err != nil {
			fmt.Println(err)
		}
	}
}
