package main

import (
	"fmt"
	"strconv"

	_ "github.com/jackc/pgx/stdlib"
)

var filename string = "./public/Repr.2022.07.03.txt"

//var testfilename string = "./public/test.txt"

func main() {
	vai := vails{}
	err := vai.readFromFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println("Readed data len:", len(vai.raw))

	ldb := lightningDB{}
	err = ldb.openDB()
	if err != nil {
		panic(err)
	}
	defer ldb.db.Close()
	//загружаем данные в базу данных
	ldb.makeTab()
	data, err := ldb.loadRawToDb(vai.raw)
	if err != nil {
		panic(err)
	}
	vai.dbData = data
	//ldb := NewLightningDB(DB)

	var eps int = 80000 //метры радиуса поиска соседей
	var minPts int = 2  //количество необходимых соседей

	data, err = dbscan(data, &ldb, eps, minPts)
	if err != nil {
		fmt.Println("dbscan err:", err)
	}
	for key, el := range data {
		id, _ := strconv.Atoi(key)
		_, err := ldb.db.Exec(`UPDATE strikes SET cluster = $1 WHERE id = $2`, el.claster, id)
		if err != nil {
			fmt.Println(err)
		}
	}
}
