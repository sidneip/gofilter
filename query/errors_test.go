package query

import (
	"errors"
	"testing"
)

func TestErrFieldNotFilterable(t *testing.T) {
	err := &ErrFieldNotFilterable{Field: "SSN"}
	if err.Error() != `field "SSN" is not filterable` {
		t.Errorf("unexpected error message: %s", err.Error())
	}
	var target *ErrFieldNotFilterable
	if !errors.As(err, &target) {
		t.Error("should be assertable with errors.As")
	}
}

func TestErrFieldNotSortable(t *testing.T) {
	err := &ErrFieldNotSortable{Field: "Email"}
	if err.Error() != `field "Email" is not sortable` {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestErrInvalidValue(t *testing.T) {
	err := &ErrInvalidValue{Field: "Age", Value: "abc", ExpectedType: "int"}
	if err.Error() != `invalid value "abc" for field "Age": expected int` {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestErrLimitExceeded(t *testing.T) {
	err := &ErrLimitExceeded{Requested: 500, Max: 100}
	if err.Error() != "requested limit 500 exceeds maximum 100" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
