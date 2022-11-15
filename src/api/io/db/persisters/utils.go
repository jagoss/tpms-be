package persisters

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func IsRecordNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func ToFloat64List(embedding string) []float64 {
	sliceEmbeddig := strings.Split(embedding, " ")
	resultList := make([]float64, len(sliceEmbeddig))
	for _, stringEmb := range sliceEmbeddig {
		nbr, _ := strconv.ParseFloat(stringEmb, 64)
		resultList = append(resultList, nbr)
	}
	return resultList
}
