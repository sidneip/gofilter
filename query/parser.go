package query

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var operators = []string{"between", "contains", "gte", "gt", "lte", "lt", "ne", "in"}

var reservedParams = map[string]bool{
	"sort":  true,
	"page":  true,
	"limit": true,
}

type parsedFilter struct {
	field    string
	operator string
	value    interface{}
}

type parsedQuery struct {
	filters   []parsedFilter
	sortField string
	sortAsc   bool
	page      int
	limit     int
}

func splitParamOperator(param string) (column, operator string) {
	for _, op := range operators {
		suffix := "_" + op
		if strings.HasSuffix(param, suffix) {
			return strings.TrimSuffix(param, suffix), op
		}
	}
	return param, "eq"
}

func parseParams[T any](params url.Values, opts options) (*parsedQuery, error) {
	registry, err := parseStructTags[T]()
	if err != nil {
		return nil, err
	}

	result := &parsedQuery{
		page:  1,
		limit: opts.defaultLimit,
	}

	for param, values := range params {
		if len(values) == 0 {
			continue
		}
		raw := values[0]

		if reservedParams[param] {
			switch param {
			case "sort":
				sortField, asc, err := parseSortParam(raw, registry)
				if err != nil {
					return nil, err
				}
				result.sortField = sortField
				result.sortAsc = asc
			case "page":
				p, err := strconv.Atoi(raw)
				if err != nil || p < 1 {
					return nil, &ErrInvalidValue{Field: "page", Value: raw, ExpectedType: "positive integer"}
				}
				result.page = p
			case "limit":
				l, err := strconv.Atoi(raw)
				if err != nil || l < 1 {
					return nil, &ErrInvalidValue{Field: "limit", Value: raw, ExpectedType: "positive integer"}
				}
				if opts.maxLimit > 0 && l > opts.maxLimit {
					return nil, &ErrLimitExceeded{Requested: l, Max: opts.maxLimit}
				}
				result.limit = l
			}
			continue
		}

		col, op := splitParamOperator(param)

		info, ok := registry.byColumn[col]
		if !ok {
			return nil, &ErrFieldNotFilterable{Field: col}
		}

		coerced, err := coerceFilterValue(raw, op, info)
		if err != nil {
			return nil, &ErrInvalidValue{Field: info.structField, Value: raw, ExpectedType: info.fieldType.String()}
		}

		result.filters = append(result.filters, parsedFilter{
			field:    info.structField,
			operator: op,
			value:    coerced,
		})
	}

	return result, nil
}

func parseSortParam(raw string, registry *fieldRegistry) (string, bool, error) {
	asc := true
	field := raw
	if strings.HasPrefix(raw, "-") {
		asc = false
		field = raw[1:]
	}

	info, ok := registry.byColumn[field]
	if !ok {
		return "", false, &ErrFieldNotSortable{Field: field}
	}
	if !info.sortable {
		return "", false, &ErrFieldNotSortable{Field: field}
	}

	return info.structField, asc, nil
}

func coerceFilterValue(raw, op string, info fieldInfo) (interface{}, error) {
	switch op {
	case "in":
		parts := strings.Split(raw, ",")
		vals := make([]interface{}, 0, len(parts))
		for _, p := range parts {
			v, err := coerceValue(strings.TrimSpace(p), info.fieldType)
			if err != nil {
				return nil, err
			}
			vals = append(vals, v)
		}
		return vals, nil
	case "between":
		parts := strings.SplitN(raw, ",", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("between requires exactly 2 comma-separated values")
		}
		min, err := coerceValue(strings.TrimSpace(parts[0]), info.fieldType)
		if err != nil {
			return nil, err
		}
		max, err := coerceValue(strings.TrimSpace(parts[1]), info.fieldType)
		if err != nil {
			return nil, err
		}
		return [2]interface{}{min, max}, nil
	default:
		return coerceValue(raw, info.fieldType)
	}
}
