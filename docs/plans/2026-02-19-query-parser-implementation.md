# Query Parser Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a `query/` package that parses HTTP query strings into gofilter in-memory filters, making gofilter the only Go lib that does end-to-end `?field_op=value` → filtered/paginated results.

**Architecture:** New `query/` package imports `filter/`. Struct tags (`gofilter:"filterable"`) control field access. Django-style suffix syntax (`age_gt=25`). Type coercion via reflection. All filters combined with AND.

**Tech Stack:** Go 1.22+, generics, reflect, strconv, net/url. Zero external deps.

---

### Task 1: Error types

**Files:**
- Create: `query/errors.go`
- Test: `query/errors_test.go`

**Step 1: Write the failing test**

```go
// query/errors_test.go
package query

import (
	"errors"
	"testing"
)

func TestErrFieldNotFilterable(t *testing.T) {
	err := &ErrFieldNotFilterable{Field: "SSN"}
	if err.Error() != `field "SSN" is not filterable` {
		t.Errorf("unexpected error message: %s", err.Error())
	}
	var target *ErrFieldNotFilterable
	if !errors.As(err, &target) {
		t.Error("should be assertable with errors.As")
	}
}

func TestErrFieldNotSortable(t *testing.T) {
	err := &ErrFieldNotSortable{Field: "Email"}
	if err.Error() != `field "Email" is not sortable` {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestErrInvalidValue(t *testing.T) {
	err := &ErrInvalidValue{Field: "Age", Value: "abc", ExpectedType: "int"}
	if err.Error() != `invalid value "abc" for field "Age": expected int` {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestErrLimitExceeded(t *testing.T) {
	err := &ErrLimitExceeded{Requested: 500, Max: 100}
	if err.Error() != "requested limit 500 exceeds maximum 100" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./query/ -run TestErr -v`
Expected: FAIL — package doesn't exist yet

**Step 3: Write minimal implementation**

```go
// query/errors.go
package query

import "fmt"

type ErrFieldNotFilterable struct{ Field string }

func (e *ErrFieldNotFilterable) Error() string {
	return fmt.Sprintf("field %q is not filterable", e.Field)
}

type ErrFieldNotSortable struct{ Field string }

func (e *ErrFieldNotSortable) Error() string {
	return fmt.Sprintf("field %q is not sortable", e.Field)
}

type ErrInvalidValue struct{ Field, Value, ExpectedType string }

func (e *ErrInvalidValue) Error() string {
	return fmt.Sprintf("invalid value %q for field %q: expected %s", e.Value, e.Field, e.ExpectedType)
}

type ErrLimitExceeded struct{ Requested, Max int }

func (e *ErrLimitExceeded) Error() string {
	return fmt.Sprintf("requested limit %d exceeds maximum %d", e.Requested, e.Max)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./query/ -run TestErr -v`
Expected: PASS

**Step 5: Commit**

```bash
git add query/errors.go query/errors_test.go
git commit -m "feat(query): add typed error types for query parsing"
```

---

### Task 2: Struct tag parsing

**Files:**
- Create: `query/tags.go`
- Test: `query/tags_test.go`

**Step 1: Write the failing test**

```go
// query/tags_test.go
package query

import (
	"testing"
	"time"
)

type TagTestUser struct {
	Name    string `gofilter:"filterable,sortable,column=name"`
	Age     int    `gofilter:"filterable"`
	City    string `gofilter:"filterable,sortable"`
	Email   string
	Score   float64   `gofilter:"filterable"`
	Active  bool      `gofilter:"filterable"`
	Created time.Time `gofilter:"filterable,sortable"`
}

func TestParseStructTags(t *testing.T) {
	registry, err := parseStructTags[TagTestUser]()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have 6 filterable fields (Name, Age, City, Score, Active, Created)
	if len(registry.fields) != 6 {
		t.Errorf("expected 6 fields, got %d", len(registry.fields))
	}

	// Name should use custom column name "name"
	nameField, ok := registry.byColumn["name"]
	if !ok {
		t.Fatal("field with column 'name' not found")
	}
	if nameField.structField != "Name" {
		t.Errorf("expected structField 'Name', got %q", nameField.structField)
	}
	if !nameField.filterable {
		t.Error("Name should be filterable")
	}
	if !nameField.sortable {
		t.Error("Name should be sortable")
	}

	// Age should use default column name "age"
	ageField, ok := registry.byColumn["age"]
	if !ok {
		t.Fatal("field with column 'age' not found")
	}
	if ageField.structField != "Age" {
		t.Errorf("expected structField 'Age', got %q", ageField.structField)
	}
	if !ageField.filterable {
		t.Error("Age should be filterable")
	}
	if ageField.sortable {
		t.Error("Age should NOT be sortable")
	}

	// Email should not be in registry
	if _, ok := registry.byColumn["email"]; ok {
		t.Error("Email should not be in registry")
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"Name", "name"},
		{"Age", "age"},
		{"FirstName", "first_name"},
		{"CreatedAt", "created_at"},
		{"HTMLParser", "html_parser"},
		{"ID", "id"},
		{"UserID", "user_id"},
	}
	for _, tt := range tests {
		got := toSnakeCase(tt.input)
		if got != tt.want {
			t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./query/ -run "TestParseStructTags|TestToSnakeCase" -v`
