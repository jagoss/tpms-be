package model

import (
	"strings"
	"time"
)

const (
	UnknowBreed Breed = iota
	Other
	Affenpinscher
	AfghanHound
	AfricanHuntingDog
	Airedale
	AmericanStaffordshireTerrier
	Appenzeller
	AustralianTerrier
	Basenji
	Basset
	Beagle
	BedlingtonTerrier
	BerneseMountainDog
	BlackAndTanCoonhound
	BlenheimSpaniel
	Bloodhound
	Bluetick
	BorderCollie
	BorderTerrier
	Borzoi
	BostonBull
	BouvierDesFlandres
	Boxer
	BrabanconGriffon
	Briard
	BrittanySpaniel
	BullMastiff
	Cairn
	Cardigan
	ChesapeakeBayRetriever
	Chihuahua
	Chow
	Clumber
	CockerSpaniel
	Collie
	CurlyCoatedRetriever
	DandieDinmont
	Dhole
	Dingo
	Doberman
	EnglishFoxhound
	EnglishSetter
	EnglishSpringer
	Entlebucher
	EskimoDog
	FlatCoatedRetriever
	FrenchBulldog
	GermanShepherd
	GermanShortHairedPointer
	GiantSchnauzer
	GoldenRetriever
	GordonSetter
	GreatDane
	GreatPyrenees
	GreaterSwissMountainDog
	Groenendael
	IbizanHound
	IrishSetter
	IrishTerrier
	IrishWaterSpaniel
	IrishWolfhound
	ItalianGreyhound
	JapaneseSpaniel
	Keeshond
	Kelpie
	KerryBlueTerrier
	Komondor
	Kuvasz
	LabradorRetriever
	LakelandTerrier
	Leonberg
	Lhasa
	Malamute
	Malinois
	MalteseDog
	MexicanHairless
	MiniaturePinscher
	MiniaturePoodle
	MiniatureSchnauzer
	Newfoundland
	NorfolkTerrier
	NorwegianElkhound
	NorwichTerrier
	OldEnglishSheepdog
	Otterhound
	Papillon
	Pekinese
	Pembroke
	Pomeranian
	Pug
	Redbone
	RhodesianRidgeback
	Rottweiler
	SaintBernard
	Saluki
	Samoyed
	Schipperke
	ScotchTerrier
	ScottishDeerhound
	SealyhamTerrier
	ShetlandSheepdog
	ShihTzu
	SiberianHusky
	SilkyTerrier
	SoftCoatedWheatenTerrier
	StaffordshireBullterrier
	StandardPoodle
	StandardSchnauzer
	SussexSpaniel
	TibetanMastiff
	TibetanTerrier
	ToyPoodle
	ToyTerrier
	Vizsla
	WalkerHound
	Weimaraner
	WelshSpringerSpaniel
	WestHighlandWhiteTerrier
	Whippet
	WireHairedFoxTerrier
	YorkshireTerrier
)

const (
	UnknowSize Size = iota
	Small
	Medium
	Large
)
const (
	UnknowAge Age = iota
	Puppy
	Young
	Adult
)

const (
	UnknowColor CoatColor = iota
	Black
	Brown
	Gray
	Beige
	White
	Multicolor
	Brindle
)

const (
	UnknowCoatLength CoatLength = iota
	Short
	Long
)

const (
	UnknowTailLength TailLength = iota
	ShortTail
	LongTail
)

const (
	UnknowEar Ear = iota
	Standing
	Cut
	Floppy
)

const (
	Pending Ack = iota
	Accepted
	Rejected
)

type Dog struct {
	ID         int64
	Name       string
	Breed      Breed
	Age        Age
	Size       Size
	CoatColor  CoatColor
	CoatLength CoatLength
	TailLength TailLength
	Ear        Ear
	IsLost     bool
	Owner      *User
	Host       *User
	Latitude   float32
	Longitude  float32
	ImgUrl     string
	CreateAt   time.Time
}

