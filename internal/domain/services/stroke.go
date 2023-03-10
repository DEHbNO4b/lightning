package services

import (
	"strconv"

	"github.com/DEHbNO4b/lightning.git/internal/domain/interfaces"
	"github.com/DEHbNO4b/lightning.git/internal/domain/models"
)

func Dbscan(data map[string]models.Stroke, neigh interfaces.Neighbours, eps int, minPts int) (map[string]models.Stroke, error) {
	claster := 0
	for key, val := range data { //начинаем обход данных
		if val.Claster != 0 { //если уже просмотрен, то пропускаем
			continue
		}
		//neighbours, err := val.neighbours(neigh, eps) //находим соседей
		neighbours, err := neigh.GetNeighbourse(val.Longitude, val.Latitude, eps)
		if err != nil {
			return nil, err
		}
		delete(neighbours, key)
		if len(neighbours) < minPts { //если соседей меньше чем minPts то помечаем как шум
			stroke := data[key]
			stroke.Claster = -1
			data[key] = stroke
			continue
		}
		claster++

		stroke := data[key] //начинаем новый кластер
		stroke.Claster = claster
		data[key] = stroke
		for _, val := range neighbours {
			expandClaster(neigh, data, claster, val, eps, minPts)
		}
	}
	return data, nil
}
func expandClaster(neigh interfaces.Neighbours, data map[string]models.Stroke, claster int, s models.Stroke, eps int, minPts int) {
	d := data[strconv.Itoa(s.Id)]
	if d.Claster > 0 {
		return
	}
	d.Claster = claster
	data[strconv.Itoa(s.Id)] = d
	//n, _ := s.Neighbours(neigh, eps)
	n, _ := neigh.GetNeighbourse(s.Longitude, s.Latitude, eps)
	delete(n, strconv.Itoa(s.Id))

	if len(n) > minPts {
		for _, v := range n {
			expandClaster(neigh, data, claster, v, eps, minPts)
		}

	}
}