Expected: FAIL — functions don't exist

**Step 3: Write minimal implementation**

```go
// query/tags.go
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
```

**Step 4: Run test to verify it passes**

Run: `go test ./query/ -run "TestParseStructTags|TestToSnakeCase" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add query/tags.go query/tags_test.go
git commit -m "feat(query): add struct tag parsing with field registry"
```

---

### Task 3: Type coercion

**Files:**
- Create: `query/coerce.go`
- Test: `query/coerce_test.go`

**Step 1: Write the failing test**

```go
// query/coerce_test.go
package query

import (
	"reflect"
	"testing"
	"time"
)

func TestCoerceString(t *testing.T) {
	val, err := coerceValue("hello", reflect.TypeOf(""))
	if err != nil {
		t.Fatal(err)
	}
	if val != "hello" {
		t.Errorf("expected 'hello', got %v", val)
	}
}

func TestCoerceInt(t *testing.T) {
	val, err := coerceValue("42", reflect.TypeOf(0))
	if err != nil {
		t.Fatal(err)
	}
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestCoerceInt64(t *testing.T) {
	val, err := coerceValue("100", reflect.TypeOf(int64(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != int64(100) {
		t.Errorf("expected int64(100), got %v (%T)", val, val)
	}
}

func TestCoerceUint(t *testing.T) {
	val, err := coerceValue("10", reflect.TypeOf(uint(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != uint(10) {
		t.Errorf("expected uint(10), got %v (%T)", val, val)
	}
}

func TestCoerceFloat64(t *testing.T) {
	val, err := coerceValue("3.14", reflect.TypeOf(float64(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != 3.14 {
		t.Errorf("expected 3.14, got %v", val)
	}
}

func TestCoerceBool(t *testing.T) {
	val, err := coerceValue("true", reflect.TypeOf(false))
	if err != nil {
		t.Fatal(err)
	}
	if val != true {
		t.Errorf("expected true, got %v", val)
	}
}

func TestCoerceTimeRFC3339(t *testing.T) {
	val, err := coerceValue("2024-01-15T10:30:00Z", reflect.TypeOf(time.Time{}))
	if err != nil {
		t.Fatal(err)
	}
	tm, ok := val.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time, got %T", val)
	}
	if tm.Year() != 2024 || tm.Month() != 1 || tm.Day() != 15 {
		t.Errorf("unexpected time: %v", tm)
	}
}

func TestCoerceTimeDateOnly(t *testing.T) {
	val, err := coerceValue("2024-01-15", reflect.TypeOf(time.Time{}))
	if err != nil {
		t.Fatal(err)
	}
	tm, ok := val.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time, got %T", val)
	}
	if tm.Year() != 2024 || tm.Month() != 1 || tm.Day() != 15 {
		t.Errorf("unexpected time: %v", tm)
	}
}

func TestCoerceInvalidInt(t *testing.T) {
	_, err := coerceValue("abc", reflect.TypeOf(0))
	if err == nil {
		t.Error("expected error for invalid int")
	}
}

func TestCoerceFloat32(t *testing.T) {
	val, err := coerceValue("2.5", reflect.TypeOf(float32(0)))
	if err != nil {
		t.Fatal(err)
	}
	if val != float32(2.5) {
		t.Errorf("expected float32(2.5), got %v (%T)", val, val)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./query/ -run TestCoerce -v`
Expected: FAIL

**Step 3: Write minimal implementation**

