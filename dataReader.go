package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// структура разрядов молний
type lightning struct {
	time      time.Time
	latitude  float32
	longitude float32
	signal    int
	cloud     bool
	err       error
}

type dataReader struct {
}

func (dr dataReader) readFromFile(filename string) []lightning {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("cannot open file %s", filename)
	}
	defer file.Close()
	data := make([]lightning, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stroke := parseLightning(line)
		data = append(data, stroke)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return data
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
	year, err := strconv.Atoi(data[1])
	if err != nil {
		l.err = err
	}
	month, err := strconv.Atoi(data[2])
	if err != nil {
		l.err = err
	}
	day, err := strconv.Atoi(data[3])
	if err != nil {
		l.err = err
	}
	hour, err := strconv.Atoi(data[4])
	if err != nil {
		l.err = err
	}
	min, err := strconv.Atoi(data[5])
	if err != nil {
		l.err = err
	}
	sec, err := strconv.Atoi(data[6])
	if err != nil {
		l.err = err
	}
	nano, err := strconv.Atoi(data[7])
	if err != nil {
		l.err = err
	}

	t := time.Date(year, time.Month(month), day, hour, min, sec, nano, time.UTC)
	l.time = t
	return l
}
