package filter

import "reflect"

// And creates a new filter that passes if all the given filters pass
func And[T any](filters ...Filter[T]) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		for _, filter := range filters {
			if !filter.Apply(item) {
				return false
			}
		}
		return true
	})
}

// Or creates a new filter that passes if any of the given filters pass
func Or[T any](filters ...Filter[T]) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		for _, filter := range filters {
			if filter.Apply(item) {
				return true
			}
		}
		return false
	})
}

// Not creates a new filter that passes if the given filter does not pass
func Not[T any](filter Filter[T]) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		return !filter.Apply(item)
	})
}

// IsNil checks if a field is nil (works for pointers, slices, maps)
func IsNil[T any](fieldName string) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		switch fieldValue.Kind() {
		case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface, reflect.Chan, reflect.Func:
			return fieldValue.IsNil()
		default:
			return false
		}
	})
}

// IsNotNil checks if a field is not nil
func IsNotNil[T any](fieldName string) Filter[T] {
	return Not(IsNil[T](fieldName))
}

// IsZero checks if a field has its zero value
func IsZero[T any](fieldName string) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		return fieldValue.IsZero()
	})
}

// IsNotZero checks if a field does not have its zero value
func IsNotZero[T any](fieldName string) Filter[T] {
	return Not(IsZero[T](fieldName))
}
