package filter

import (
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
)

// ExportedGetFieldValue retrieves a field value from a struct by name.
// Supports nested fields using dot notation (e.g., "Address.City").
// This function is exported for use in custom filter implementations.
//
// Example:
//
//	val, err := filter.ExportedGetFieldValue(user, "Address.City")
//	if err == nil {
//	    fmt.Println(val.String())  // "São Paulo"
//	}
func ExportedGetFieldValue(item interface{}, fieldPath string) (reflect.Value, error) {
	return getFieldValue(item, fieldPath)
}

// StringMatchMode defines different modes for string matching
type StringMatchMode int

const (
	// ExactMatch requires the string to match exactly (with optional case insensitivity)
	ExactMatch StringMatchMode = iota
	// ContainsMatch checks if the field contains the search string
	ContainsMatch
	// PrefixMatch checks if the field starts with the search string
	PrefixMatch
	// SuffixMatch checks if the field ends with the search string
	SuffixMatch
)

// StringMatchOptions configures how string matching is performed
type StringMatchOptions struct {
	// Mode determines the type of string matching to perform
	Mode StringMatchMode
	// IgnoreCase determines whether the matching should be case-insensitive
	IgnoreCase bool
}

// StringMatch returns a filter with configurable string matching behavior.
// Supports exact match, contains, prefix, and suffix modes with optional case insensitivity.
//
// Example:
//
//	filter.StringMatch[User]("Email", "gmail.com", filter.StringMatchOptions{
//	    Mode: filter.SuffixMatch, IgnoreCase: true,
//	})  // users with email ending in "gmail.com" (case-insensitive)
func StringMatch[T any](fieldName string, value string, options StringMatchOptions) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil || fieldValue.Kind() != reflect.String {
			return false
		}

		fieldStr := fieldValue.String()

		switch options.Mode {
		case ExactMatch:
			if options.IgnoreCase {
				return strings.EqualFold(fieldStr, value)
			}
			return fieldStr == value
		case ContainsMatch:
			if options.IgnoreCase {
				return strings.Contains(strings.ToLower(fieldStr), strings.ToLower(value))
			}
			return strings.Contains(fieldStr, value)
		case PrefixMatch:
			if options.IgnoreCase {
				return strings.HasPrefix(strings.ToLower(fieldStr), strings.ToLower(value))
			}
			return strings.HasPrefix(fieldStr, value)
		case SuffixMatch:
			if options.IgnoreCase {
				return strings.HasSuffix(strings.ToLower(fieldStr), strings.ToLower(value))
			}
			return strings.HasSuffix(fieldStr, value)
		default:
			return false
		}
	})
}

// ArrayContains checks if an array field contains a specific value
func ArrayContains[T any](fieldName string, value interface{}, ignoreCase bool) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		if fieldValue.Kind() != reflect.Slice && fieldValue.Kind() != reflect.Array {
			return false
		}

		targetValue := reflect.ValueOf(value)

		for i := 0; i < fieldValue.Len(); i++ {
			elemValue := fieldValue.Index(i)

			if ignoreCase && elemValue.Kind() == reflect.String && targetValue.Kind() == reflect.String {
				if strings.EqualFold(elemValue.String(), targetValue.String()) {
					return true
				}
			} else {
				equal, err := compareValues(elemValue, targetValue)
				if err == nil && equal {
					return true
				}
			}
		}

		return false
	})
}

// ArrayContainsAny checks if an array contains any of the provided values
func ArrayContainsAny[T any](fieldName string, values []interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		if fieldValue.Kind() != reflect.Slice && fieldValue.Kind() != reflect.Array {
			return false
		}

		for i := 0; i < fieldValue.Len(); i++ {
			elemValue := fieldValue.Index(i)

			for _, val := range values {
				targetValue := reflect.ValueOf(val)
				equal, err := compareValues(elemValue, targetValue)
				if err == nil && equal {
					return true
				}
			}
		}

		return false
	})
}

// ArrayContainsAll checks if an array contains all of the provided values
func ArrayContainsAll[T any](fieldName string, values []interface{}) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		if fieldValue.Kind() != reflect.Slice && fieldValue.Kind() != reflect.Array {
			return false
		}

		for _, val := range values {
			targetValue := reflect.ValueOf(val)
			found := false

			for i := 0; i < fieldValue.Len(); i++ {
				elemValue := fieldValue.Index(i)
				equal, err := compareValues(elemValue, targetValue)
				if err == nil && equal {
					found = true
					break
				}
			}

			if !found {
				return false
			}
		}

		return true
	})
}

// Between returns a filter that checks if a field value is within a range (inclusive).
// Works with numeric types, strings, and any comparable type.
//
// Example:
//
//	filter.Between[User]("Age", 18, 65)  // users where 18 <= Age <= 65
func Between[T any](fieldName string, min, max interface{}) Filter[T] {
	return And[T](
		Gte[T](fieldName, min),
		Lte[T](fieldName, max),
	)
}

