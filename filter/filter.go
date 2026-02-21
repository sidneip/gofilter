// Package filter provides a composable, type-safe filtering engine for Go slices.
// It enables building complex filter expressions using a fluent API without
// requiring a database.
//
// The package supports:
//   - Comparison operators (Eq, Ne, Gt, Lt, Gte, Lte)
//   - String matching (Contains, StringMatch with regex support)
//   - Collection operations (In, Between, ArrayContains)
//   - Logical composition (And, Or, Not)
//   - Geospatial queries (WithinRadius, WithinBoundingBox)
//   - Map field queries (HasKey, KeyValueEquals)
//   - Custom filter functions
//
// Example:
//
//	users := []User{{Name: "Ana", Age: 25}, {Name: "Bob", Age: 30}}
//	adults := filter.Apply(users, filter.Gte[User]("Age", 18))
package filter

// Filter is the core interface that all filters must implement.
// It defines a single method Apply that returns true if an item
// passes the filter criteria.
type Filter[T any] interface {
	// Apply checks if the given item passes the filter criteria
	Apply(item T) bool
}

// FilterFunc is a function adapter that implements the Filter interface.
// It allows using simple functions as filters.
//
// Example:
//
//	customFilter := filter.FilterFunc[User](func(u User) bool {
//	    return u.Age > 18 && strings.HasPrefix(u.Name, "A")
//	})
type FilterFunc[T any] func(T) bool

// Apply implements the Filter interface for FilterFunc.
func (f FilterFunc[T]) Apply(item T) bool {
	return f(item)
}

// Apply filters a slice using the provided filter and returns a new slice
// containing only the items that pass the filter criteria.
// The original slice is not modified.
//
// Example:
//
//	adults := filter.Apply(users, filter.Gte[User]("Age", 18))
func Apply[T any](items []T, filter Filter[T]) []T {
	result := make([]T, 0)

	for _, item := range items {
		if filter.Apply(item) {
			result = append(result, item)
		}
	}

	return result
}
