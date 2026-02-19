package query

import "fmt"

type ErrFieldNotFilterable struct{ Field string }

func (e *ErrFieldNotFilterable) Error() string {
	return fmt.Sprintf("field %q is not filterable", e.Field)
}

type ErrFieldNotSortable struct{ Field string }

func (e *ErrFieldNotSortable) Error() string {
	return fmt.Sprintf("field %q is not sortable", e.Field)
}

type ErrInvalidValue struct{ Field, Value, ExpectedType string }

func (e *ErrInvalidValue) Error() string {
	return fmt.Sprintf("invalid value %q for field %q: expected %s", e.Value, e.Field, e.ExpectedType)
}

type ErrLimitExceeded struct{ Requested, Max int }

func (e *ErrLimitExceeded) Error() string {
	return fmt.Sprintf("requested limit %d exceeds maximum %d", e.Requested, e.Max)
}
