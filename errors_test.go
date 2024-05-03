package sanitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrRequiredIs(mainTest *testing.T) {
	testCases := []struct {
		name     string
		arg      error
		expected bool
	}{
		{name: "base case"},
		{name: "nil *ErrRequired", arg: (*ErrRequired)(nil), expected: true},
		{name: "*ErrRequired", arg: &ErrRequired{}, expected: true},
	}

	t := assert.New(mainTest)
	for _, tc := range testCases {
		actual := (&ErrRequired{}).Is(tc.arg)
		t.Equal(tc.expected, actual, tc.name)
	}
}

func TestErrOutOfRange(mainTest *testing.T) {
	one := 1
	two := 2

	testCases := []struct {
		name     string
		arg      ErrOutOfRange[int]
		expected string
	}{

		{
			name:     "base case",
			expected: "length of field must be within a valid range (that range was not specified)",
		},
		{
			name:     "issue: string/slice too big",
			arg:      ErrOutOfRange[int]{Max: &one, Value: 0},
			expected: "length of field must be less than 1, got 0within a valid range (that range was not specified)",
		},

		{
			name:     "issue: string/slice too small",
			arg:      ErrOutOfRange[int]{Min: &two, Value: 1},
			expected: "length of field must be greater than 2, got 1",
		},

		{
			name:     "issue: string/slice not between range",
			arg:      ErrOutOfRange[int]{Min: &one, Max: &two, Value: 3},
			expected: "length of field must be between 1 and 2, got 3",
		},
		{
			name: "issue: integer generic too big",
			arg: ErrOutOfRange[int]{
				Value: 2,
				Max:   &one,
				issue: 1,
			},
			expected: "value of field must be less than 1, got 2within a valid range (that range was not specified)",
		},
		{
			name: "issue: integer generic too small",
			arg: ErrOutOfRange[int]{
				Value: 1,
				Min:   &two,
				issue: 1,
			},
			expected: "value of field must be greater than 2, got 1",
		},
		{
			name: "issue: integer generic not between range",
			arg: ErrOutOfRange[int]{
				Value: 4,
				Min:   &one,
				Max:   &one,
				issue: 1,
			},
			expected: "value of field must be between 1 and 1, got 4",
		},
	}

	t := assert.New(mainTest)
	for _, tc := range testCases {
		tc.arg.FieldName = "field"
		actual := tc.arg.Error()
		t.Equal(tc.expected, actual, tc.name)
	}
}
