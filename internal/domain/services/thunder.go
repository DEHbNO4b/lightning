package services

import "github.com/DEHbNO4b/lightning.git/internal/domain/interfaces"

type ThunderService struct {
	ThunderStorage interfaces.CalcThunder
}

func NewThunderService(ict interfaces.CalcThunder) *ThunderService {
	return &ThunderService{ThunderStorage: ict}
}
func (ts *ThunderService) CalcAllThanders() error {
	err := ts.ThunderStorage.CalcThundersPolygon()
	if err != nil {
		return err
	}
	err = ts.ThunderStorage.CalcThundersArea()
	if err != nil {
		return err
	}
	err = ts.ThunderStorage.CalcThundersCapacity()
	if err != nil {
		return err
	}
	err = ts.ThunderStorage.CalcTimes()
	if err != nil {
		return err
	}

	return nil
}