```go
// query/coerce.go
package query

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

func coerceValue(raw string, targetType reflect.Type) (interface{}, error) {
	if targetType == timeType {
		return parseTime(raw)
	}

	switch targetType.Kind() {
	case reflect.String:
		return raw, nil
	case reflect.Int:
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as int: %w", raw, err)
		}
		return int(v), nil
	case reflect.Int8:
		v, err := strconv.ParseInt(raw, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as int8: %w", raw, err)
		}
		return int8(v), nil
	case reflect.Int16:
		v, err := strconv.ParseInt(raw, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as int16: %w", raw, err)
		}
		return int16(v), nil
	case reflect.Int32:
		v, err := strconv.ParseInt(raw, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as int32: %w", raw, err)
		}
		return int32(v), nil
	case reflect.Int64:
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as int64: %w", raw, err)
		}
		return v, nil
	case reflect.Uint:
		v, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as uint: %w", raw, err)
		}
		return uint(v), nil
	case reflect.Uint8:
		v, err := strconv.ParseUint(raw, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as uint8: %w", raw, err)
		}
		return uint8(v), nil
	case reflect.Uint16:
		v, err := strconv.ParseUint(raw, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as uint16: %w", raw, err)
		}
		return uint16(v), nil
	case reflect.Uint32:
		v, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as uint32: %w", raw, err)
		}
		return uint32(v), nil
	case reflect.Uint64:
		v, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as uint64: %w", raw, err)
		}
		return v, nil
	case reflect.Float32:
		v, err := strconv.ParseFloat(raw, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as float32: %w", raw, err)
		}
		return float32(v), nil
	case reflect.Float64:
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as float64: %w", raw, err)
		}
		return v, nil
	case reflect.Bool:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %q as bool: %w", raw, err)
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", targetType)
	}
}

func parseTime(raw string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02",
	}
	for _, f := range formats {
		t, err := time.Parse(f, raw)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse %q as time", raw)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./query/ -run TestCoerce -v`
Expected: PASS

**Step 5: Commit**

```bash
git add query/coerce.go query/coerce_test.go
git commit -m "feat(query): add type coercion from string to Go types"
```

---

### Task 4: Query param parser (operator detection + field resolution)

**Files:**
- Create: `query/parser.go`
- Test: `query/parser_test.go`

**Step 1: Write the failing test**

