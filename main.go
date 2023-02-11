package main

import (
	"fmt"

	"github.com/asmarques/geodist"
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

	lis := geodist.Point{Lat: float64(ld.data[0].latitude), Long: float64(ld.data[0].longitude)}
	sfo := geodist.Point{Lat: float64(ld.data[3000].latitude), Long: float64(ld.data[3000].longitude)}

	d := geodist.HaversineDistance(lis, sfo)
	fmt.Printf("Haversine: %.2f km\n", d)

	d, err = geodist.VincentyDistance(lis, sfo)
	if err != nil {
		fmt.Printf("Failed to converge: %v", err)
	}

	fmt.Printf("Vincenty: %.6f km\n", d)
}