type DogModel struct {
	ID         int64
	Name       string
	Breed      Breed
	Age        Age
	Size       Size
	CoatColor  CoatColor
	CoatLength CoatLength
	TailLength TailLength
	Ear        Ear
	IsLost     bool
	OwnerID    string
	HostID     string
	Latitude   float32
	Longitude  float32
	ImgUrl     string
	CreateAt   time.Time
	DeleteAt   time.Time
}

type Breed int

func (b Breed) String() string {
	switch b {
	case Affenpinscher:
		return "Affenpinscher"
	case AfghanHound:
		return "Afghan Hound"
	case AfricanHuntingDog:
		return "African Hunting Dog"
	case Airedale:
		return "Airedale"
	case AmericanStaffordshireTerrier:
		return "American Staffordshire Terrier"
	case Appenzeller:
		return "Appenzeller"
	case AustralianTerrier:
		return "Australian Terrier"
	case Basenji:
		return "Basenji"
	case Basset:
		return "Basset"
	case Beagle:
		return "Beagle"
	case BedlingtonTerrier:
		return "Bedlington Terrier"
	case BerneseMountainDog:
		return "Bernese Mountain Dog"
	case BlackAndTanCoonhound:
		return "Black-And-Tan Coonhound"
	case BlenheimSpaniel:
		return "Blenheim Spaniel"
	case Bloodhound:
		return "Bloodhound"
	case Bluetick:
		return "Bluetick"
	case BorderCollie:
		return "Border Collie"
	case BorderTerrier:
		return "Border Terrier"
	case Borzoi:
		return "Borzoi"
	case BostonBull:
		return "Boston Bull"
	case BouvierDesFlandres:
		return "Bouvier Des Flandres"
	case Boxer:
		return "Boxer"
	case BrabanconGriffon:
		return "Brabancon Griffon"
	case Briard:
		return "Briard"
	case BrittanySpaniel:
		return "Brittany Spaniel"
	case BullMastiff:
		return "Bull Mastiff"
	case Cairn:
		return "Cairn"
	case Cardigan:
		return "Cardigan"
	case ChesapeakeBayRetriever:
		return "Chesapeake Bay Retriever"
	case Chihuahua:
		return "Chihuahua"
	case Chow:
		return "Chow"
	case Clumber:
		return "Clumber"
	case CockerSpaniel:
		return "Cocker Spaniel"
	case Collie:
		return "Collie"
	case CurlyCoatedRetriever:
		return "Curly-Coated Retriever"
	case DandieDinmont:
		return "Dandie Dinmont"
	case Dhole:
		return "Dhole"
	case Dingo:
		return "Dingo"
	case Doberman:
		return "Doberman"
	case EnglishFoxhound:
		return "English Foxhound"
	case EnglishSetter:
		return "English Setter"
	case EnglishSpringer:
		return "English Springer"
	case Entlebucher:
		return "Entlebucher"
	case EskimoDog:
		return "Eskimo Dog"
	case FlatCoatedRetriever:
		return "Flat-Coated Retriever"
	case FrenchBulldog:
		return "French Bulldog"
	case GermanShepherd:
		return "German Shepherd"
	case GermanShortHairedPointer:
		return "German Short-Haired Pointer"
	case GiantSchnauzer:
		return "Giant Schnauzer"
	case GoldenRetriever:
		return "Golden Retriever"
	case GordonSetter:
		return "Gordon Setter"
	case GreatDane:
		return "Great Dane"
	case GreatPyrenees:
		return "Great Pyrenees"
	case GreaterSwissMountainDog:
		return "Greater Swiss Mountain Dog"
	case Groenendael:
		return "Groenendael"
	case IbizanHound:
		return "Ibizan Hound"
	case IrishSetter:
		return "Irish Setter"
	case IrishTerrier:
		return "Irish Terrier"
	case IrishWaterSpaniel:
		return "Irish Water Spaniel"
	case IrishWolfhound:
		return "Irish Wolfhound"
	case ItalianGreyhound:
		return "Italian Greyhound"
	case JapaneseSpaniel:
		return "Japanese Spaniel"
	case Keeshond:
		return "Keeshond"
	case Kelpie:
		return "Kelpie"
	case KerryBlueTerrier:
		return "Kerry Blue Terrier"
	case Komondor:
		return "Komondor"
	case Kuvasz:
		return "Kuvasz"
	case LabradorRetriever:
		return "Labrador Retriever"
	case LakelandTerrier:
		return "Lakeland Terrier"
	case Leonberg:
		return "Leonberg"
	case Lhasa:
		return "Lhasa"
	case Malamute:
		return "Malamute"
	case Malinois:
		return "Malinois"
	case MalteseDog:
		return "Maltese Dog"
	case MexicanHairless:
		return "Mexican Hairless"
	case MiniaturePinscher:
		return "Miniature Pinscher"
	case MiniaturePoodle:
		return "Miniature Poodle"
	case MiniatureSchnauzer:
		return "Miniature Schnauzer"
	case Newfoundland:
		return "Newfoundland"
	case NorfolkTerrier:
		return "Norfolk Terrier"
	case NorwegianElkhound:
		return "Norwegian Elkhound"
	case NorwichTerrier:
		return "Norwich Terrier"
	case OldEnglishSheepdog:
		return "Old English Sheepdog"
	case Otterhound:
		return "Otterhound"
	case Papillon:
		return "Papillon"
	case Pekinese:
		return "Pekinese"
	case Pembroke:
		return "Pembroke"
	case Pomeranian:
		return "Pomeranian"
	case Pug:
		return "Pug"
	case Redbone:
		return "Redbone"
	case RhodesianRidgeback:
		return "Rhodesian Ridgeback"
	case Rottweiler:
		return "Rottweiler"
	case SaintBernard:
		return "Saint Bernard"
	case Saluki:
		return "Saluki"
	case Samoyed:
		return "Samoyed"
	case Schipperke:
		return "Schipperke"
	case ScotchTerrier:
		return "Scotch Terrier"
	case ScottishDeerhound:
		return "Scottish Deerhound"
	case SealyhamTerrier:
		return "Sealyham Terrier"
	case ShetlandSheepdog:
		return "Shetland Sheepdog"
	case ShihTzu:
		return "Shih-Tzu"
	case SiberianHusky:
		return "Siberian Husky"
	case SilkyTerrier:
		return "Silky Terrier"
	case SoftCoatedWheatenTerrier:
		return "Soft-Coated Wheaten Terrier"
	case StaffordshireBullterrier:
		return "Staffordshire Bullterrier"
	case StandardPoodle:
		return "Standard Poodle"
	case StandardSchnauzer:
		return "Standard Schnauzer"
	case SussexSpaniel:
		return "Sussex Spaniel"
	case TibetanMastiff:
		return "Tibetan Mastiff"
	case TibetanTerrier:
		return "Tibetan Terrier"
	case ToyPoodle:
		return "Toy Poodle"
	case ToyTerrier:
		return "Toy Terrier"
	case Vizsla:
		return "Vizsla"
	case WalkerHound:
		return "Walker Hound"
	case Weimaraner:
		return "Weimaraner"
	case WelshSpringerSpaniel:
		return "Welsh Springer Spaniel"
	case WestHighlandWhiteTerrier:
		return "West Highland White Terrier"
	case Whippet:
		return "Whippet"
	case WireHairedFoxTerrier:
		return "Wire-Haired Fox Terrier"
	case YorkshireTerrier:
		return "Yorkshire Terrier"
	case Other:
		return "Otra"
	default:
		return "Desconocido"
	}
}

