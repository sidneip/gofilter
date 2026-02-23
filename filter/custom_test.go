package filter

import (
	"reflect"
	"testing"
	"time"
)

type TestUser struct {
	Name      string
	Age       int
	Email     string
	Active    bool
	Score     float64
	Tags      []string
	CreatedAt time.Time
	DeletedAt *time.Time
	Metadata  map[string]interface{}
}

// TestNe tests the Ne (not equal) filter
func TestNe(t *testing.T) {
	users := []TestUser{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 30},
	}

	result := Apply(users, Ne[TestUser]("Age", 30))
	if len(result) != 1 || result[0].Name != "Bob" {
		t.Errorf("Expected 1 user not aged 30, got %d", len(result))
	}

	// Test Ne with string
	result = Apply(users, Ne[TestUser]("Name", "Alice"))
	if len(result) != 2 {
		t.Errorf("Expected 2 users not named Alice, got %d", len(result))
	}
}

// TestIsNil tests the IsNil filter
func TestIsNil(t *testing.T) {
	type UserWithSlice struct {
		Name  string
		Tags  []string
		Extra map[string]string
	}

	users := []UserWithSlice{
		{Name: "Alice", Tags: nil, Extra: nil},
		{Name: "Bob", Tags: []string{"admin"}, Extra: map[string]string{"key": "value"}},
	}

	// Test nil slice
	result := Apply(users, IsNil[UserWithSlice]("Tags"))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 user with nil Tags, got %d", len(result))
	}

	// Test nil map
	result = Apply(users, IsNil[UserWithSlice]("Extra"))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 user with nil Extra, got %d", len(result))
	}
}

// TestIsNotNil tests the IsNotNil filter
func TestIsNotNil(t *testing.T) {
	type UserWithSlice struct {
		Name string
		Tags []string
	}

	users := []UserWithSlice{
		{Name: "Alice", Tags: nil},
		{Name: "Bob", Tags: []string{"admin"}},
	}

	result := Apply(users, IsNotNil[UserWithSlice]("Tags"))
	if len(result) != 1 || result[0].Name != "Bob" {
		t.Errorf("Expected 1 user with non-nil Tags, got %d", len(result))
	}
}

// TestIsZero tests the IsZero filter
func TestIsZero(t *testing.T) {
	users := []TestUser{
		{Name: "Alice", Score: 0},
		{Name: "Bob", Score: 85.5},
		{Name: "Charlie", Score: 0},
	}

	result := Apply(users, IsZero[TestUser]("Score"))
	if len(result) != 2 {
		t.Errorf("Expected 2 users with zero score, got %d", len(result))
	}

	// Test IsZero with string
	users = []TestUser{
		{Name: "Alice", Email: ""},
		{Name: "Bob", Email: "bob@test.com"},
	}

	result = Apply(users, IsZero[TestUser]("Email"))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 user with empty email, got %d", len(result))
	}
}

// TestIsNotZero tests the IsNotZero filter
func TestIsNotZero(t *testing.T) {
	users := []TestUser{
		{Name: "Alice", Score: 0},
		{Name: "Bob", Score: 85.5},
	}

	result := Apply(users, IsNotZero[TestUser]("Score"))
	if len(result) != 1 || result[0].Name != "Bob" {
		t.Errorf("Expected 1 user with non-zero score, got %d", len(result))
	}
}

