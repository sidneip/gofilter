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
	// Only "Daniel" contains "an" (case-sensitive)
	if len(result) != 1 {
		t.Errorf("expected 1 user with 'an' in name, got %d", len(result))
	}
	if len(result) > 0 && result[0].Name != "Daniel" {
		t.Errorf("expected Daniel, got %s", result[0].Name)
	}
}

func TestApplyIn(t *testing.T) {
	params := url.Values{"city_in": {"SP,RJ"}}
	result, err := Apply(testUsers(), params)
	if err != nil {
		t.Fatal(err)
	}
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
