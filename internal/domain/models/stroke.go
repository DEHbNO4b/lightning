package models

import (
	"time"
)

type Stroke struct {
	Time      time.Time
	Latitude  float32
	Longitude float32
	Signal    int
	Cloud     bool
	Err       error
	Claster   int
	Id        int
}