// TestStringMatch tests the StringMatch filter with various modes
func TestStringMatch(t *testing.T) {
	users := []TestUser{
		{Name: "Alice", Email: "alice@gmail.com"},
		{Name: "Bob", Email: "bob@GMAIL.COM"},
		{Name: "Charlie", Email: "charlie@yahoo.com"},
	}

	// Exact match
	result := Apply(users, StringMatch[TestUser]("Email", "alice@gmail.com", StringMatchOptions{Mode: ExactMatch}))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 user with exact email match, got %d", len(result))
	}

	// Exact match case insensitive
	result = Apply(users, StringMatch[TestUser]("Email", "ALICE@GMAIL.COM", StringMatchOptions{Mode: ExactMatch, IgnoreCase: true}))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 user with case-insensitive exact match, got %d", len(result))
	}

	// Contains match
	result = Apply(users, StringMatch[TestUser]("Email", "gmail", StringMatchOptions{Mode: ContainsMatch}))
	if len(result) != 1 {
		t.Errorf("Expected 1 user with gmail in email (case sensitive), got %d", len(result))
	}

	// Contains match case insensitive
	result = Apply(users, StringMatch[TestUser]("Email", "GMAIL", StringMatchOptions{Mode: ContainsMatch, IgnoreCase: true}))
	if len(result) != 2 {
		t.Errorf("Expected 2 users with gmail in email (case insensitive), got %d", len(result))
	}

	// Prefix match
	result = Apply(users, StringMatch[TestUser]("Email", "alice", StringMatchOptions{Mode: PrefixMatch}))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 user with email starting with alice, got %d", len(result))
	}

	// Prefix match case insensitive
	result = Apply(users, StringMatch[TestUser]("Email", "ALICE", StringMatchOptions{Mode: PrefixMatch, IgnoreCase: true}))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 user with email starting with ALICE (case insensitive), got %d", len(result))
	}

	// Suffix match
	result = Apply(users, StringMatch[TestUser]("Email", "yahoo.com", StringMatchOptions{Mode: SuffixMatch}))
	if len(result) != 1 || result[0].Name != "Charlie" {
		t.Errorf("Expected 1 user with email ending with yahoo.com, got %d", len(result))
	}

	// Suffix match case insensitive
	result = Apply(users, StringMatch[TestUser]("Email", "YAHOO.COM", StringMatchOptions{Mode: SuffixMatch, IgnoreCase: true}))
	if len(result) != 1 || result[0].Name != "Charlie" {
		t.Errorf("Expected 1 user with email ending with YAHOO.COM (case insensitive), got %d", len(result))
	}

	// Invalid mode (should return false)
	result = Apply(users, StringMatch[TestUser]("Email", "test", StringMatchOptions{Mode: StringMatchMode(99)}))
	if len(result) != 0 {
		t.Errorf("Expected 0 users with invalid mode, got %d", len(result))
	}

	// Non-string field (should return false)
	result = Apply(users, StringMatch[TestUser]("Age", "30", StringMatchOptions{Mode: ExactMatch}))
	if len(result) != 0 {
		t.Errorf("Expected 0 users when matching non-string field, got %d", len(result))
	}
}

// TestArrayContains tests the ArrayContains filter
func TestArrayContains(t *testing.T) {
	users := []TestUser{
		{Name: "Alice", Tags: []string{"admin", "active"}},
		{Name: "Bob", Tags: []string{"user", "ACTIVE"}},
		{Name: "Charlie", Tags: []string{"guest"}},
	}

	// Case sensitive
	result := Apply(users, ArrayContains[TestUser]("Tags", "active", false))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 user with 'active' tag (case sensitive), got %d", len(result))
	}

	// Case insensitive
	result = Apply(users, ArrayContains[TestUser]("Tags", "active", true))
	if len(result) != 2 {
		t.Errorf("Expected 2 users with 'active' tag (case insensitive), got %d", len(result))
	}

	// Non-array field
	result = Apply(users, ArrayContains[TestUser]("Name", "test", false))
	if len(result) != 0 {
		t.Errorf("Expected 0 users when checking non-array field, got %d", len(result))
	}
}

// TestArrayContainsAny tests the ArrayContainsAny filter
func TestArrayContainsAny(t *testing.T) {
	users := []TestUser{
		{Name: "Alice", Tags: []string{"admin", "active"}},
		{Name: "Bob", Tags: []string{"user"}},
		{Name: "Charlie", Tags: []string{"guest", "inactive"}},
	}

	result := Apply(users, ArrayContainsAny[TestUser]("Tags", []interface{}{"admin", "guest"}))
	if len(result) != 2 {
		t.Errorf("Expected 2 users with admin or guest tag, got %d", len(result))
	}

	// No match
	result = Apply(users, ArrayContainsAny[TestUser]("Tags", []interface{}{"superadmin"}))
	if len(result) != 0 {
		t.Errorf("Expected 0 users with superadmin tag, got %d", len(result))
	}
}

