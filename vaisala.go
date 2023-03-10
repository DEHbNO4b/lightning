package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DEHbNO4b/lightning.git/internal/domain/models"
)

type vails struct {
	raw    []models.Stroke
	dbData map[string]models.Stroke
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
func parseStroke(s string) models.Stroke {
	l := models.Stroke{}
	data := strings.Split(s, "\t")
	//parse cloud
	cloud, err := strconv.ParseBool(data[21])
	if err != nil {
		l.Err = err
	}
	l.Cloud = cloud
	//parse signal
	sig, err := strconv.Atoi(data[10])
	if err != nil {
		l.Err = err
	}
	l.Signal = sig
	//parse longitude
	if long, err := strconv.ParseFloat(data[9], 32); err == nil {
		l.Longitude = float32(long)
	}
	//parse latitude
	if lat, err := strconv.ParseFloat(data[8], 32); err == nil {
		l.Latitude = float32(lat)
	}

	//parse time
	year, err := strconv.Atoi(data[1])
	if err != nil {
		l.Err = err
	}
	month, err := strconv.Atoi(data[2])
	if err != nil {
		l.Err = err
	}
	day, err := strconv.Atoi(data[3])
	if err != nil {
		l.Err = err
	}
	hour, err := strconv.Atoi(data[4])
	if err != nil {
		l.Err = err
	}
	min, err := strconv.Atoi(data[5])
	if err != nil {
		l.Err = err
	}
	sec, err := strconv.Atoi(data[6])
	if err != nil {
		l.Err = err
	}
	nano, err := strconv.Atoi(data[7])
	if err != nil {
		l.Err = err
	}

	t := time.Date(year, time.Month(month), day, hour, min, sec, nano, time.UTC)
	l.Time = t
	return l
}
