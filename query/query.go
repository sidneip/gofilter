package query

import (
	"net/url"

	"github.com/sidneip/gofilter/filter"
)

type PageResult[T any] struct {
	Items   []T  `json:"items"`
	Total   int  `json:"total"`
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	HasNext bool `json:"has_next"`
}

type options struct {
	defaultLimit   int
	maxLimit       int
	defaultSort    string
	defaultSortAsc bool
}

type Option func(*options)

func defaultOptions() options {
	return options{
		defaultLimit: 20,
	}
}

func WithMaxLimit(max int) Option {
	return func(o *options) {
		o.maxLimit = max
	}
}

func WithDefaultSort(field string, ascending bool) Option {
	return func(o *options) {
		o.defaultSort = field
		o.defaultSortAsc = ascending
	}
}

func WithDefaultLimit(limit int) Option {
	return func(o *options) {
		o.defaultLimit = limit
	}
}

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
