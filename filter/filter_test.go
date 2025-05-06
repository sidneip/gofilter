package filter

import (
	"testing"
)

type Person struct {
	Name    string
	Age     int
	Hobbies []string
	Address *Address
}

type Address struct {
	City    string
	Country string
}

func TestEq(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	// Test equals for string field
	result := Apply(people, Eq[Person]("Name", "Alice"))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 person named Alice, got %d", len(result))
	}

	// Test equals for int field
	result = Apply(people, Eq[Person]("Age", 25))
	if len(result) != 1 || result[0].Name != "Bob" {
		t.Errorf("Expected 1 person aged 25, got %d", len(result))
	}

	// Test equals with no match
	result = Apply(people, Eq[Person]("Name", "Charlie"))
	if len(result) != 0 {
		t.Errorf("Expected 0 people, got %d", len(result))
	}
}

func TestGt(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}

	// Test greater than
	result := Apply(people, Gt[Person]("Age", 28))
	if len(result) != 2 {
		t.Errorf("Expected 2 people older than 28, got %d", len(result))
	}
}

func TestLt(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}

	// Test less than
	result := Apply(people, Lt[Person]("Age", 30))
	if len(result) != 1 || result[0].Name != "Bob" {
		t.Errorf("Expected 1 person younger than 30, got %d", len(result))
	}
}

func TestContains(t *testing.T) {
	people := []Person{
		{Name: "Alice", Hobbies: []string{"reading", "swimming"}},
		{Name: "Bob", Hobbies: []string{"gaming"}},
	}

	// Test contains for slice
	result := Apply(people, Contains[Person]("Hobbies", "reading"))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 person with reading hobby, got %d", len(result))
	}

	// Test contains for string
	result = Apply(people, Contains[Person]("Name", "li"))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 person with 'li' in name, got %d", len(result))
	}
}

func TestAnd(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 30, Hobbies: []string{"reading", "swimming"}},
		{Name: "Bob", Age: 25, Hobbies: []string{"gaming"}},
		{Name: "Charlie", Age: 35, Hobbies: []string{"reading", "hiking"}},
	}

	// Test And operator
	result := Apply(people, And[Person](
		Gt[Person]("Age", 25),
		Contains[Person]("Hobbies", "reading"),
	))

	if len(result) != 2 {
		t.Errorf("Expected 2 people older than 25 who like reading, got %d", len(result))
	}

	for _, p := range result {
		if p.Age <= 25 {
			t.Errorf("Expected all people to be older than 25, got %d", p.Age)
		}

		hasReading := false
		for _, h := range p.Hobbies {
			if h == "reading" {
				hasReading = true
				break
			}
		}

		if !hasReading {
			t.Errorf("Expected all people to have reading as a hobby")
		}
	}
}

func TestOr(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 30, Hobbies: []string{"reading", "swimming"}},
		{Name: "Bob", Age: 25, Hobbies: []string{"gaming"}},
		{Name: "Charlie", Age: 35, Hobbies: []string{"hiking"}},
	}

	// Test Or operator
	result := Apply(people, Or[Person](
		Eq[Person]("Name", "Alice"),
		Contains[Person]("Hobbies", "hiking"),
	))

	if len(result) != 2 {
		t.Errorf("Expected 2 people (Alice or hikers), got %d", len(result))
	}

	names := make(map[string]bool)
	for _, p := range result {
		names[p.Name] = true
	}

	if !names["Alice"] || !names["Charlie"] {
		t.Errorf("Expected Alice and Charlie, got different people")
	}
}

func TestNested(t *testing.T) {
	people := []Person{
		{
			Name:    "Alice",
			Age:     30,
			Address: &Address{City: "New York", Country: "USA"},
		},
		{
			Name:    "Bob",
			Age:     25,
			Address: &Address{City: "London", Country: "UK"},
		},
	}

	// Test nested field access
	result := Apply(people, Eq[Person]("Address.Country", "USA"))
	if len(result) != 1 || result[0].Name != "Alice" {
		t.Errorf("Expected 1 person from USA, got %d", len(result))
	}
}

func TestNot(t *testing.T) {
	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}

	// Test Not operator
	result := Apply(people, Not[Person](Eq[Person]("Name", "Bob")))
	if len(result) != 2 {
		t.Errorf("Expected 2 people not named Bob, got %d", len(result))
	}

	for _, p := range result {
		if p.Name == "Bob" {
			t.Errorf("Expected no one named Bob, but found Bob")
		}
	}
}
