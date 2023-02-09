package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

var dsn string = "postgres://postgres:917836@localhost:5432/lightning?"
var filename string = "./public/Repr.2022.07.03.txt"

func main() {
	//получаем слайс молниевых структур
	dr := dataReader{}
	data := dr.readFromFile(filename)

	// fmt.Printf("%#v", data[5000])
	// fmt.Println()
	// fmt.Printf("%#v", data[5000].time)
	// fmt.Println()
	fmt.Println("datalen:", len(data))

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close() // (3)

	query := `INSERT INTO strikes (time,latitude,longitude)
	values($1,$2,$3)`
	result, err := db.Exec(query, data[0].time, data[0].latitude, data[0].longitude)
	fmt.Printf("%#v", result)
	fmt.Printf("%#v", err)
}
