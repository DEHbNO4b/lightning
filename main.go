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

	var eps int = 100000 //метры радиуса поиска соседей
	var minPts int = 2   //количество необходимых соседей

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
	for key, val := range data {

		if val.claster != 0 {
			continue
		}
		neighbours, err := val.neighbours(db, eps)
		if err != nil {
			return nil, err
		}
		if len(neighbours) < minPts {
			// stroke := data[key]
			// stroke.claster = -1
			// data[key] = stroke
			continue
		}
		claster++
		count := 0
		stroke := data[key]
		stroke.claster = claster
		data[key] = stroke
		count++
		seed := neighbours
		delete(seed, key)
		for key, val := range seed {
			if data[key].claster > 0 {
				continue
			}

			n, err := val.neighbours(db, eps)
			if err != nil {
				return nil, err
			}
			if len(n) >= minPts {
				for k, v := range n {
					seed[k] = v
				}
				stroke := data[key]
				stroke.claster = claster
				data[key] = stroke
				s := seed[key]
				s.claster = claster
				seed[key] = s
				count++
			}

		}
		fmt.Printf("cluster= %d , count = %d\n", claster, count)
	}
	return data, nil
}