// DateBefore returns a filter that checks if a date field is before the specified date
func DateBefore[T any](fieldName string, date time.Time) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		// Extract time.Time from field
		var fieldTime time.Time

		switch fieldValue.Kind() {
		case reflect.Struct:
			if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
				fieldTime = fieldValue.Interface().(time.Time)
			} else {
				return false
			}
		case reflect.String:
			// Try to parse string as date (try multiple formats)
			formats := []string{
				time.RFC3339,
				"2006-01-02T15:04:05",
				"2006-01-02",
				"01/02/2006",
				"02/01/2006",
			}

			for _, format := range formats {
				parsedTime, err := time.Parse(format, fieldValue.String())
				if err == nil {
					fieldTime = parsedTime
					break
				}
			}

			// If we couldn't parse the date, return false
			if fieldTime.IsZero() {
				return false
			}
		default:
			return false
		}

		return fieldTime.Before(date)
	})
}

// DateAfter returns a filter that checks if a date field is after the specified date
func DateAfter[T any](fieldName string, date time.Time) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil {
			return false
		}

		// Extract time.Time from field
		var fieldTime time.Time

		switch fieldValue.Kind() {
		case reflect.Struct:
			if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
				fieldTime = fieldValue.Interface().(time.Time)
			} else {
				return false
			}
		case reflect.String:
			// Try to parse string as date (try multiple formats)
			formats := []string{
				time.RFC3339,
				"2006-01-02T15:04:05",
				"2006-01-02",
				"01/02/2006",
				"02/01/2006",
			}

			for _, format := range formats {
				parsedTime, err := time.Parse(format, fieldValue.String())
				if err == nil {
					fieldTime = parsedTime
					break
				}
			}

			// If we couldn't parse the date, return false
			if fieldTime.IsZero() {
				return false
			}
		default:
			return false
		}

		return fieldTime.After(date)
	})
}

// DateBetween returns a filter that checks if a date field is between two dates (inclusive)
func DateBetween[T any](fieldName string, start, end time.Time) Filter[T] {
	return And[T](
		DateAfter[T](fieldName, start.Add(-1*time.Second)), // Make it inclusive of start time
		DateBefore[T](fieldName, end.Add(1*time.Second)),   // Make it inclusive of end time
	)
}

// Sort returns a sorted copy of the slice based on a field value.
// The original slice is not modified. Set ascending to true for A-Z/0-9 order.
//
// Example:
//
//	sorted := filter.Sort(users, "Name", true)   // sort by Name ascending
//	sorted := filter.Sort(users, "Age", false)   // sort by Age descending
func Sort[T any](items []T, fieldName string, ascending bool) []T {
	result := make([]T, len(items))
	copy(result, items)

	sort.Slice(result, func(i, j int) bool {
		fieldValueI, err := getFieldValue(result[i], fieldName)
		if err != nil {
			return false
		}

		fieldValueJ, err := getFieldValue(result[j], fieldName)
		if err != nil {
			return false
		}

		less, err := compareValuesLess(fieldValueI, fieldValueJ)
		if err != nil {
			return false
		}

		if ascending {
			return less
		}
		return !less
	})

	return result
}

// Custom creates a filter from a user-provided function.
// Use this when built-in operators don't cover your use case.
//
// Example:
//
//	filter.Custom[User](func(u User) bool {
//	    return u.Age > 18 && strings.Contains(u.Email, "@company.com")
//	})
func Custom[T any](fn func(T) bool) Filter[T] {
	return FilterFunc[T](fn)
}

// RegexMatch returns a filter that checks if a string field matches a regular expression.
// If the pattern is invalid, the filter will never match (no panic).
//
// Example:
//
//	filter.RegexMatch[User]("Email", `^[a-z]+@gmail\.com$`)  // Gmail users
func RegexMatch[T any](fieldName, pattern string) Filter[T] {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		// If the pattern is invalid, the filter will never match
		return FilterFunc[T](func(T) bool { return false })
	}

	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, fieldName)
		if err != nil || fieldValue.Kind() != reflect.String {
			return false
		}

		return regex.MatchString(fieldValue.String())
	})
}

// NestedArrayAny filters based on a condition in any element of a nested array
func NestedArrayAny[T any](arrayField string, conditionFn func(elem reflect.Value) bool) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, arrayField)
		if err != nil || fieldValue.Kind() != reflect.Slice {
			return false
		}

		for i := 0; i < fieldValue.Len(); i++ {
			if conditionFn(fieldValue.Index(i)) {
				return true
			}
		}

		return false
	})
}

// NestedArrayAll filters based on a condition in all elements of a nested array
func NestedArrayAll[T any](arrayField string, conditionFn func(elem reflect.Value) bool) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		fieldValue, err := getFieldValue(item, arrayField)
		if err != nil || fieldValue.Kind() != reflect.Slice {
			return false
		}

		if fieldValue.Len() == 0 {
			return false
		}

		for i := 0; i < fieldValue.Len(); i++ {
			if !conditionFn(fieldValue.Index(i)) {
				return false
			}
		}

		return true
	})
}
