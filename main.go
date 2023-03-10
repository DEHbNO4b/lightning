package main

import (
	"fmt"

	"github.com/DEHbNO4b/lightning.git/internal/domain/services"
	"github.com/DEHbNO4b/lightning.git/internal/maindb"
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

	ldb := maindb.LightningDB{}
	err = ldb.OpenDB()
	if err != nil {
		panic(err)
	}
	defer ldb.Db.Close()
	//загружаем данные в базу данных
	ldb.MakeTab()
	data, err := ldb.LoadRawToDb(vai.raw)
	if err != nil {
		panic(err)
	}
	vai.dbData = data

	var eps int = 80000 //метры радиуса поиска соседей
	var minPts int = 2  //количество необходимых соседей

	data, err = services.Dbscan(data, &ldb, eps, minPts)
	if err != nil {
		fmt.Println("dbscan err:", err)
	}
	err = ldb.LoadClasterToDb(data)
	if err != nil {
		fmt.Println("load claster to db err:", err)
	}
	pgts := maindb.NewPgThunderStorage(ldb.Db)

	ts := services.NewThunderService(pgts)
	err = ts.CalcAllThanders()
	if err != nil {
		panic(err)
	}
}
