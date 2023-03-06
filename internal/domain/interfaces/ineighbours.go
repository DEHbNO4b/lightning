package interfaces

type ineighbours interface {
	getNeighbourse(long, lat float32, eps int) (map[string]Stroke, error)
}
