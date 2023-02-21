package main

import (
	"strconv"
	"time"
)

type stroke struct {
	time      time.Time
	latitude  float32
	longitude float32
	signal    int
	cloud     bool
	err       error
	claster   int
	id        int
}
type ineighbours interface {
	getNeighbourse(long, lat float32, eps int) (map[string]stroke, error)
}

func (s *stroke) neighbours(n ineighbours, eps int) (map[string]stroke, error) {
	return n.getNeighbourse(s.longitude, s.latitude, eps)

}
func dbscan(data map[string]stroke, neigh ineighbours, eps int, minPts int) (map[string]stroke, error) {
	claster := 0
	for key, val := range data { //начинаем обход данных
		if val.claster != 0 { //если уже просмотрен, то пропускаем
			continue
		}
		neighbours, err := val.neighbours(neigh, eps) //находим соседей
		if err != nil {
			return nil, err
		}
		delete(neighbours, key)
		if len(neighbours) < minPts { //если соседей меньше чем minPts то помечаем как шум
			stroke := data[key]
			stroke.claster = -1
			data[key] = stroke
			continue
		}
		claster++

		stroke := data[key] //начинаем новый кластер
		stroke.claster = claster
		data[key] = stroke
		for _, val := range neighbours {
			expandClaster(neigh, data, claster, val, eps, minPts)
		}
	}
	return data, nil
}
func expandClaster(neigh ineighbours, data map[string]stroke, claster int, s stroke, eps int, minPts int) {
	d := data[strconv.Itoa(s.id)]
	if d.claster > 0 {
		return
	}
	d.claster = claster
	data[strconv.Itoa(s.id)] = d
	n, _ := s.neighbours(neigh, eps)
	delete(n, strconv.Itoa(s.id))

	if len(n) > minPts {
		for _, v := range n {
			expandClaster(neigh, data, claster, v, eps, minPts)
		}

	}
}
