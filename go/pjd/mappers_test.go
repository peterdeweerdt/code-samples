package pjd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type snakeTest struct {
	in  string
	out string
}

func TestSnakeCaseMapper(t *testing.T) {
	tests := []snakeTest{
		{"Customer", "customer"},
		{"CustomerID", "customer_id"},
		{"ID", "id"},
		{"HTTPRequest", "http_request"},
		{"TheHTTPRequest", "the_http_request"},
		{"CustomerName", "customer_name"},
		{"ImageURL", "image_url"},
		{"POSID", "posid"},
		{"PosID", "pos_id"},
		{"testCamel", "test_camel"},
	}

	for _, test := range tests {
		result := ToSnakeCase(test.in)
		assert.Equal(t, test.out, result)
	}
}
