package main

import "time"

//структура разрядов молний
type lightning struct {
	time      time.Time
	nano      time.Duration
	longitude float32
	latitude  float32
	signal    int
	cloud     bool
}

func main() {

	//открываем файл и считываем данные
	//получаем слайс молниевых структур

	//
}
