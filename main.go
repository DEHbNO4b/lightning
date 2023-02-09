package main

import (
	"fmt"
)

var filename string = "./public/Repr.2022.07.03.txt"

func main() {
	//получаем слайс молниевых структур
	dr := dataReader{}
	data := dr.readFromFile(filename)

	fmt.Printf("%#v", data[0])
	fmt.Println()
	fmt.Println("datalen:", len(data))

}
