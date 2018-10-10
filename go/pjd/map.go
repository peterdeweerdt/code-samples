package pjd

import (
	"strconv"
)

func MapStringsToInt64s(strings []string) ([]int64, error) {
	ints := make([]int64, len(strings))
	for i, s := range strings {
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		ints[i] = val
	}
	return ints, nil
}
