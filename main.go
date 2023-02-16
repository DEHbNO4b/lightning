package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/stdlib"
)

var filename string = "./public/Repr.2022.07.03.txt"

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
	_, err = db.Exec(queryMakeIndex)
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

func main() {
	strokes, err := readFromFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println("datalen:", len(strokes))

	db, err := openDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//загружаем данные в базу данных
	data, err := loadDataToDb(strokes, db)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", data["1"])
	// var eps int = 50000 //метры радиуса поиска соседей
	// var minPts int = 3  //количество необходимых соседей

	// err = ld.calculateDbscan(eps, minPts)
	// if err != nil {
	// 	fmt.Println("dbscan err:", err)
	// }

}
