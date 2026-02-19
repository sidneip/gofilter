package query

import (
	"reflect"
	"testing"
	"time"
)

func TestCoerceString(t *testing.T) {
	val, err := coerceValue("hello", reflect.TypeOf(""))
	if err != nil {
		t.Fatal(err)
	}
	if val != "hello" {
		t.Errorf("expected 'hello', got %v", val)
	}
}

func TestCoerceInt(t *testing.T) {
	val, err := coerceValue("42", reflect.TypeOf(0))
	if err != nil {
		t.Fatal(err)
	}
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestCoerceInt64(t *testing.T) {
	val, err := coerceValue("100", reflect.TypeOf(int64(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != int64(100) {
		t.Errorf("expected int64(100), got %v (%T)", val, val)
	}
}

func TestCoerceUint(t *testing.T) {
	val, err := coerceValue("10", reflect.TypeOf(uint(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != uint(10) {
		t.Errorf("expected uint(10), got %v (%T)", val, val)
	}
}

func TestCoerceFloat64(t *testing.T) {
	val, err := coerceValue("3.14", reflect.TypeOf(float64(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != 3.14 {
		t.Errorf("expected 3.14, got %v", val)
	}
}

func TestCoerceBool(t *testing.T) {
	val, err := coerceValue("true", reflect.TypeOf(false))
	if err != nil {
		t.Fatal(err)
	}
	if val != true {
		t.Errorf("expected true, got %v", val)
	}
}

func TestCoerceTimeRFC3339(t *testing.T) {
	val, err := coerceValue("2024-01-15T10:30:00Z", reflect.TypeOf(time.Time{}))
	if err != nil {
		t.Fatal(err)
	}
	tm, ok := val.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time, got %T", val)
	}
	if tm.Year() != 2024 || tm.Month() != 1 || tm.Day() != 15 {
		t.Errorf("unexpected time: %v", tm)
	}
}

func TestCoerceTimeDateOnly(t *testing.T) {
	val, err := coerceValue("2024-01-15", reflect.TypeOf(time.Time{}))
	if err != nil {
		t.Fatal(err)
	}
	tm, ok := val.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time, got %T", val)
	}
	if tm.Year() != 2024 || tm.Month() != 1 || tm.Day() != 15 {
		t.Errorf("unexpected time: %v", tm)
	}
}

func TestCoerceInvalidInt(t *testing.T) {
	_, err := coerceValue("abc", reflect.TypeOf(0))
	if err == nil {
		t.Error("expected error for invalid int")
	}
}

func TestCoerceFloat32(t *testing.T) {
	val, err := coerceValue("2.5", reflect.TypeOf(float32(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != float32(2.5) {
		t.Errorf("expected float32(2.5), got %v (%T)", val, val)
	}
}
