package query

import (
	"net/url"
	"testing"
)

type ParserTestUser struct {
	Name string `gofilter:"filterable,sortable,column=name"`
	Age  int    `gofilter:"filterable,sortable"`
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
		param   string
		wantOp  string
		wantCol string
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