func ParseBreed(breed string) Breed {
	switch breed {
	case "Affenpinscher":
		return Affenpinscher
	case "Afghan Hound":
		return AfghanHound
	case "African Hunting Dog":
		return AfricanHuntingDog
	case "Airedale":
		return Airedale
	case "American Staffordshire Terrier":
		return AmericanStaffordshireTerrier
	case "Appenzeller":
		return Appenzeller
	case "Australian Terrier":
		return AustralianTerrier
	case "Basenji":
		return Basenji

	case "Basset":
		return Basset

	case "Beagle":
		return Beagle

	case "Bedlington Terrier":
		return BedlingtonTerrier

	case "Bernese Mountain Dog":
		return BerneseMountainDog

	case "Black-And-Tan Coonhound":
		return BlackAndTanCoonhound

	case "Blenheim Spaniel":
		return BlenheimSpaniel

	case "Bloodhound":
		return Bloodhound

	case "Bluetick":
		return Bluetick

	case "Border Collie":
		return BorderCollie

	case "Border Terrier":
		return BorderTerrier

	case "Borzoi":
		return Borzoi

	case "Boston Bull":
		return BostonBull

	case "Bouvier Des Flandres":
		return BouvierDesFlandres

	case "Boxer":
		return Boxer

	case "Brabancon Griffon":
		return BrabanconGriffon

	case "Briard":
		return Briard

	case "Brittany Spaniel":
		return BrittanySpaniel

	case "Bull Mastiff":
		return BullMastiff

	case "Cairn":
		return Cairn

	case "Cardigan":
		return Cardigan

	case "Chesapeake Bay Retriever":
		return ChesapeakeBayRetriever

	case "Chihuahua":
		return Chihuahua

	case "Chow":
		return Chow

	case "Clumber":
		return Clumber

	case "Cocker Spaniel":
		return CockerSpaniel

	case "Collie":
		return Collie

	case "Curly-Coated Retriever":
		return CurlyCoatedRetriever

	case "Dandie Dinmont":
		return DandieDinmont

	case "Dhole":
		return Dhole

	case "Dingo":
		return Dingo

	case "Doberman":
		return Doberman

	case "English Foxhound":
		return EnglishFoxhound

	case "English Setter":
		return EnglishSetter

	case "English Springer":
		return EnglishSpringer

	case "Entlebucher":
		return Entlebucher

	case "Eskimo Dog":
		return EskimoDog

	case "Flat-Coated Retriever":
		return FlatCoatedRetriever

	case "French Bulldog":
		return FrenchBulldog

	case "German Shepherd":
		return GermanShepherd

	case "German Short-Haired Pointer":
		return GermanShortHairedPointer

	case "Giant Schnauzer":
		return GiantSchnauzer

	case "Golden Retriever":
		return GoldenRetriever

	case "Gordon Setter":
		return GordonSetter

	case "Great Dane":
		return GreatDane

	case "Great Pyrenees":
		return GreatPyrenees

	case "Greater Swiss Mountain Dog":
		return GreaterSwissMountainDog

	case "Groenendael":
		return Groenendael

	case "Ibizan Hound":
		return IbizanHound

	case "Irish Setter":
		return IrishSetter

	case "Irish Terrier":
		return IrishTerrier

	case "Irish Water Spaniel":
		return IrishWaterSpaniel

	case "Irish Wolfhound":
		return IrishWolfhound

	case "Italian Greyhound":
		return ItalianGreyhound

	case "Japanese Spaniel":
		return JapaneseSpaniel

	case "Keeshond":
		return Keeshond

	case "Kelpie":
		return Kelpie

	case "Kerry Blue Terrier":
		return KerryBlueTerrier

	case "Komondor":
		return Komondor

	case "Kuvasz":
		return Kuvasz

	case "Labrador Retriever":
		return LabradorRetriever

	case "Lakeland Terrier":
		return LakelandTerrier

	case "Leonberg":
		return Leonberg

	case "Lhasa":
		return Lhasa

	case "Malamute":
		return Malamute

	case "Malinois":
		return Malinois

	case "Maltese Dog":
		return MalteseDog

	case "Mexican Hairless":
		return MexicanHairless

	case "Miniature Pinscher":
		return MiniaturePinscher

	case "Miniature Poodle":
		return MiniaturePoodle

	case "Miniature Schnauzer":
		return MiniatureSchnauzer

	case "Newfoundland":
		return Newfoundland

	case "Norfolk Terrier":
		return NorfolkTerrier

	case "Norwegian Elkhound":
		return NorwegianElkhound

	case "Norwich Terrier":
		return NorwichTerrier

	case "Old English Sheepdog":
		return OldEnglishSheepdog

	case "Otterhound":
		return Otterhound

	case "Papillon":
		return Papillon

	case "Pekinese":
		return Pekinese

	case "Pembroke":
		return Pembroke

	case "Pomeranian":
		return Pomeranian

	case "Pug":
		return Pug

	case "Redbone":
		return Redbone

	case "Rhodesian Ridgeback":
		return RhodesianRidgeback

	case "Rottweiler":
		return Rottweiler

	case "Saint Bernard":
		return SaintBernard

	case "Saluki":
		return Saluki

	case "Samoyed":
		return Samoyed

	case "Schipperke":
		return Schipperke

	case "Scotch Terrier":
		return ScotchTerrier

	case "Scottish Deerhound":
		return ScottishDeerhound

	case "Sealyham Terrier":
		return SealyhamTerrier

	case "Shetland Sheepdog":
		return ShetlandSheepdog

	case "Shih-Tzu":
		return ShihTzu

	case "Siberian Husky":
		return SiberianHusky

	case "Silky Terrier":
		return SilkyTerrier

	case "Soft-Coated Wheaten Terrier":
		return SoftCoatedWheatenTerrier

	case "Staffordshire Bullterrier":
		return StaffordshireBullterrier

	case "Standard Poodle":
		return StandardPoodle

	case "Standard Schnauzer":
		return StandardSchnauzer

	case "Sussex Spaniel":
		return SussexSpaniel

	case "Tibetan Mastiff":
		return TibetanMastiff

	case "Tibetan Terrier":
		return TibetanTerrier

	case "Toy Poodle":
		return ToyPoodle

	case "Toy Terrier":
		return ToyTerrier

	case "Vizsla":
		return Vizsla

	case "Walker Hound":
		return WalkerHound

	case "Weimaraner":
		return Weimaraner

	case "Welsh Springer Spaniel":
		return WelshSpringerSpaniel

	case "West Highland White Terrier":
		return WestHighlandWhiteTerrier

	case "Whippet":
		return Whippet

	case "Wire-Haired Fox Terrier":
		return WireHairedFoxTerrier

	case "Yorkshire Terrier":
		return YorkshireTerrier
	case "Otra":
		return Other
	default:
		return UnknowBreed
	}
}

