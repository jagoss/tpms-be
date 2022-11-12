package interfaces

type Prediction interface {
	CalculateVector(dogID int64, imgsUrl string) error
	FindMatches(dogID uint) ([]uint, error)
}
