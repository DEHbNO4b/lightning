package main

import (
	"database/sql"
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

	var eps int = 80000 //метры радиуса поиска соседей
	var minPts int = 2  //количество необходимых соседей

	data, err = dbscan(data, DB, eps, minPts)
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

func dbscan(data map[string]stroke, db *sql.DB, eps int, minPts int) (map[string]stroke, error) {
	claster := 0
	for key, val := range data { //начинаем обход данных
		if val.claster != 0 { //если уже просмотрен, то пропускаем
			continue
		}
		neighbours, err := val.neighbours(db, eps) //находим соседей
		if err != nil {
			return nil, err
		}
		delete(neighbours, key)
		if len(neighbours) < minPts { //если соседей меньше чем minPts то помечаем как шум
			stroke := data[key]
			stroke.claster = -1
			data[key] = stroke
			continue
		}
		claster++

		stroke := data[key] //начинаем новый кластер
		stroke.claster = claster
		data[key] = stroke
		for _, val := range neighbours {
			expandClaster(db, data, claster, val, eps, minPts)
		}
	}
	return data, nil
}
func expandClaster(db *sql.DB, data map[string]stroke, claster int, s stroke, eps int, minPts int) {
	d := data[strconv.Itoa(s.id)]
	if d.claster > 0 {
		return
	}
	d.claster = claster
	data[strconv.Itoa(s.id)] = d
	n, _ := s.neighbours(db, eps)
	delete(n, strconv.Itoa(s.id))

	if len(n) > minPts {
		for _, v := range n {
			expandClaster(db, data, claster, v, eps, minPts)
		}

	}
}
