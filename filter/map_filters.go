package filter

import (
	"reflect"
)

// HasKey checks if a map field contains the specified key
func HasKey[T any](fieldName string, key interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		if fieldValue.Kind() != reflect.Map {
			return false
		}

		keyValue := reflect.ValueOf(key)

		// Make sure the key is of the correct type for the map
		if !keyValue.Type().AssignableTo(fieldValue.Type().Key()) {
			return false
		}

		return fieldValue.MapIndex(keyValue).IsValid()
	})
}

// HasValue checks if a map field contains the specified value
func HasValue[T any](fieldName string, value interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		if fieldValue.Kind() != reflect.Map {
			return false
		}

		targetValue := reflect.ValueOf(value)

		// Iterate through all map values
		iter := fieldValue.MapRange()
		for iter.Next() {
			mapValue := iter.Value()
			equal, err := compareValues(mapValue, targetValue)
			if err == nil && equal {
				return true
			}
		}

		return false
	})
}

// KeyValueEquals checks if a specific key in a map has a specific value
func KeyValueEquals[T any](fieldName string, key, value interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		if fieldValue.Kind() != reflect.Map {
			return false
		}

		keyValue := reflect.ValueOf(key)

		// Make sure the key is of the correct type for the map
		if !keyValue.Type().AssignableTo(fieldValue.Type().Key()) {
			return false
		}

		mapValueForKey := fieldValue.MapIndex(keyValue)
		if !mapValueForKey.IsValid() {
			return false
		}

		// Handle the case where mapValueForKey is an interface containing the actual value
		if mapValueForKey.Kind() == reflect.Interface {
			mapValueForKey = mapValueForKey.Elem()
		}

		targetValue := reflect.ValueOf(value)

		// If we have interface values, do a direct comparison of the underlying values
		if mapValueForKey.Kind() == reflect.Interface || targetValue.Kind() == reflect.Interface {
			return mapValueForKey.Interface() == targetValue.Interface()
		}

		// Otherwise use the general comparison function
		equal, err := compareValues(mapValueForKey, targetValue)
		if err != nil {
			return false
		}

		return equal
	})
}

// MapContainsAll checks if a map contains all the specified key-value pairs
func MapContainsAll[T any](fieldName string, kvPairs map[interface{}]interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		if fieldValue.Kind() != reflect.Map {
			return false
		}

		for k, v := range kvPairs {
			keyValue := reflect.ValueOf(k)

			// Make sure the key is of the correct type for the map
			if !keyValue.Type().AssignableTo(fieldValue.Type().Key()) {
				return false
			}

			mapValueForKey := fieldValue.MapIndex(keyValue)
			if !mapValueForKey.IsValid() {
				return false
			}

			targetValue := reflect.ValueOf(v)
			equal, err := compareValues(mapValueForKey, targetValue)
			if err != nil || !equal {
				return false
			}
		}

		return true
	})
}

// MapContainsAny checks if a map contains any of the specified key-value pairs
func MapContainsAny[T any](fieldName string, kvPairs map[interface{}]interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		if fieldValue.Kind() != reflect.Map {
			return false
		}

		for k, v := range kvPairs {
			keyValue := reflect.ValueOf(k)

			// Make sure the key is of the correct type for the map
			if !keyValue.Type().AssignableTo(fieldValue.Type().Key()) {
				continue
			}

			mapValueForKey := fieldValue.MapIndex(keyValue)
			if !mapValueForKey.IsValid() {
				continue
			}

			targetValue := reflect.ValueOf(v)
			equal, err := compareValues(mapValueForKey, targetValue)
			if err == nil && equal {
				return true
			}
		}

		return false
	})
}

// MapSizeEquals checks if a map has exactly the specified number of entries
func MapSizeEquals[T any](fieldName string, size int) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		if fieldValue.Kind() != reflect.Map {
			return false
		}

		return fieldValue.Len() == size
	})
}

// MapSizeGreaterThan checks if a map has more than the specified number of entries
func MapSizeGreaterThan[T any](fieldName string, size int) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		if fieldValue.Kind() != reflect.Map {
			return false
		}

		return fieldValue.Len() > size
	})
}

// MapSizeLessThan checks if a map has fewer than the specified number of entries
func MapSizeLessThan[T any](fieldName string, size int) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		if fieldValue.Kind() != reflect.Map {
			return false
		}

		return fieldValue.Len() < size
	})
}
