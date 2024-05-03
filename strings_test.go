package sanitize

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrValidate(mainTest *testing.T) {
	one, two := 1, 2

	type testCase struct {
		name        string
		arg         StrValidator
		val         string
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "base case",
		},
		{
			name:        "force utf8 by default",
			val:         string([]byte{'\xc3', '\x28'}),
			expectedErr: &ErrNotUTF8{FieldName: "field"},
		},
		{
			name:        "errors if not passed anything when required",
			arg:         StrValidator{Required: true},
			expectedErr: &ErrRequired{FieldName: "field"},
		},
		{
			name: "minlength",
			arg:  StrValidator{Required: true, MinLen: 2},
			val:  "1",
			expectedErr: &ErrOutOfRange[int]{
				FieldName: "field",
				Value:     1,
				Min:       &two,
			},
		},
		{
			name: "minlength check doesn't happen if required is false",
			arg:  StrValidator{MinLen: 2},
		},
		{
			name: "maxlength",
			arg:  StrValidator{MaxLen: 1},
			val:  "asd",
			expectedErr: &ErrOutOfRange[int]{
				FieldName: "field",
				Value:     3,
				Max:       &one,
			},
		},
		{
			name: "allow non-utf8",
			arg:  StrValidator{AllowNonUTF8: true},
			val:  "퟿",
		},
		{
			name: "match regex",
			arg: StrValidator{
				Regex: regexp.MustCompile(`[0-9]`),
			},
			val: "a",
			expectedErr: &ErrFailedRegexp{
				FieldName: "field",
				Regexp:    "[0-9]",
			},
		},
	}

	t := assert.New(mainTest)
	for _, tc := range testCases {
		actualErr := tc.arg.Validate("field", tc.val)
		t.Equal(tc.expectedErr, actualErr, tc.name)
	}
}
