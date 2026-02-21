// Package query provides HTTP query parameter parsing for filtering, sorting,
// and paginating slices of structs. It bridges URL query strings to the filter
// package, enabling REST APIs to expose filterable endpoints with minimal code.
//
// Example usage:
//
//	type User struct {
//	    Name string `gofilter:"filterable,sortable"`
//	    Age  int    `gofilter:"filterable,sortable"`
//	}
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    result, err := query.ApplyPaginated(users, r.URL.Query())
//	    // result contains filtered, sorted, paginated items
//	}
package query

import (
	"net/url"

	"github.com/sidneip/gofilter/filter"
)

// PageResult represents a paginated response containing filtered items
// along with pagination metadata. It is designed to be JSON-serialized
// directly in HTTP responses.
type PageResult[T any] struct {
	// Items contains the filtered and paginated slice of results
	Items []T `json:"items"`
	// Total is the count of all items matching the filter (before pagination)
	Total int `json:"total"`
	// Page is the current page number (1-based)
	Page int `json:"page"`
	// Limit is the maximum number of items per page
	Limit int `json:"limit"`
	// HasNext indicates whether there are more pages available
	HasNext bool `json:"has_next"`
}

type options struct {
	defaultLimit   int
	maxLimit       int
	defaultSort    string
	defaultSortAsc bool
}

// Option is a functional option for configuring query behavior.
// Use WithMaxLimit, WithDefaultSort, and WithDefaultLimit to create options.
type Option func(*options)

func defaultOptions() options {
	return options{
		defaultLimit: 20,
	}
}

// WithMaxLimit sets the maximum allowed limit for pagination.
// Requests exceeding this limit will return ErrLimitExceeded.
//
// Example:
//
//	query.Apply(items, params, query.WithMaxLimit(100))
func WithMaxLimit(max int) Option {
	return func(o *options) {
		o.maxLimit = max
	}
}

// WithDefaultSort sets the default sort field and direction when no sort
// parameter is provided in the query string. If ascending is true, items
// are sorted in ascending order; otherwise, descending.
//
// Example:
//
//	query.Apply(items, params, query.WithDefaultSort("Name", true))
func WithDefaultSort(field string, ascending bool) Option {
	return func(o *options) {
		o.defaultSort = field
		o.defaultSortAsc = ascending
	}
}

// WithDefaultLimit sets the default number of items per page when no limit
// parameter is provided in the query string. The default is 20.
//
// Example:
//
//	query.Apply(items, params, query.WithDefaultLimit(50))
func WithDefaultLimit(limit int) Option {
	return func(o *options) {
		o.defaultLimit = limit
	}
}

// Apply filters and sorts a slice based on URL query parameters.
// It parses the query string for filter operators (eq, gt, lt, contains, etc.),
// applies them to the slice, and returns the filtered result.
//
// Supported query syntax:
//   - field=value        → equals
//   - field_gt=value     → greater than
//   - field_gte=value    → greater than or equal
//   - field_lt=value     → less than
//   - field_lte=value    → less than or equal
//   - field_ne=value     → not equal
//   - field_contains=val → substring match
//   - field_in=a,b,c     → value in list
//   - field_between=a,b  → value between a and b
//   - sort=field         → sort ascending
//   - sort=-field        → sort descending
//
// Returns an error if the query contains invalid parameters or values.
func Apply[T any](items []T, params url.Values, opts ...Option) ([]T, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}

	parsed, err := parseParams[T](params, o)
	if err != nil {
		return nil, err
	}

	result := items
	if len(parsed.filters) > 0 {
		filters := make([]filter.Filter[T], 0, len(parsed.filters))
		for _, pf := range parsed.filters {
			f := buildFilter[T](pf)
			filters = append(filters, f)
		}
		result = filter.Apply(result, filter.And(filters...))
	}

	sortField := parsed.sortField
	sortAsc := parsed.sortAsc
	if sortField == "" && o.defaultSort != "" {
		sortField = o.defaultSort
		sortAsc = o.defaultSortAsc
	}
	if sortField != "" {
		result = filter.Sort(result, sortField, sortAsc)
	}

	return result, nil
}

// ApplyPaginated filters, sorts, and paginates a slice based on URL query parameters.
// It extends Apply with pagination support, returning a PageResult containing
// the items for the requested page along with pagination metadata.
//
// Additional query parameters for pagination:
//   - page=N  → page number (1-based, default: 1)
//   - limit=N → items per page (default: 20)
//
// Example:
//
//	// GET /users?city=SP&sort=-age&page=2&limit=10
//	result, err := query.ApplyPaginated(users, r.URL.Query(),
//	    query.WithMaxLimit(100),
//	)
//	// result.Items contains up to 10 users from São Paulo, sorted by age desc
//	// result.Total contains the total count matching the filter
//	// result.HasNext indicates if there are more pages
func ApplyPaginated[T any](items []T, params url.Values, opts ...Option) (*PageResult[T], error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}

	parsed, err := parseParams[T](params, o)
	if err != nil {
		return nil, err
	}

	result := items
	if len(parsed.filters) > 0 {
		filters := make([]filter.Filter[T], 0, len(parsed.filters))
		for _, pf := range parsed.filters {
			f := buildFilter[T](pf)
			filters = append(filters, f)
		}
		result = filter.Apply(result, filter.And(filters...))
	}

	sortField := parsed.sortField
	sortAsc := parsed.sortAsc
	if sortField == "" && o.defaultSort != "" {
		sortField = o.defaultSort
		sortAsc = o.defaultSortAsc
	}
	if sortField != "" {
		result = filter.Sort(result, sortField, sortAsc)
	}

	total := len(result)
	page := parsed.page
	limit := parsed.limit
	if limit <= 0 {
		limit = o.defaultLimit
	}

	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	return &PageResult[T]{
		Items:   result[start:end],
		Total:   total,
		Page:    page,
		Limit:   limit,
		HasNext: end < total,
	}, nil
}

func buildFilter[T any](pf parsedFilter) filter.Filter[T] {
	switch pf.operator {
	case "eq":
		return filter.Eq[T](pf.field, pf.value)
	case "ne":
		return filter.Ne[T](pf.field, pf.value)
	case "gt":
		return filter.Gt[T](pf.field, pf.value)
	case "gte":
		return filter.Gte[T](pf.field, pf.value)
	case "lt":
		return filter.Lt[T](pf.field, pf.value)
	case "lte":
		return filter.Lte[T](pf.field, pf.value)
	case "contains":
		return filter.Contains[T](pf.field, pf.value)
	case "in":
		vals, ok := pf.value.([]interface{})
		if !ok {
			return filter.FilterFunc[T](func(T) bool { return false })
		}
		return filter.In[T](pf.field, vals)
	case "between":
		vals, ok := pf.value.([2]interface{})
		if !ok {
			return filter.FilterFunc[T](func(T) bool { return false })
		}
		return filter.Between[T](pf.field, vals[0], vals[1])
	default:
		return filter.FilterFunc[T](func(T) bool { return false })
	}
}