type Size int

func (s Size) String() string {
	switch s {
	case Small:
		return "Chico"
	case Medium:
		return "Mediano"
	case Large:
		return "Grande"
	default:
		return "Desconocido"
	}
}

func ParseSize(size string) Size {
	switch strings.ToUpper(size) {
	case "CHICO":
		return Small
	case "MEDIANO":
		return Medium
	case "GRANDE":
		return Large
	default:
		return UnknowSize
	}
}

type Age int

func (a Age) String() string {
	switch a {
	case Puppy:
		return "Cachorro"
	case Young:
		return "Joven"
	case Adult:
		return "Adulto"
	default:
		return "Desconocido"
	}
}

func ParseAge(age string) Age {
	switch strings.ToUpper(age) {
	case "CACHORRO":
		return Puppy
	case "JOVEN":
		return Young
	case "ADULTO":
		return Adult
	default:
		return UnknowAge
	}
}

type CoatColor int

func (cc CoatColor) String() string {
	switch cc {
	case Black:
		return "Negro"
	case Brown:
		return "Marron"
	case Gray:
		return "Gris"
	case Beige:
		return "Beige"
	case White:
		return "Blanco"
	case Multicolor:
		return "Mas de un color"
	case Brindle:
		return "Atigrado"
	default:
		return "Desconocido"
	}
}

