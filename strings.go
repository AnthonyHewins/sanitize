package sanitize

import (
	"regexp"
	"unicode/utf8"
)

type InvalidReasonStr byte

const (
	NotValidUTF8 InvalidReasonStr = iota
	DidNotMatchRegex
)

type ErrFailedRegexp struct {
	FieldName string
	Regexp    string
}

func (e *ErrFailedRegexp) Error() string {
	return e.FieldName + " did not match required pattern"
}

func (e *ErrFailedRegexp) Is(err error) bool {
	_, ok := err.(*ErrFailedRegexp)
	return ok
}

type ErrNotUTF8 struct {
	FieldName string
}

func (e *ErrNotUTF8) Error() string {
	return e.FieldName + " contained invalid UTF-8 characters"
}

func (e *ErrNotUTF8) Is(err error) bool {
	_, ok := err.(*ErrNotUTF8)
	return ok
}

type StrValidator struct {
	// Force the string to be:
	// - Non-empty if you're validating a string
	// - Non-nil if you're validating a *string
	Required bool

	// Length checks. If you set max length to anything other than zero,
	// it will check for a max length. It always checks for min length,
	// since it defaults to 0, which will always pass if unset
	MinLen, MaxLen int

	// Force UTF8, unless specified otherwise
	AllowNonUTF8 bool

	// Use this regex to validate
	Regex *regexp.Regexp
}

func (s *StrValidator) ValidatePtr(fieldName string, val *string) error {
	if val == nil {
		if !s.Required {
			return nil
		}

		return &ErrRequired{FieldName: fieldName}
	}

	return s.coreValidate(fieldName, *val)
}

func (s *StrValidator) Validate(fieldName string, val string) error {
	if val == "" {
		if !s.Required {
			return nil
		}

		return &ErrRequired{FieldName: fieldName}
	}

	return s.coreValidate(fieldName, val)
}

func (s *StrValidator) coreValidate(fieldName string, val string) error {
	n := len(val)
	if max := s.MaxLen; (max != 0 && n > max) || n < s.MinLen {
		return s.rangeErr(fieldName, n)
	}

	if !s.AllowNonUTF8 && !utf8.ValidString(val) {
		return &ErrNotUTF8{FieldName: fieldName}
	}

	if s.Regex != nil && !s.Regex.MatchString(val) {
		return &ErrFailedRegexp{
			FieldName: fieldName,
			Regexp:    s.Regex.String(),
		}
	}

	return nil
}

func (s *StrValidator) rangeErr(fieldName string, n int) *ErrOutOfRange[int] {
	var max, min *int
	if s.MaxLen != 0 {
		max = &s.MaxLen
	}

	if s.MinLen != 0 {
		min = &s.MinLen
	}

	return &ErrOutOfRange[int]{
		FieldName: fieldName,
		Value:     n,
		Max:       max,
		Min:       min,
		issue:     length,
	}
}