// TestArrayContainsAll tests the ArrayContainsAll filter
func TestArrayContainsAll(t *testing.T) {
	users := []TestUser{
		{Name: "Alice", Tags: []string{"admin", "active", "premium"}},
		{Name: "Bob", Tags: []string{"admin", "active"}},
		{Name: "Charlie", Tags: []string{"guest"}},
	}

	result := Apply(users, ArrayContainsAll[TestUser]("Tags", []interface{}{"admin", "active"}))
	if len(result) != 2 {
		t.Errorf("Expected 2 users with both admin and active tags, got %d", len(result))
	}

	result = Apply(users, ArrayContainsAll[TestUser]("Tags", []interface{}{"admin", "active", "premium"}))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 user with all three tags, got %d", len(result))
	}

	// Empty array should return false
	users = []TestUser{{Name: "Empty", Tags: []string{}}}
	result = Apply(users, ArrayContainsAll[TestUser]("Tags", []interface{}{"any"}))
	if len(result) != 0 {
		t.Errorf("Expected 0 users with empty tags array, got %d", len(result))
	}
}

// TestDateBefore tests the DateBefore filter
func TestDateBefore(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	users := []TestUser{
		{Name: "Alice", CreatedAt: yesterday},
		{Name: "Bob", CreatedAt: tomorrow},
	}

	result := Apply(users, DateBefore[TestUser]("CreatedAt", now))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 user created before now, got %d", len(result))
	}
}

// TestDateAfter tests the DateAfter filter
func TestDateAfter(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	users := []TestUser{
		{Name: "Alice", CreatedAt: yesterday},
		{Name: "Bob", CreatedAt: tomorrow},
	}

	result := Apply(users, DateAfter[TestUser]("CreatedAt", now))
	if len(result) != 1 || result[0].Name != "Bob" {
		t.Errorf("Expected 1 user created after now, got %d", len(result))
	}
}

// TestDateBetween tests the DateBetween filter
func TestDateBetween(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	twoDaysAgo := now.Add(-48 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	users := []TestUser{
		{Name: "Alice", CreatedAt: twoDaysAgo},
		{Name: "Bob", CreatedAt: yesterday},
		{Name: "Charlie", CreatedAt: tomorrow},
	}

	result := Apply(users, DateBetween[TestUser]("CreatedAt", twoDaysAgo, now))
	if len(result) != 2 {
		t.Errorf("Expected 2 users created between two days ago and now, got %d", len(result))
	}
}

// TestDateFilterWithStringField tests date filters with string date fields
func TestDateFilterWithStringField(t *testing.T) {
	type Event struct {
		Name string
		Date string
	}

	events := []Event{
		{Name: "Event1", Date: "2024-01-15"},
		{Name: "Event2", Date: "2024-06-15"},
		{Name: "Event3", Date: "invalid-date"},
	}

	cutoff, _ := time.Parse("2006-01-02", "2024-03-01")
	result := Apply(events, DateBefore[Event]("Date", cutoff))
	if len(result) != 1 || result[0].Name != "Event1" {
		t.Errorf("Expected 1 event before March 2024, got %d", len(result))
	}

	result = Apply(events, DateAfter[Event]("Date", cutoff))
	if len(result) != 1 || result[0].Name != "Event2" {
		t.Errorf("Expected 1 event after March 2024, got %d", len(result))
	}
}

// TestRegexMatch tests the RegexMatch filter
func TestRegexMatch(t *testing.T) {
	users := []TestUser{
		{Name: "Alice", Email: "alice@gmail.com"},
		{Name: "Bob", Email: "bob@yahoo.com"},
		{Name: "Charlie", Email: "charlie123@gmail.com"},
	}

	// Match gmail emails
	result := Apply(users, RegexMatch[TestUser]("Email", `^[a-z]+@gmail\.com$`))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 user with simple gmail email, got %d", len(result))
	}

	// Match emails with numbers
	result = Apply(users, RegexMatch[TestUser]("Email", `\d+`))
	if len(result) != 1 || result[0].Name != "Charlie" {
		t.Errorf("Expected 1 user with numbers in email, got %d", len(result))
	}

	// Invalid regex (should return empty)
	result = Apply(users, RegexMatch[TestUser]("Email", `[invalid`))
	if len(result) != 0 {
		t.Errorf("Expected 0 users with invalid regex, got %d", len(result))
	}

	// Non-string field
	result = Apply(users, RegexMatch[TestUser]("Age", `\d+`))
	if len(result) != 0 {
		t.Errorf("Expected 0 users when matching non-string field, got %d", len(result))
	}
}