func ParseCoatColor(coatColor string) CoatColor {
	switch strings.ToUpper(coatColor) {
	case "NEGRO":
		return Black
	case "MARRON":
		return Brown
	case "GRIS":
		return Gray
	case "BEIGE":
		return Beige
	case "BLANCO":
		return White
	case "MAS DE UN COLOR":
		return Multicolor
	case "ATIGRADO":
		return Brindle
	default:
		return UnknowColor
	}
}

type CoatLength int

func (cl CoatLength) String() string {
	switch cl {
	case Short:
		return "Corto"
	case Long:
		return "Largo"
	default:
		return "Desconocido"
	}
}

func ParseCoatLength(length string) CoatLength {
	switch strings.ToUpper(length) {
	case "CORTO":
		return Short
	case "LARGO":
		return Long
	default:
		return UnknowCoatLength
	}
}

type TailLength int

func (tl TailLength) String() string {
	switch tl {
	case ShortTail:
		return "Corta"
	case LongTail:
		return "Larga"
	default:
		return "Desconocido"
	}
}

func ParseTailLength(length string) TailLength {
	switch strings.ToUpper(length) {
	case "CORTA":
		return ShortTail
	case "LARGA":
		return LongTail
	default:
		return UnknowTailLength
	}
}

type Ear int

