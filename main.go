package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

var filename string = "./public/Repr.2022.07.03.txt"

func main() {
	strokes, err := readFromFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println("datalen:", len(strokes))

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
	fmt.Printf("%#v", data["1"])

	s := stroke{latitude: 47, longitude: 47}
	n, err := s.neighbours(DB, 100000)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("count neighbors:", len(n))

	var eps int = 50000 //метры радиуса поиска соседей
	var minPts int = 3  //количество необходимых соседей

	data, err = dbscan(data, DB, eps, minPts)
	if err != nil {
		fmt.Println("dbscan err:", err)
	}

	fmt.Printf("%#v", data["4000"])

}

func dbscan(data map[string]stroke, db *sql.DB, eps int, minPts int) (map[string]stroke, error) {
	for key, val := range data {
		claster := 0
		if val.claster != 0 {
			continue
		}
		neighbours, err := val.neighbours(db, eps)
		if err != nil {
			return nil, err
		}
		if len(neighbours) < minPts {
			stroke := data[key]
			stroke.claster = -1
			data[key] = stroke
			continue
		}
		claster++
		stroke := data[key]
		stroke.claster = claster
		data[key] = stroke
		seed := neighbours
		delete(seed, key)
		for key, val := range seed {
			if val.claster != 0 {
				continue
			}
			stroke := data[key]
			stroke.claster = claster
			data[key] = stroke
			n, err := stroke.neighbours(db, eps)
			if err != nil {
				return nil, err
			}
			if len(n) >= minPts {
				for k, v := range n {
					seed[k] = v
				}
			}
		}

	}
	return data, nil
}
