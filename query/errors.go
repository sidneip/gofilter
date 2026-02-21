package query

import "fmt"

// ErrFieldNotFilterable is returned when a query attempts to filter
// on a field that does not have the "filterable" tag.
type ErrFieldNotFilterable struct{ Field string }

func (e *ErrFieldNotFilterable) Error() string {
	return fmt.Sprintf("field %q is not filterable", e.Field)
}

// ErrFieldNotSortable is returned when a query attempts to sort
// by a field that does not have the "sortable" tag.
type ErrFieldNotSortable struct{ Field string }

func (e *ErrFieldNotSortable) Error() string {
	return fmt.Sprintf("field %q is not sortable", e.Field)
}

// ErrInvalidValue is returned when a query parameter value cannot
// be coerced to the expected field type.
type ErrInvalidValue struct{ Field, Value, ExpectedType string }

func (e *ErrInvalidValue) Error() string {
	return fmt.Sprintf("invalid value %q for field %q: expected %s", e.Value, e.Field, e.ExpectedType)
}

// ErrLimitExceeded is returned when the requested pagination limit
// exceeds the maximum allowed by WithMaxLimit.
type ErrLimitExceeded struct{ Requested, Max int }

func (e *ErrLimitExceeded) Error() string {
	return fmt.Sprintf("requested limit %d exceeds maximum %d", e.Requested, e.Max)
}
