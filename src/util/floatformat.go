package util

import (
	"strconv"
)


func GetFloatFormat(num float64, pointNum int) (r float64) {
	f := strconv.FormatFloat(num, 'f', pointNum, 64)
	b,_ := strconv.ParseFloat(f, 64)

	return b
}