package models

import "time"

type Thunder struct {
	Id           int
	Claster      int
	Polygon      [][]float32
	Area         float32
	CountStrikes int
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
}