// TestNestedArrayAny tests the NestedArrayAny filter
func TestNestedArrayAny(t *testing.T) {
	users := []TestUser{
		{Name: "Alice", Tags: []string{"admin", "active"}},
		{Name: "Bob", Tags: []string{"user", "approved"}},
		{Name: "Charlie", Tags: []string{"guest"}},
	}

	// Check if any tag starts with 'a'
	result := Apply(users, NestedArrayAny[TestUser]("Tags", func(elem reflect.Value) bool {
		if elem.Kind() == reflect.String {
			s := elem.String()
			return len(s) > 0 && s[0] == 'a'
		}
		return false
	}))

	// Alice has "admin" and "active", Bob has "approved"
	if len(result) != 2 {
		t.Errorf("Expected 2 users with tags starting with 'a', got %d", len(result))
	}

	// Non-slice field
	result = Apply(users, NestedArrayAny[TestUser]("Name", func(elem reflect.Value) bool {
		return true
	}))
	if len(result) != 0 {
		t.Errorf("Expected 0 users when checking non-slice field, got %d", len(result))
	}
}

// TestNestedArrayAll tests the NestedArrayAll filter
func TestNestedArrayAll(t *testing.T) {
	users := []TestUser{
		{Name: "Alice", Tags: []string{"active", "approved"}},
		{Name: "Bob", Tags: []string{"active", "pending"}},
		{Name: "Charlie", Tags: []string{}},
	}

	// Check if all tags start with 'a'
	result := Apply(users, NestedArrayAll[TestUser]("Tags", func(elem reflect.Value) bool {
		if elem.Kind() == reflect.String {
			s := elem.String()
			return len(s) > 0 && s[0] == 'a'
		}
		return false
	}))

	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 user with all tags starting with 'a', got %d", len(result))
	}

	// Empty array should return false
	result = Apply(users, NestedArrayAll[TestUser]("Tags", func(elem reflect.Value) bool {
		return true
	}))
	// Alice and Bob have tags, Charlie has empty
	if len(result) != 2 {
		t.Errorf("Expected 2 users with non-empty tags, got %d", len(result))
	}
}

// TestExportedGetFieldValue tests the ExportedGetFieldValue function
func TestExportedGetFieldValue(t *testing.T) {
	user := TestUser{
		Name:  "Alice",
		Age:   30,
		Email: "alice@test.com",
	}

	// Get string field
	val, err := ExportedGetFieldValue(user, "Name")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if val.String() != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", val.String())
	}

	// Get int field
	val, err = ExportedGetFieldValue(user, "Age")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if val.Int() != 30 {
		t.Errorf("Expected 30, got %d", val.Int())
	}

	// Get non-existent field
	_, err = ExportedGetFieldValue(user, "NonExistent")
	if err == nil {
		t.Errorf("Expected error for non-existent field")
	}

	// Get field from non-struct
	_, err = ExportedGetFieldValue("not a struct", "Field")
	if err == nil {
		t.Errorf("Expected error for non-struct value")
	}
}

