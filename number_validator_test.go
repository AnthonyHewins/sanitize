package sanitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberValidate(mainTest *testing.T) {
	zero, one, two := 0, 1, 2
	testCases := []struct {
		name        string
		v           NumberValidator[int]
		arg         *int
		expectedErr error
	}{
		{
			name: "not required -> no error on nil",
		},
		{
			name:        "required -> err on nil",
			v:           NumberValidator[int]{Required: true},
			expectedErr: &ErrRequired{FieldName: "field"},
		},
		{
			name: "too small -> err",
			v: NumberValidator[int]{
				Required: true,
				MinVal:   &one,
			},
			arg: new(int),
			expectedErr: &ErrOutOfRange[int]{
				FieldName: "field",
				issue:     1,
				Min:       &one,
			},
		},
		{
			name: "too big -> err",
			v: NumberValidator[int]{
				Required: true,
				MaxVal:   new(int),
			},
			arg: &one,
			expectedErr: &ErrOutOfRange[int]{
				FieldName: "field",
				Value:     1,
				Max:       &zero,
				issue:     1,
			},
		},
		{
			name: "valid",
			v: NumberValidator[int]{
				Required: true,
				MinVal:   new(int),
				MaxVal:   &two,
			},
			arg: &one,
		},
	}

	t := assert.New(mainTest)
	for _, tc := range testCases {
		actualErr := tc.v.ValidatePtr("field", tc.arg)
		t.Equal(tc.expectedErr, actualErr, tc.name)
	}
}
