package pjd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMappingStringsToInt64s(t *testing.T) {
	cases := []struct {
		in            []string
		out           []int64
		shouldSucceed bool
	}{
		{[]string{"1", "2", "3", "4"}, []int64{1, 2, 3, 4}, true},
		{[]string{"1"}, []int64{1}, true},
		{[]string{}, []int64{}, true},
		{[]string{"1.3"}, nil, false},
		{[]string{""}, nil, false},
		{[]string{"x"}, nil, false},
	}

	for _, c := range cases {
		result, err := MapStringsToInt64s(c.in)
		if c.shouldSucceed {
			assert.NoError(t, err)
			assert.Equal(t, c.out, result)
		} else {
			assert.Error(t, err)
			assert.Equal(t, c.out, result)
		}
	}
}
