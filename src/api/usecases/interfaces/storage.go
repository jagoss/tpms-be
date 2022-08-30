package interfaces

type Storage interface {
	SaveImgs(imgs [][]byte) (string, error)
	GetAllImgsName() ([]string, error)
	GetImgs(filePaths string) ([][]byte, error)
}
