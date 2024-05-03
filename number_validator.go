package sanitize

import (
	"golang.org/x/exp/constraints"
)

type NumberValidator[X constraints.Integer | constraints.Float] struct {
	// Require that the int value is passed.
	// This only can be enforced with ValidatePtr
	Required bool

	MaxVal *X // Max value allowed, if any
	MinVal *X // Min value allowed, if any
}

func (i *NumberValidator[X]) Validate(fieldName string, val X) error {
	return i.coreValidate(fieldName, val)
}

func (i *NumberValidator[X]) ValidatePtr(fieldName string, val *X) error {
	if val == nil {
		if !i.Required {
			return nil
		}

		return &ErrRequired{FieldName: fieldName}
	}

	return i.coreValidate(fieldName, *val)
}

func (i *NumberValidator[X]) coreValidate(fieldName string, deref X) error {
	if min := i.MinVal; min != nil && deref < *min {
		return i.rangeErr(fieldName, deref)
	}

	if max := i.MaxVal; max != nil && deref > *max {
		return i.rangeErr(fieldName, deref)
	}

	return nil
}

func (i *NumberValidator[X]) rangeErr(fieldName string, val X) *ErrOutOfRange[X] {
	return &ErrOutOfRange[X]{
		FieldName: fieldName,
		Value:     val,
		Min:       i.MinVal,
		Max:       i.MaxVal,
		issue:     number,
	}
}
