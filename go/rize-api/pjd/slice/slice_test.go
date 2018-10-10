package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	cases := []struct {
		slice       []string
		value       string
		expectation bool
	}{
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{"a", "b", "c"}, "A", false},
	}

	for _, c := range cases {
		assert.Equal(t, c.expectation, ContainsString(c.slice, c.value))
	}
}