```go
// query/parser_test.go
package query

import (
	"net/url"
	"testing"
)

type ParserTestUser struct {
	Name string `gofilter:"filterable,sortable,column=name"`
	Age  int    `gofilter:"filterable"`
	City string `gofilter:"filterable,sortable"`
}

func TestParseSimpleEq(t *testing.T) {
	params := url.Values{"name": {"Ana"}}
	parsed, err := parseParams[ParserTestUser](params, defaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(parsed.filters))
	}
	if parsed.filters[0].field != "Name" {
		t.Errorf("expected field 'Name', got %q", parsed.filters[0].field)
	}
	if parsed.filters[0].operator != "eq" {
		t.Errorf("expected operator 'eq', got %q", parsed.filters[0].operator)
	}
}

func TestParseOperatorSuffix(t *testing.T) {
	tests := []struct {
		param    string
		wantOp   string
		wantCol  string
	}{
		{"age_gt", "gt", "age"},
		{"age_gte", "gte", "age"},
		{"age_lt", "lt", "age"},
		{"age_lte", "lte", "age"},
		{"city_ne", "ne", "city"},
		{"name_contains", "contains", "name"},
		{"city_in", "in", "city"},
		{"age_between", "between", "age"},
		{"name", "eq", "name"},
	}
	for _, tt := range tests {
		col, op := splitParamOperator(tt.param)
		if op != tt.wantOp {
			t.Errorf("splitParamOperator(%q): op = %q, want %q", tt.param, op, tt.wantOp)
		}
		if col != tt.wantCol {
			t.Errorf("splitParamOperator(%q): col = %q, want %q", tt.param, col, tt.wantCol)
		}
	}
}

func TestParseReservedParams(t *testing.T) {
	params := url.Values{
		"name":  {"Ana"},
		"sort":  {"-age"},
		"page":  {"2"},
		"limit": {"10"},
	}
	parsed, err := parseParams[ParserTestUser](params, defaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.filters) != 1 {
		t.Fatalf("expected 1 filter (sort/page/limit are reserved), got %d", len(parsed.filters))
	}
	if parsed.sortField != "Age" {
		t.Errorf("expected sort field 'Age', got %q", parsed.sortField)
	}
	if parsed.sortAsc {
		t.Error("expected descending sort")
	}
	if parsed.page != 2 {
		t.Errorf("expected page 2, got %d", parsed.page)
	}
	if parsed.limit != 10 {
		t.Errorf("expected limit 10, got %d", parsed.limit)
	}
}

func TestParseFieldNotFilterable(t *testing.T) {
	type Restricted struct {
		Name string `gofilter:"filterable"`
		SSN  string
	}
	params := url.Values{"ssn": {"123"}}
	_, err := parseParams[Restricted](params, defaultOptions())
	if err == nil {
		t.Fatal("expected error for non-filterable field")
	}
	if _, ok := err.(*ErrFieldNotFilterable); !ok {
		t.Errorf("expected ErrFieldNotFilterable, got %T: %v", err, err)
	}
}

func TestParseFieldNotSortable(t *testing.T) {
	type Limited struct {
		Name string `gofilter:"filterable"`
	}
	params := url.Values{"sort": {"name"}}
	_, err := parseParams[Limited](params, defaultOptions())
	if err == nil {
		t.Fatal("expected error for non-sortable field")
	}
	if _, ok := err.(*ErrFieldNotSortable); !ok {
		t.Errorf("expected ErrFieldNotSortable, got %T: %v", err, err)
	}
}

func TestParseInvalidValue(t *testing.T) {
	params := url.Values{"age": {"abc"}}
	_, err := parseParams[ParserTestUser](params, defaultOptions())
	if err == nil {
		t.Fatal("expected error for invalid int value")
	}
	if _, ok := err.(*ErrInvalidValue); !ok {
		t.Errorf("expected ErrInvalidValue, got %T: %v", err, err)
	}
}

func TestParseLimitExceeded(t *testing.T) {
	params := url.Values{"limit": {"500"}}
	opts := defaultOptions()
	opts.maxLimit = 100
	_, err := parseParams[ParserTestUser](params, opts)
	if err == nil {
		t.Fatal("expected error for limit exceeded")
	}
	if _, ok := err.(*ErrLimitExceeded); !ok {
		t.Errorf("expected ErrLimitExceeded, got %T: %v", err, err)
	}
}

func TestParseInOperator(t *testing.T) {
	params := url.Values{"city_in": {"SP,RJ,MG"}}
	parsed, err := parseParams[ParserTestUser](params, defaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(parsed.filters))
	}
	f := parsed.filters[0]
	if f.operator != "in" {
		t.Errorf("expected 'in' operator, got %q", f.operator)
	}
	vals, ok := f.value.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{} value, got %T", f.value)
	}
	if len(vals) != 3 {
		t.Errorf("expected 3 values, got %d", len(vals))
	}
}

func TestParseBetweenOperator(t *testing.T) {
	params := url.Values{"age_between": {"18,30"}}
	parsed, err := parseParams[ParserTestUser](params, defaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(parsed.filters))
	}
	f := parsed.filters[0]
	if f.operator != "between" {
		t.Errorf("expected 'between' operator, got %q", f.operator)
	}
	vals, ok := f.value.([2]interface{})
	if !ok {
		t.Fatalf("expected [2]interface{} value, got %T", f.value)
	}
	if vals[0] != 18 || vals[1] != 30 {
		t.Errorf("expected [18, 30], got %v", vals)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./query/ -run "TestParse|TestParseOperator" -v`
Expected: FAIL

**Step 3: Write minimal implementation**

```go
// query/parser.go
package query

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var operators = []string{"gte", "gt", "lte", "lt", "ne", "contains", "between", "in"}

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
```

**Step 4: Run test to verify it passes**

Run: `go test ./query/ -run "TestParse" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add query/parser.go query/parser_test.go
git commit -m "feat(query): add query param parser with operator detection"
```

---

### Task 5: Apply and ApplyPaginated (main API)

**Files:**
- Create: `query/query.go`
- Test: `query/query_test.go`

**Step 1: Write the failing test**

