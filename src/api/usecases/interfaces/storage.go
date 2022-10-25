package interfaces

type Storage interface {
	SaveImgs(imgs []string) (string, error)
	GetAllImgsName() ([]string, error)
	GetImgs(filePaths string) ([]string, error)
}
