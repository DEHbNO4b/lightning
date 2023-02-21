package main

import (
	"fmt"
	"strconv"

	_ "github.com/jackc/pgx/stdlib"
)

var filename string = "./public/Repr.2022.07.03.txt"

//var testfilename string = "./public/test.txt"

func main() {
	strokes, err := readFromFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println("Readed data len:", len(strokes))

	DB, err := openDB()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	//загружаем данные в базу данных
	data, err := loadDataToDb(strokes, DB)
	if err != nil {
		panic(err)
	}

	neigh := NewNeighbours(DB)

	var eps int = 80000 //метры радиуса поиска соседей
	var minPts int = 2  //количество необходимых соседей

	data, err = dbscan(data, &neigh, eps, minPts)
	if err != nil {
		fmt.Println("dbscan err:", err)
	}
	for key, el := range data {
		id, _ := strconv.Atoi(key)
		_, err := DB.Exec(`UPDATE strikes SET cluster = $1 WHERE id = $2`, el.claster, id)
		if err != nil {
			fmt.Println(err)
		}
	}
}