```go
// query/query_test.go
package query

import (
	"net/url"
	"testing"
)

type User struct {
	Name string `gofilter:"filterable,sortable,column=name"`
	Age  int    `gofilter:"filterable,sortable"`
	City string `gofilter:"filterable,sortable"`
}

func testUsers() []User {
	return []User{
		{Name: "Ana", Age: 20, City: "SP"},
		{Name: "Bruno", Age: 17, City: "RJ"},
		{Name: "Carla", Age: 25, City: "SP"},
		{Name: "Daniel", Age: 30, City: "MG"},
		{Name: "Elena", Age: 22, City: "RJ"},
	}
}

func TestApplyEq(t *testing.T) {
	params := url.Values{"city": {"SP"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 users from SP, got %d", len(result))
	}
}

func TestApplyGt(t *testing.T) {
	params := url.Values{"age_gt": {"20"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	// Carla(25), Daniel(30), Elena(22)
	if len(result) != 3 {
		t.Errorf("expected 3 users with age > 20, got %d", len(result))
	}
}

func TestApplyMultipleFilters(t *testing.T) {
	params := url.Values{"city": {"SP"}, "age_gt": {"20"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	// Only Carla (SP, age 25)
	if len(result) != 1 {
		t.Errorf("expected 1 user, got %d", len(result))
	}
	if len(result) > 0 && result[0].Name != "Carla" {
		t.Errorf("expected Carla, got %s", result[0].Name)
	}
}

func TestApplyContains(t *testing.T) {
	params := url.Values{"name_contains": {"an"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	// Daniel, Elena (contain "an")
	if len(result) != 2 {
		t.Errorf("expected 2 users with 'an' in name, got %d", len(result))
	}
}

func TestApplyIn(t *testing.T) {
	params := url.Values{"city_in": {"SP,RJ"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	// Ana(SP), Bruno(RJ), Carla(SP), Elena(RJ)
	if len(result) != 4 {
		t.Errorf("expected 4 users in SP or RJ, got %d", len(result))
	}
}

func TestApplyBetween(t *testing.T) {
	params := url.Values{"age_between": {"20,25"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	// Ana(20), Carla(25), Elena(22)
	if len(result) != 3 {
		t.Errorf("expected 3 users age 20-25, got %d", len(result))
	}
}

func TestApplySort(t *testing.T) {
	params := url.Values{"sort": {"age"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 5 {
		t.Fatalf("expected 5 users, got %d", len(result))
	}
	if result[0].Name != "Bruno" {
		t.Errorf("expected Bruno (17) first, got %s (%d)", result[0].Name, result[0].Age)
	}
	if result[4].Name != "Daniel" {
		t.Errorf("expected Daniel (30) last, got %s (%d)", result[4].Name, result[4].Age)
	}
}

func TestApplySortDesc(t *testing.T) {
	params := url.Values{"sort": {"-age"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	if result[0].Name != "Daniel" {
		t.Errorf("expected Daniel (30) first, got %s", result[0].Name)
	}
}

func TestApplyPaginatedBasic(t *testing.T) {
	params := url.Values{"sort": {"name"}, "page": {"1"}, "limit": {"2"}}
	page, err := ApplyPaginated(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 5 {
		t.Errorf("expected total 5, got %d", page.Total)
	}
	if len(page.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(page.Items))
	}
	if page.Page != 1 {
		t.Errorf("expected page 1, got %d", page.Page)
	}
	if !page.HasNext {
		t.Error("expected HasNext to be true")
	}
}

func TestApplyPaginatedLastPage(t *testing.T) {
	params := url.Values{"sort": {"name"}, "page": {"3"}, "limit": {"2"}}
	page, err := ApplyPaginated(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 {
		t.Errorf("expected 1 item on last page, got %d", len(page.Items))
	}
	if page.HasNext {
		t.Error("expected HasNext to be false on last page")
	}
}

func TestApplyPaginatedWithFilter(t *testing.T) {
	params := url.Values{"city": {"SP"}, "page": {"1"}, "limit": {"1"}}
	page, err := ApplyPaginated(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 2 {
		t.Errorf("expected total 2 (SP users), got %d", page.Total)
	}
	if len(page.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(page.Items))
	}
	if !page.HasNext {
		t.Error("expected HasNext to be true")
	}
}

func TestApplyWithOptions(t *testing.T) {
	params := url.Values{}
	result, err := Apply(testUsers(), params,
		WithDefaultSort("Age", true),
	)
	if err != nil {
		t.Fatal(err)
	}
	if result[0].Name != "Bruno" {
		t.Errorf("expected Bruno first with default sort by Age asc, got %s", result[0].Name)
	}
}

func TestApplyNe(t *testing.T) {
	params := url.Values{"city_ne": {"SP"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	// Bruno(RJ), Daniel(MG), Elena(RJ)
	if len(result) != 3 {
		t.Errorf("expected 3 users not from SP, got %d", len(result))
	}
}

func TestApplyLte(t *testing.T) {
	params := url.Values{"age_lte": {"20"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	// Ana(20), Bruno(17)
	if len(result) != 2 {
		t.Errorf("expected 2 users with age <= 20, got %d", len(result))
	}
}

func TestApplyGte(t *testing.T) {
	params := url.Values{"age_gte": {"25"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	// Carla(25), Daniel(30)
	if len(result) != 2 {
		t.Errorf("expected 2 users with age >= 25, got %d", len(result))
	}
}

func TestApplyLt(t *testing.T) {
	params := url.Values{"age_lt": {"20"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	// Bruno(17)
	if len(result) != 1 {
		t.Errorf("expected 1 user with age < 20, got %d", len(result))
	}
}

func TestApplyEmptyParams(t *testing.T) {
	params := url.Values{}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 5 {
		t.Errorf("expected all 5 users with empty params, got %d", len(result))
	}
}

func TestApplyMaxLimit(t *testing.T) {
	params := url.Values{"limit": {"500"}}
	_, err := Apply(testUsers(), params, WithMaxLimit(100))
	if err == nil {
		t.Fatal("expected error for limit > maxLimit")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./query/ -run "TestApply" -v`
