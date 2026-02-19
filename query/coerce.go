package query

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

func coerceValue(raw string, targetType reflect.Type) (interface{}, error) {
	if targetType == timeType {
		return parseTime(raw)
	}

	switch targetType.Kind() {
	case reflect.String:
		return raw, nil
	case reflect.Int:
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as int: %w", raw, err)
		}
		return int(v), nil
	case reflect.Int8:
		v, err := strconv.ParseInt(raw, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as int8: %w", raw, err)
		}
		return int8(v), nil
	case reflect.Int16:
		v, err := strconv.ParseInt(raw, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as int16: %w", raw, err)
		}
		return int16(v), nil
	case reflect.Int32:
		v, err := strconv.ParseInt(raw, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as int32: %w", raw, err)
		}
		return int32(v), nil
	case reflect.Int64:
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as int64: %w", raw, err)
		}
		return v, nil
	case reflect.Uint:
		v, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as uint: %w", raw, err)
		}
		return uint(v), nil
	case reflect.Uint8:
		v, err := strconv.ParseUint(raw, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as uint8: %w", raw, err)
		}
		return uint8(v), nil
	case reflect.Uint16:
		v, err := strconv.ParseUint(raw, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as uint16: %w", raw, err)
		}
		return uint16(v), nil
	case reflect.Uint32:
		v, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as uint32: %w", raw, err)
		}
		return uint32(v), nil
	case reflect.Uint64:
		v, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as uint64: %w", raw, err)
		}
		return v, nil
	case reflect.Float32:
		v, err := strconv.ParseFloat(raw, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as float32: %w", raw, err)
		}
		return float32(v), nil
	case reflect.Float64:
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as float64: %w", raw, err)
		}
		return v, nil
	case reflect.Bool:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as bool: %w", raw, err)
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", targetType)
	}
}

func parseTime(raw string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02",
	}
	for _, f := range formats {
		t, err := time.Parse(f, raw)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse %q as time", raw)
}
