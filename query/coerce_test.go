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

func TestCoerceInt8(t *testing.T) {
	val, err := coerceValue("127", reflect.TypeOf(int8(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != int8(127) {
		t.Errorf("expected int8(127), got %v (%T)", val, val)
	}

	// Overflow
	_, err = coerceValue("200", reflect.TypeOf(int8(0)))
	if err == nil {
		t.Error("expected error for int8 overflow")
	}
}

func TestCoerceInt16(t *testing.T) {
	val, err := coerceValue("1000", reflect.TypeOf(int16(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != int16(1000) {
		t.Errorf("expected int16(1000), got %v (%T)", val, val)
	}

	_, err = coerceValue("abc", reflect.TypeOf(int16(0)))
	if err == nil {
		t.Error("expected error for invalid int16")
	}
}

func TestCoerceInt32(t *testing.T) {
	val, err := coerceValue("100000", reflect.TypeOf(int32(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != int32(100000) {
		t.Errorf("expected int32(100000), got %v (%T)", val, val)
	}

	_, err = coerceValue("abc", reflect.TypeOf(int32(0)))
	if err == nil {
		t.Error("expected error for invalid int32")
	}
}

func TestCoerceUint8(t *testing.T) {
	val, err := coerceValue("255", reflect.TypeOf(uint8(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != uint8(255) {
		t.Errorf("expected uint8(255), got %v (%T)", val, val)
	}

	_, err = coerceValue("300", reflect.TypeOf(uint8(0)))
	if err == nil {
		t.Error("expected error for uint8 overflow")
	}
}

func TestCoerceUint16(t *testing.T) {
	val, err := coerceValue("65535", reflect.TypeOf(uint16(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != uint16(65535) {
		t.Errorf("expected uint16(65535), got %v (%T)", val, val)
	}

	_, err = coerceValue("abc", reflect.TypeOf(uint16(0)))
	if err == nil {
		t.Error("expected error for invalid uint16")
	}
}

func TestCoerceUint32(t *testing.T) {
	val, err := coerceValue("4294967295", reflect.TypeOf(uint32(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != uint32(4294967295) {
		t.Errorf("expected uint32(4294967295), got %v (%T)", val, val)
	}

	_, err = coerceValue("abc", reflect.TypeOf(uint32(0)))
	if err == nil {
		t.Error("expected error for invalid uint32")
	}
}

func TestCoerceUint64(t *testing.T) {
	val, err := coerceValue("18446744073709551615", reflect.TypeOf(uint64(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != uint64(18446744073709551615) {
		t.Errorf("expected uint64(max), got %v (%T)", val, val)
	}

	_, err = coerceValue("abc", reflect.TypeOf(uint64(0)))
	if err == nil {
		t.Error("expected error for invalid uint64")
	}
}

func TestCoerceInvalidFloat32(t *testing.T) {
	_, err := coerceValue("abc", reflect.TypeOf(float32(0)))
	if err == nil {
		t.Error("expected error for invalid float32")
	}
}

func TestCoerceInvalidFloat64(t *testing.T) {
	_, err := coerceValue("abc", reflect.TypeOf(float64(0)))
	if err == nil {
		t.Error("expected error for invalid float64")
	}
}

func TestCoerceInvalidBool(t *testing.T) {
	_, err := coerceValue("maybe", reflect.TypeOf(false))
	if err == nil {
		t.Error("expected error for invalid bool")
	}
}

func TestCoerceInvalidTime(t *testing.T) {
	_, err := coerceValue("not-a-date", reflect.TypeOf(time.Time{}))
	if err == nil {
		t.Error("expected error for invalid time")
	}
}

func TestCoerceUnsupportedType(t *testing.T) {
	_, err := coerceValue("test", reflect.TypeOf(struct{}{}))
	if err == nil {
		t.Error("expected error for unsupported type")
	}
}

func TestCoerceInvalidUint(t *testing.T) {
	_, err := coerceValue("-1", reflect.TypeOf(uint(0)))
	if err == nil {
		t.Error("expected error for negative uint")
	}
}

func TestCoerceInvalidInt64(t *testing.T) {
	_, err := coerceValue("abc", reflect.TypeOf(int64(0)))
	if err == nil {
		t.Error("expected error for invalid int64")
	}
}
