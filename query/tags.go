package query

import (
	"reflect"
	"strings"
	"unicode"
)

type fieldInfo struct {
	structField string
	column      string
	filterable  bool
	sortable    bool
	fieldType   reflect.Type
}

type fieldRegistry struct {
	fields   []fieldInfo
	byColumn map[string]fieldInfo
}

func parseStructTags[T any]() (*fieldRegistry, error) {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	reg := &fieldRegistry{
		byColumn: make(map[string]fieldInfo),
	}

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		tag := sf.Tag.Get("gofilter")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		info := fieldInfo{
			structField: sf.Name,
			column:      toSnakeCase(sf.Name),
			fieldType:   sf.Type,
		}

		for _, part := range parts {
			part = strings.TrimSpace(part)
			switch {
			case part == "filterable":
				info.filterable = true
			case part == "sortable":
				info.sortable = true
			case strings.HasPrefix(part, "column="):
				info.column = strings.TrimPrefix(part, "column=")
			}
		}

		if !info.filterable {
			continue
		}

		reg.fields = append(reg.fields, info)
		reg.byColumn[info.column] = info
	}

	return reg, nil
}

func toSnakeCase(s string) string {
	var result []rune
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := runes[i-1]
				if unicode.IsLower(prev) {
					result = append(result, '_')
				} else if unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					result = append(result, '_')
				}
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