// TestSortWithErrors tests Sort function edge cases
func TestSortWithErrors(t *testing.T) {
	users := []TestUser{
		{Name: "Charlie", Age: 35},
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	// Sort descending
	result := Sort(users, "Age", false)
	if result[0].Name != "Charlie" || result[2].Name != "Bob" {
		t.Errorf("Expected descending sort by age")
	}

	// Sort by non-existent field (should not crash)
	result = Sort(users, "NonExistent", true)
	if len(result) != 3 {
		t.Errorf("Expected 3 users even with invalid field")
	}
}

// TestBetweenWithStrings tests Between with string values
func TestBetweenWithStrings(t *testing.T) {
	users := []TestUser{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Charlie"},
		{Name: "Diana"},
	}

	result := Apply(users, Between[TestUser]("Name", "B", "D"))
	if len(result) != 2 {
		t.Errorf("Expected 2 users with names between B and D, got %d", len(result))
	}
}

// TestCustomFilter tests the Custom filter
func TestCustomFilter(t *testing.T) {
	users := []TestUser{
		{Name: "Alice", Age: 30, Active: true},
		{Name: "Bob", Age: 25, Active: false},
		{Name: "Charlie", Age: 35, Active: true},
	}

	result := Apply(users, Custom(func(u TestUser) bool {
		return u.Active && u.Age > 28
	}))

	if len(result) != 2 {
		t.Errorf("Expected 2 active users over 28, got %d", len(result))
	}
}

// TestCompareValuesWithDifferentTypes tests compareValues with various types
func TestCompareValuesWithDifferentTypes(t *testing.T) {
	type Item struct {
		IntVal    int
		Int8Val   int8
		Int16Val  int16
		Int32Val  int32
		Int64Val  int64
		UintVal   uint
		Uint8Val  uint8
		Uint16Val uint16
		Uint32Val uint32
		Uint64Val uint64
		FloatVal  float32
		Float64   float64
		BoolVal   bool
		StringVal string
	}

	items := []Item{
		{IntVal: 10, Int8Val: 10, Int16Val: 10, Int32Val: 10, Int64Val: 10,
			UintVal: 10, Uint8Val: 10, Uint16Val: 10, Uint32Val: 10, Uint64Val: 10,
			FloatVal: 10.5, Float64: 10.5, BoolVal: true, StringVal: "test"},
		{IntVal: 20, Int8Val: 20, Int16Val: 20, Int32Val: 20, Int64Val: 20,
			UintVal: 20, Uint8Val: 20, Uint16Val: 20, Uint32Val: 20, Uint64Val: 20,
			FloatVal: 20.5, Float64: 20.5, BoolVal: false, StringVal: "other"},
	}

	// Test int8
	result := Apply(items, Eq[Item]("Int8Val", int8(10)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with Int8Val=10, got %d", len(result))
	}

	// Test int16
	result = Apply(items, Eq[Item]("Int16Val", int16(10)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with Int16Val=10, got %d", len(result))
	}

	// Test int32
	result = Apply(items, Eq[Item]("Int32Val", int32(10)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with Int32Val=10, got %d", len(result))
	}

	// Test uint
	result = Apply(items, Eq[Item]("UintVal", uint(10)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with UintVal=10, got %d", len(result))
	}

	// Test uint8
	result = Apply(items, Eq[Item]("Uint8Val", uint8(10)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with Uint8Val=10, got %d", len(result))
	}

	// Test uint16
	result = Apply(items, Eq[Item]("Uint16Val", uint16(10)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with Uint16Val=10, got %d", len(result))
	}

	// Test uint32
	result = Apply(items, Eq[Item]("Uint32Val", uint32(10)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with Uint32Val=10, got %d", len(result))
	}

	// Test uint64
	result = Apply(items, Eq[Item]("Uint64Val", uint64(10)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with Uint64Val=10, got %d", len(result))
	}

	// Test float32
	result = Apply(items, Eq[Item]("FloatVal", float32(10.5)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with FloatVal=10.5, got %d", len(result))
	}

	// Test bool
	result = Apply(items, Eq[Item]("BoolVal", true))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with BoolVal=true, got %d", len(result))
	}

	// Test less than comparisons with different types
	result = Apply(items, Lt[Item]("Int8Val", int8(15)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with Int8Val<15, got %d", len(result))
	}

	result = Apply(items, Lt[Item]("UintVal", uint(15)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with UintVal<15, got %d", len(result))
	}

	result = Apply(items, Lt[Item]("FloatVal", float32(15.0)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with FloatVal<15, got %d", len(result))
	}

	// Test string comparison
	result = Apply(items, Lt[Item]("StringVal", "u"))
	if len(result) != 2 {
		t.Errorf("Expected 2 items with StringVal<'u', got %d", len(result))
	}
}

// TestGteAndLteOperators tests Gte and Lte with various scenarios
func TestGteAndLteOperators(t *testing.T) {
	type Item struct {
		Value int
	}

	items := []Item{{Value: 10}, {Value: 20}, {Value: 30}}

	// Gte - should include exact match
	result := Apply(items, Gte[Item]("Value", 20))
	if len(result) != 2 {
		t.Errorf("Expected 2 items with Value>=20, got %d", len(result))
	}

	// Lte - should include exact match
	result = Apply(items, Lte[Item]("Value", 20))
	if len(result) != 2 {
		t.Errorf("Expected 2 items with Value<=20, got %d", len(result))
	}
}

// TestInvalidFieldAccess tests error handling for invalid field access
func TestInvalidFieldAccess(t *testing.T) {
	users := []TestUser{{Name: "Alice"}}

	// Non-existent field
	result := Apply(users, Eq[TestUser]("NonExistent", "value"))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for non-existent field, got %d", len(result))
	}

	// Gte on non-existent field
	result = Apply(users, Gte[TestUser]("NonExistent", 10))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for Gte on non-existent field, got %d", len(result))
	}

	// Lte on non-existent field
	result = Apply(users, Lte[TestUser]("NonExistent", 10))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for Lte on non-existent field, got %d", len(result))
	}
}

// TestGetFieldValueEdgeCases tests edge cases for getFieldValue
func TestGetFieldValueEdgeCases(t *testing.T) {
	// Test nested field with pointer
	type Address struct {
		City    string
		Country string
	}
	type PersonWithAddress struct {
		Name    string
		Address *Address
	}

	// Test with nil nested pointer
	people := []PersonWithAddress{
		{Name: "Alice", Address: nil},
		{Name: "Bob", Address: &Address{City: "NYC", Country: "USA"}},
	}

	// Should return 0 for nil nested pointer
	result := Apply(people, Eq[PersonWithAddress]("Address.City", "NYC"))
	if len(result) != 1 || result[0].Name != "Bob" {
		t.Errorf("Expected 1 person with Address.City=NYC, got %d", len(result))
	}

	// Test nested field access with non-pointer
	type PersonEmbedded struct {
		Name    string
		Address Address
	}

	embedded := []PersonEmbedded{
		{Name: "Alice", Address: Address{City: "LA", Country: "USA"}},
		{Name: "Bob", Address: Address{City: "NYC", Country: "USA"}},
	}

	result2 := Apply(embedded, Eq[PersonEmbedded]("Address.City", "NYC"))
	if len(result2) != 1 || result2[0].Name != "Bob" {
		t.Errorf("Expected 1 person with embedded Address.City=NYC, got %d", len(result2))
	}

	// Test nested path where intermediate is not a struct
	type PersonWithString struct {
		Name   string
		Simple string
	}

	simples := []PersonWithString{{Name: "Alice", Simple: "test"}}
	result3 := Apply(simples, Eq[PersonWithString]("Simple.Something", "value"))
	if len(result3) != 0 {
		t.Errorf("Expected 0 results for invalid nested path, got %d", len(result3))
	}
}

// TestCompareValuesTypeConversion tests type conversion in compareValues
func TestCompareValuesTypeConversion(t *testing.T) {
	type Item struct {
		Int64Val int64
		IntVal   int
	}

	items := []Item{
		{Int64Val: 100, IntVal: 100},
		{Int64Val: 200, IntVal: 200},
	}

	// int64 compared with int
	result := Apply(items, Eq[Item]("Int64Val", int64(100)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with Int64Val=100, got %d", len(result))
	}
}

// TestCompareValuesUnsupportedType tests unsupported types in compareValues
func TestCompareValuesUnsupportedType(t *testing.T) {
	type Item struct {
		Data []byte
	}

	items := []Item{
		{Data: []byte("test")},
	}

	// Comparing []byte is unsupported
	result := Apply(items, Eq[Item]("Data", []byte("test")))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for unsupported type comparison, got %d", len(result))
	}
}

// TestCompareValuesLessUnsupportedType tests unsupported types in compareValuesLess
func TestCompareValuesLessUnsupportedType(t *testing.T) {
	type Item struct {
		Data []byte
	}

	items := []Item{
		{Data: []byte("a")},
		{Data: []byte("b")},
	}

	// Lt on []byte is unsupported
	result := Apply(items, Lt[Item]("Data", []byte("c")))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for Lt on unsupported type, got %d", len(result))
	}

	// Gt on []byte is unsupported
	result = Apply(items, Gt[Item]("Data", []byte("a")))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for Gt on unsupported type, got %d", len(result))
	}
}

// TestCompareValuesBool tests bool comparison
func TestCompareValuesBool(t *testing.T) {
	type Item struct {
		Active bool
	}

	items := []Item{
		{Active: true},
		{Active: false},
	}

	result := Apply(items, Eq[Item]("Active", true))
	if len(result) != 1 || !items[0].Active {
		t.Errorf("Expected 1 active item, got %d", len(result))
	}

	result = Apply(items, Eq[Item]("Active", false))
	if len(result) != 1 {
		t.Errorf("Expected 1 inactive item, got %d", len(result))
	}
}

// TestCompareValuesUint tests uint comparison
func TestCompareValuesUint(t *testing.T) {
	type Item struct {
		Count uint
	}

	items := []Item{
		{Count: 10},
		{Count: 20},
		{Count: 30},
	}

	result := Apply(items, Eq[Item]("Count", uint(20)))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with Count=20, got %d", len(result))
	}

	result = Apply(items, Lt[Item]("Count", uint(25)))
	if len(result) != 2 {
		t.Errorf("Expected 2 items with Count<25, got %d", len(result))
	}

	result = Apply(items, Gt[Item]("Count", uint(15)))
	if len(result) != 2 {
		t.Errorf("Expected 2 items with Count>15, got %d", len(result))
	}
}

// TestCompareValuesStringComparison tests string lexicographic comparison
func TestCompareValuesStringComparison(t *testing.T) {
	type Item struct {
		Name string
	}

	items := []Item{
		{Name: "apple"},
		{Name: "banana"},
		{Name: "cherry"},
	}

	// String less than comparison
	result := Apply(items, Lt[Item]("Name", "banana"))
	if len(result) != 1 || result[0].Name != "apple" {
		t.Errorf("Expected 1 item with Name<'banana', got %d", len(result))
	}

	// String greater than comparison
	result = Apply(items, Gt[Item]("Name", "banana"))
	if len(result) != 1 || result[0].Name != "cherry" {
		t.Errorf("Expected 1 item with Name>'banana', got %d", len(result))
	}
}

// TestContainsNonStringValue tests Contains with non-string value on string field
func TestContainsNonStringValue(t *testing.T) {
	type Item struct {
		Name string
	}

	items := []Item{
		{Name: "test123"},
	}

	// Contains with int on string field should return false
	result := Apply(items, Contains[Item]("Name", 123))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for Contains with non-string value, got %d", len(result))
	}
}

// TestContainsOnNonStringNonSlice tests Contains on non-string, non-slice field
func TestContainsOnNonStringNonSlice(t *testing.T) {
	type Item struct {
		Count int
	}

	items := []Item{
		{Count: 123},
	}

	// Contains on int field should return false
	result := Apply(items, Contains[Item]("Count", 1))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for Contains on int field, got %d", len(result))
	}
}

// TestInWithIncompatibleTypes tests In with incompatible types
func TestInWithIncompatibleTypes(t *testing.T) {
	type Item struct {
		Value int
	}

	items := []Item{{Value: 10}}

	// In with incompatible type in values list
	result := Apply(items, In[Item]("Value", []interface{}{"abc", "def"}))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for In with incompatible types, got %d", len(result))
	}
}

// TestNeWithIncompatibleType tests Ne with incompatible types
func TestNeWithIncompatibleType(t *testing.T) {
	type Item struct {
		Value int
	}

	items := []Item{{Value: 10}}

	// Ne with incompatible type
	result := Apply(items, Ne[Item]("Value", "string"))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for Ne with incompatible type, got %d", len(result))
	}
}

// TestGtLtWithIncompatibleType tests Gt/Lt with incompatible types
func TestGtLtWithIncompatibleType(t *testing.T) {
	type Item struct {
		Value int
	}

	items := []Item{{Value: 10}}

	// Gt with incompatible type
	result := Apply(items, Gt[Item]("Value", "string"))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for Gt with incompatible type, got %d", len(result))
	}

	// Lt with incompatible type
	result = Apply(items, Lt[Item]("Value", "string"))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for Lt with incompatible type, got %d", len(result))
	}
}

