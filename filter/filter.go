package filter

// Filter is the core interface that all filters must implement
type Filter[T any] interface {
	// Apply checks if the given item passes the filter criteria
	Apply(item T) bool
}

// FilterFunc is a function type that implements the Filter interface
type FilterFunc[T any] func(T) bool

// Apply implements the Filter interface for FilterFunc
func (f FilterFunc[T]) Apply(item T) bool {
	return f(item)
}

// Apply applies a filter to a slice of items and returns a new slice with only the items that pass the filter
func Apply[T any](items []T, filter Filter[T]) []T {
	result := make([]T, 0)

	for _, item := range items {
		if filter.Apply(item) {
			result = append(result, item)
		}
	}

	return result
}
