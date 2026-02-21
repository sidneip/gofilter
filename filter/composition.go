package filter

import "reflect"

// And creates a composite filter that passes only if ALL given filters pass.
// Use this to combine multiple conditions with logical AND.
//
// Example:
//
//	filter.And[User](
//	    filter.Gte[User]("Age", 18),
//	    filter.Eq[User]("Active", true),
//	)  // users who are 18+ AND active
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

// Or creates a composite filter that passes if ANY of the given filters pass.
// Use this to combine multiple conditions with logical OR.
//
// Example:
//
//	filter.Or[User](
//	    filter.Eq[User]("City", "SP"),
//	    filter.Eq[User]("City", "RJ"),
//	)  // users from SP OR RJ
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

// Not creates a filter that inverts the result of the given filter.
// Use this to negate a condition.
//
// Example:
//
//	filter.Not(filter.Eq[User]("Status", "banned"))  // users who are NOT banned
func Not[T any](filter Filter[T]) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		return !filter.Apply(item)
	})
}

// IsNil returns a filter that checks if a field is nil.
// Works with pointers, slices, maps, interfaces, channels, and functions.
//
// Example:
//
//	filter.IsNil[User]("DeletedAt")  // users where DeletedAt is nil
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

// IsNotNil returns a filter that checks if a field is not nil.
//
// Example:
//
//	filter.IsNotNil[User]("Avatar")  // users who have an avatar
func IsNotNil[T any](fieldName string) Filter[T] {
	return Not(IsNil[T](fieldName))
}

// IsZero returns a filter that checks if a field has its zero value.
// Zero values are: 0 for numbers, "" for strings, false for bools, nil for pointers, etc.
//
// Example:
//
//	filter.IsZero[User]("Score")  // users with Score == 0
func IsZero[T any](fieldName string) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		return fieldValue.IsZero()
	})
}

// IsNotZero returns a filter that checks if a field has a non-zero value.
//
// Example:
//
//	filter.IsNotZero[User]("Score")  // users with Score != 0
func IsNotZero[T any](fieldName string) Filter[T] {
	return Not(IsZero[T](fieldName))
}