// TestIsNilWithPointerField tests IsNil with pointer field
// Note: getFieldValue dereferences pointers, so IsNil on pointer fields
// returns false because getFieldValue returns an error for nil pointers
func TestIsNilWithPointerField(t *testing.T) {
	type Item struct {
		Data *string
	}

	s := "value"
	items := []Item{
		{Data: nil},
		{Data: &s},
	}

	// IsNil on pointer field returns false because getFieldValue
	// returns an error when trying to access nil pointer
	result := Apply(items, IsNil[Item]("Data"))
	if len(result) != 0 {
		t.Errorf("Expected 0 items because getFieldValue errors on nil pointer, got %d", len(result))
	}
}

// TestIsNilWithNonNilableField tests IsNil with non-nilable field
func TestIsNilWithNonNilableField(t *testing.T) {
	type Item struct {
		Value int
	}

	items := []Item{{Value: 0}, {Value: 10}}

	// Int is not nilable, should return false
	result := Apply(items, IsNil[Item]("Value"))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for IsNil on non-nilable field, got %d", len(result))
	}
}

// TestIsZeroWithVariousTypes tests IsZero with different types
func TestIsZeroWithVariousTypes(t *testing.T) {
	type Item struct {
		IntVal    int
		UintVal   uint
		FloatVal  float64
		StringVal string
		BoolVal   bool
	}

	items := []Item{
		{IntVal: 0, UintVal: 0, FloatVal: 0, StringVal: "", BoolVal: false},
		{IntVal: 1, UintVal: 1, FloatVal: 1.0, StringVal: "a", BoolVal: true},
	}

	result := Apply(items, IsZero[Item]("IntVal"))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with zero IntVal, got %d", len(result))
	}

	result = Apply(items, IsZero[Item]("UintVal"))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with zero UintVal, got %d", len(result))
	}

	result = Apply(items, IsZero[Item]("FloatVal"))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with zero FloatVal, got %d", len(result))
	}

	result = Apply(items, IsZero[Item]("StringVal"))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with zero StringVal, got %d", len(result))
	}

	result = Apply(items, IsZero[Item]("BoolVal"))
	if len(result) != 1 {
		t.Errorf("Expected 1 item with zero BoolVal, got %d", len(result))
	}
}
