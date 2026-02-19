package query

import (
	"testing"
	"time"
)

type TagTestUser struct {
	Name    string    `gofilter:"filterable,sortable,column=name"`
	Age     int       `gofilter:"filterable"`
	City    string    `gofilter:"filterable,sortable"`
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

	if len(registry.fields) != 6 {
		t.Errorf("expected 6 fields, got %d", len(registry.fields))
	}

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
