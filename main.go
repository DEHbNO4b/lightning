package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var filename string = "./source/Repr.2022.07.03.txt"

// структура разрядов молний
type lightning struct {
	time      time.Time
	longitude float32
	latitude  float32
	signal    int
	cloud     bool
	err       error
}

func parseLightning(s string) lightning {
	l := lightning{}
	data := strings.Split(s, "\t")
	//parse cloud
	cloud, err := strconv.ParseBool(data[21])
	if err != nil {
		l.err = err
	}
	l.cloud = cloud
	//parse signal
	sig, err := strconv.Atoi(data[10])
	if err != nil {
		l.err = err
	}
	l.signal = sig
	//parse longitude
	if long, err := strconv.ParseFloat(data[9], 32); err == nil {
		l.longitude = float32(long)
	}
	//parse latitude
	if lat, err := strconv.ParseFloat(data[8], 32); err == nil {
		l.latitude = float32(lat)
	}

	//parse time
	year, _ := strconv.Atoi(data[1])
	month, _ := strconv.Atoi(data[2])
	day, _ := strconv.Atoi(data[3])
	hour, _ := strconv.Atoi(data[4])
	min, _ := strconv.Atoi(data[5])
	sec, _ := strconv.Atoi(data[6])
	nano, _ := strconv.Atoi(data[7])
	t := time.Date(year, time.Month(month), day, hour, min, sec, nano, time.UTC)
	l.time = t
	return l
}
func main() {

	//открываем файл и считываем данные
	//получаем слайс молниевых структур
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("cannot open file %s", filename)
	}
	defer file.Close()
	l := make([]lightning, 1)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stroke := parseLightning(line)
		l = append(l, stroke)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	for _, el := range l {
		fmt.Println(el)
	}
	fmt.Println("размер слайса l:", len(l))
}
