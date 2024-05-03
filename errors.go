package sanitize

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

// ErrRequired appears when a field has a required value, meaning
// it should not be nil or should not be the zero value, but was
type ErrRequired struct {
	FieldName string
}

func (e *ErrRequired) Error() string {
	return e.FieldName + " is required"
}

func (e *ErrRequired) Is(err error) bool {
	_, ok := err.(*ErrRequired)
	return ok
}

type issue byte

const (
	length issue = iota
	number
)

// ErrOutOfRange is an error when a slice/string is too long/short,
// or when an integer/float value is too big/small
type ErrOutOfRange[X constraints.Float | constraints.Integer] struct {
	FieldName string
	Value     X
	Min       *X
	Max       *X
	issue
}

func (e *ErrOutOfRange[X]) Error() string {
	var sb strings.Builder
	switch e.issue {
	case length:
		sb.WriteString("length of ")
	case number:
		sb.WriteString("value of ")
	}

	sb.WriteString(e.FieldName + " must be ")

	if min := e.Min; min != nil {
		if max := e.Max; max != nil {
			sb.WriteString(fmt.Sprintf("between %v and %v, got %v", *min, *max, e.Value))
			return sb.String()
		}

		sb.WriteString(fmt.Sprintf("greater than %v, got %v", *min, e.Value))
		return sb.String()
	}

	if max := e.Max; max != nil {
		sb.WriteString(fmt.Sprintf("less than %v, got %v", *max, e.Value))
	}

	sb.WriteString("within a valid range (that range was not specified)")
	return sb.String()
}

func (e *ErrOutOfRange[X]) Is(err error) bool {
	_, ok := err.(*ErrOutOfRange[X])
	return ok
}
