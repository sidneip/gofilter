package filter

import (
	"fmt"
	"reflect"
	"strings"
)

// getFieldValue gets the value of a field from a struct by name
// Supports nested fields with dot notation (e.g., "Address.City")
func getFieldValue(item interface{}, fieldPath string) (reflect.Value, error) {
	value := reflect.ValueOf(item)

	// Handle pointers
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return reflect.Value{}, fmt.Errorf("nil pointer")
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("item is not a struct")
	}

	// Handle nested fields with dot notation
	fields := strings.Split(fieldPath, ".")
	for i, field := range fields {
		value = value.FieldByName(field)

		if !value.IsValid() {
			return reflect.Value{}, fmt.Errorf("field %s not found", field)
		}

		// Handle pointer to struct for nested fields
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				return reflect.Value{}, fmt.Errorf("nil pointer for field %s", field)
			}
			value = value.Elem()
		}

		// If not the last field, ensure it's a struct
		if i < len(fields)-1 && value.Kind() != reflect.Struct {
			return reflect.Value{}, fmt.Errorf("%s is not a struct", field)
		}
	}

	return value, nil
}

// compareValues compares two values and returns true if they are equal
func compareValues(a, b reflect.Value) (bool, error) {
	// Handle different types
	if a.Type() != b.Type() {
		// Try to convert b to a's type
		if b.Type().ConvertibleTo(a.Type()) {
			b = b.Convert(a.Type())
		} else {
			return false, fmt.Errorf("cannot compare values of different types: %s and %s", a.Type(), b.Type())
		}
	}

	// Compare based on kind
	switch a.Kind() {
	case reflect.String:
		return a.String() == b.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() == b.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return a.Uint() == b.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return a.Float() == b.Float(), nil
	case reflect.Bool:
		return a.Bool() == b.Bool(), nil
	// Add more cases as needed
	default:
		return false, fmt.Errorf("unsupported type for comparison: %s", a.Type())
	}
}

// compareValuesLess compares two values and returns true if a < b
func compareValuesLess(a, b reflect.Value) (bool, error) {
	// Convert if needed
	if a.Type() != b.Type() {
		if b.Type().ConvertibleTo(a.Type()) {
			b = b.Convert(a.Type())
		} else {
			return false, fmt.Errorf("cannot compare values of different types: %s and %s", a.Type(), b.Type())
		}
	}

	// Compare based on kind
	switch a.Kind() {
	case reflect.String:
		return a.String() < b.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() < b.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return a.Uint() < b.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return a.Float() < b.Float(), nil
	default:
		return false, fmt.Errorf("unsupported type for less than comparison: %s", a.Type())
	}
}

// Eq returns a filter that checks if a field equals a value
func Eq[T any](fieldName string, value interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		targetValue := reflect.ValueOf(value)
		equal, err := compareValues(fieldValue, targetValue)
		if err != nil {
			return false
		}

		return equal
	})
}

// Ne returns a filter that checks if a field does not equal a value
func Ne[T any](fieldName string, value interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		targetValue := reflect.ValueOf(value)
		equal, err := compareValues(fieldValue, targetValue)
		if err != nil {
			return false
		}

		return !equal
	})
}

// Gt returns a filter that checks if a field is greater than a value
func Gt[T any](fieldName string, value interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		targetValue := reflect.ValueOf(value)
		less, err := compareValuesLess(targetValue, fieldValue)
		if err != nil {
			return false
		}

		return less
	})
}

// Lt returns a filter that checks if a field is less than a value
func Lt[T any](fieldName string, value interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		targetValue := reflect.ValueOf(value)
		less, err := compareValuesLess(fieldValue, targetValue)
		if err != nil {
			return false
		}

		return less
	})
}

// Gte returns a filter that checks if a field is greater than or equal to a value
func Gte[T any](fieldName string, value interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		targetValue := reflect.ValueOf(value)

		less, err := compareValuesLess(fieldValue, targetValue)
		if err != nil {
			return false
		}

		equal, err := compareValues(fieldValue, targetValue)
		if err != nil {
			return false
		}

		return !less || equal
	})
}

// Lte returns a filter that checks if a field is less than or equal to a value
func Lte[T any](fieldName string, value interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		targetValue := reflect.ValueOf(value)

		less, err := compareValuesLess(fieldValue, targetValue)
		if err != nil {
			return false
		}

		equal, err := compareValues(fieldValue, targetValue)
		if err != nil {
			return false
		}

		return less || equal
	})
}

// Contains returns a filter that checks if a field contains a value
// Works for strings, slices, and arrays
func Contains[T any](fieldName string, value interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		// For strings
		if fieldValue.Kind() == reflect.String {
			if valueStr, ok := value.(string); ok {
				return strings.Contains(fieldValue.String(), valueStr)
			}
			return false
		}

		// For slices and arrays
		if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
			targetValue := reflect.ValueOf(value)

			for i := 0; i < fieldValue.Len(); i++ {
				elem := fieldValue.Index(i)
				equal, err := compareValues(elem, targetValue)
				if err == nil && equal {
					return true
				}
			}
		}

		return false
	})
}

// In returns a filter that checks if a field's value is in a slice of values
func In[T any](fieldName string, values []interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		for _, value := range values {
			targetValue := reflect.ValueOf(value)
			equal, err := compareValues(fieldValue, targetValue)
			if err == nil && equal {
				return true
			}
		}

		return false
	})
}