func (e Ear) String() string {
	switch e {
	case Standing:
		return "Parada"
	case Cut:
		return "Cortada"
	case Floppy:
		return "Caida"
	default:
		return "Desconocido"
	}
}

func ParseEar(ear string) Ear {
	switch strings.ToUpper(ear) {
	case "PARADA":
		return Standing
	case "CORTADA":
		return Cut
	case "CAIDA":
		return Floppy
	default:
		return UnknowEar
	}
}

type DogResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Breed      string  `json:"breed"`
	Age        string  `json:"age"`
	Size       string  `json:"size"`
	CoatColor  string  `json:"coatColor"`
	CoatLength string  `json:"coatLength"`
	TailLength string  `json:"tailLength"`
	Ear        string  `json:"ear"`
	IsLost     bool    `json:"isLost"`
	Owner      string  `json:"owner"`
	Host       string  `json:"host"`
	Latitude   float32 `json:"latitude"`
	Longitude  float32 `json:"longitude"`
	ImgsUrl    string  `json:"imgsUrl"`
	ProfileImg string  `json:"profileImg"`
}

type DogRequest struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Breed      string   `json:"breed"`
	Age        string   `json:"age"`
	Size       string   `json:"size"`
	CoatColor  string   `json:"coatColor"`
	CoatLength string   `json:"coatLength"`
	TailLength string   `json:"tailLength"`
	Ear        string   `json:"ear"`
	IsLost     bool     `json:"isLost"`
	Owner      string   `json:"owner"`
	Host       string   `json:"host"`
	Latitude   float32  `json:"latitude"`
	Longitude  float32  `json:"longitude"`
	ImgUrl     string   `json:"imgUrl"`
	Imgs       []string `json:"imgs"`
}

type PossibleMatch struct {
	DogID         uint
	PossibleDogID uint
	Ack           Ack
}

type Ack int

func (a Ack) String() string {
	switch a {
	case Pending:
		return "PENDING"
	case Accepted:
		return "ACCEPTED"
	case Rejected:
		return "REJECTED"
	default:
		return "PENDING"
	}
}

func ParseAck(ack string) Ack {
	switch strings.ToUpper(ack) {
	case "ACCEPTED":
		return Accepted
	case "REJECTED":
		return Rejected
	case "PENDING":
		return Pending
	default:
		return Pending
	}
}