Expected: FAIL

**Step 3: Write minimal implementation**

```go
// query/query.go
package query

import (
	"net/url"

	"github.com/sidneip/gofilter/filter"
)

type PageResult[T any] struct {
	Items   []T `json:"items"`
	Total   int `json:"total"`
	Page    int `json:"page"`
	Limit   int `json:"limit"`
	HasNext bool `json:"has_next"`
}

type options struct {
	defaultLimit    int
	maxLimit        int
	defaultSort     string
	defaultSortAsc  bool
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
```

**Step 4: Run test to verify it passes**

Run: `go test ./query/ -v`
Expected: ALL PASS

**Step 5: Commit**

```bash
git add query/query.go query/query_test.go
git commit -m "feat(query): add Apply and ApplyPaginated for end-to-end query filtering"
```

---

### Task 6: Run all tests and verify

**Step 1: Run full test suite**

Run: `go test ./... -v`
Expected: ALL PASS (both `filter/` and `query/` packages)

**Step 2: Run with race detector**

Run: `go test ./... -race`
Expected: PASS, no race conditions

**Step 3: Commit if any fixes were needed**

---

### Task 7: Add example

**Files:**
- Create: `examples/query/main.go`

**Step 1: Write example**

```go
// examples/query/main.go
package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/sidneip/gofilter/query"
)

type User struct {
	Name string `json:"name" gofilter:"filterable,sortable,column=name"`
	Age  int    `json:"age" gofilter:"filterable,sortable"`
	City string `json:"city" gofilter:"filterable,sortable"`
}

func main() {
	users := []User{
		{Name: "Ana", Age: 20, City: "SP"},
		{Name: "Bruno", Age: 17, City: "RJ"},
		{Name: "Carla", Age: 25, City: "SP"},
		{Name: "Daniel", Age: 30, City: "MG"},
		{Name: "Elena", Age: 22, City: "RJ"},
	}

	// Simulate: GET /users?city=SP&age_gt=18&sort=-age
	params := url.Values{
		"city":   {"SP"},
		"age_gt": {"18"},
		"sort":   {"-age"},
	}

	fmt.Println("=== Simple Apply ===")
	result, err := query.Apply(users, params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	printJSON(result)

	// Simulate: GET /users?age_gte=18&sort=name&page=1&limit=2
	params2 := url.Values{
		"age_gte": {"18"},
		"sort":    {"name"},
		"page":    {"1"},
		"limit":   {"2"},
	}

	fmt.Println("\n=== Paginated Apply ===")
	page, err := query.ApplyPaginated(users, params2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	printJSON(page)
}

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}
```

**Step 2: Run example**

Run: `go run examples/query/main.go`
Expected: JSON output showing filtered users

**Step 3: Commit**

```bash
git add examples/query/main.go
git commit -m "feat: add query parsing example"
```
