package interfaces

import (
	"github.com/DEHbNO4b/lightning.git/internal/domain/models"
)

type Neighbours interface {
	GetNeighbourse(long, lat float32, eps int) (map[string]models.Stroke, error)
}
