package model

const (
	Bulldog Breed = iota
	GermanShepherd
	Labrador
	GoldenRetriever
	FrenchBulldog
	Pug
	BorderCollie

	Small Size = iota
	Medium
	Large

	Puppy Age = iota
	Adult
	Senior
)

type Dog struct {
	Name  string
	Breed Breed
	Age   Age
	Size  Size
}

type Breed int16

func (b Breed) String() string {
	switch b {
	case Bulldog:
		return "Bulldog"
	case GermanShepherd:
		return "Pastor Aleman"
	case Labrador:
		return "Labrador"
	case GoldenRetriever:
		return "Golden Retriever"
	case FrenchBulldog:
		return "Bulldog Frances"
	case Pug:
		return "Pug"
	case BorderCollie:
		return "Border Collie"
	default:
		return "perro"
	}
}

type Size int

func (s Size) String() string {
	switch s {
	case Small:
		return "chico"
	case Medium:
		return "mediano"
	case Large:
		return "grande"
	default:
		return "desconocido"
	}
}

type Age int

func (a Age) String() string {
	switch a {
	case Puppy:
		return "cachorro"
	case Adult:
		return "adulto"
	case Senior:
		return "senior"
	default:
		return "desconocido"
	}
}
