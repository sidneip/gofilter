package filter_test

import (
	"fmt"

	"github.com/sidneip/gofilter/filter"
)

type User struct {
	Name   string
	Age    int
	City   string
	Active bool
}

var users = []User{
	{Name: "Ana", Age: 25, City: "SP", Active: true},
	{Name: "Bob", Age: 30, City: "RJ", Active: false},
	{Name: "Carlos", Age: 35, City: "SP", Active: true},
	{Name: "Diana", Age: 28, City: "MG", Active: true},
}

func ExampleApply() {
	// Filter users from São Paulo
	result := filter.Apply(users, filter.Eq[User]("City", "SP"))

	for _, u := range result {
		fmt.Printf("%s (%s)\n", u.Name, u.City)
	}
	// Output:
	// Ana (SP)
	// Carlos (SP)
}

func ExampleAnd() {
	// Filter active users over 25
	result := filter.Apply(users, filter.And(
		filter.Gt[User]("Age", 25),
		filter.Eq[User]("Active", true),
	))

	for _, u := range result {
		fmt.Printf("%s: %d\n", u.Name, u.Age)
	}
	// Output:
	// Carlos: 35
	// Diana: 28
}

func ExampleOr() {
	// Filter users from SP or RJ
	result := filter.Apply(users, filter.Or(
		filter.Eq[User]("City", "SP"),
		filter.Eq[User]("City", "RJ"),
	))

	for _, u := range result {
		fmt.Printf("%s (%s)\n", u.Name, u.City)
	}
	// Output:
	// Ana (SP)
	// Bob (RJ)
	// Carlos (SP)
}

func ExampleContains() {
	// Filter users with "an" in their name
	result := filter.Apply(users, filter.Contains[User]("Name", "an"))

	for _, u := range result {
		fmt.Println(u.Name)
	}
	// Output:
	// Diana
}

func ExampleIn() {
	// Filter users from specific cities
	result := filter.Apply(users, filter.In[User]("City", []interface{}{"SP", "MG"}))

	for _, u := range result {
		fmt.Printf("%s (%s)\n", u.Name, u.City)
	}
	// Output:
	// Ana (SP)
	// Carlos (SP)
	// Diana (MG)
}

func ExampleBetween() {
	// Filter users between 25 and 30 years old
	result := filter.Apply(users, filter.Between[User]("Age", 25, 30))

	for _, u := range result {
		fmt.Printf("%s: %d\n", u.Name, u.Age)
	}
	// Output:
	// Ana: 25
	// Bob: 30
	// Diana: 28
}

func ExampleSort() {
	// Sort users by age (ascending)
	sorted := filter.Sort(users, "Age", true)

	for _, u := range sorted {
		fmt.Printf("%s: %d\n", u.Name, u.Age)
	}
	// Output:
	// Ana: 25
	// Diana: 28
	// Bob: 30
	// Carlos: 35
}

func ExampleCustom() {
	// Custom filter: active users with age > 26
	result := filter.Apply(users, filter.Custom(func(u User) bool {
		return u.Active && u.Age > 26
	}))

	for _, u := range result {
		fmt.Printf("%s: %d, active: %v\n", u.Name, u.Age, u.Active)
	}
	// Output:
	// Carlos: 35, active: true
	// Diana: 28, active: true
}
