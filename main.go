package main

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

var filename string = "./public/Repr.2022.07.03.txt"

func main() {
	//получаем структуру lightningData для работы с данными молний
	ld, err := NewLightningData()
	if err != nil {
		fmt.Println(err)
	}
	defer ld.db.Close()
	err = ld.readFromFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("datalen:", len(ld.data))
	err = ld.loadDataToDb()
	if err != nil {
		panic(err)
	}

	fmt.Println("start distanse")
	fmt.Printf("%#v", ld.data[0])
}
