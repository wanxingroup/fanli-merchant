package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckMobile(t *testing.T) {

	tests := []struct {
		input string
		want  bool
	}{
		{
			input: "18916039393",
			want:  true,
		},
		{
			input: "13902212597",
			want:  true,
		},
		{
			input: "18601614130",
			want:  true,
		},
		{
			input: "13221024759",
			want:  true,
		},
		{
			input: "17181234567",
			want:  true,
		},
		{
			input: "18621073521",
			want:  true,
		},
		{
			input: "12345678901",
			want:  false,
		},
		{
			input: "1234567890",
			want:  false,
		},
		{
			input: "123456789012",
			want:  false,
		},
	}

	for _, test := range tests {

		assert.Equal(t, test.want, CheckMobile(test.input), test)
	}
}
