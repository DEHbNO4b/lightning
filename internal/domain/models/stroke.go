package models

import (
	"time"
)

type Stroke struct {
	time      time.Time
	latitude  float32
	longitude float32
	signal    int
	cloud     bool
	err       error
	claster   int
	id        int
}

func (s *Stroke) neighbours(n ineighbours, eps int) (map[string]Stroke, error) {
	return n.getNeighbourse(s.longitude, s.latitude, eps)

}
