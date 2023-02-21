package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type vails struct {
	raw    []stroke
	dbData map[string]stroke
}

func (v *vails) readFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stroke := parseStroke(line)
		v.raw = append(v.raw, stroke)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
func parseStroke(s string) stroke {
	l := stroke{}
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
