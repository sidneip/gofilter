package main

import (
	"fmt"

	"github.com/sidneip/gofilter/filter"
)

type Address struct {
	City    string
	Country string
}

type Person struct {
	Name    string
	Age     int
	Address Address
	Hobbies []string
}

func main() {
	people := []Person{
		{
			Name: "Ana",
			Age:  20,
			Address: Address{
				City:    "New York",
				Country: "USA",
			},
			Hobbies: []string{"reading", "swimming"},
		},
		{
			Name: "Bruno",
			Age:  17,
			Address: Address{
				City:    "London",
				Country: "UK",
			},
			Hobbies: []string{"soccer"},
		},
		{
			Name: "Carla",
			Age:  25,
			Address: Address{
				City:    "Paris",
				Country: "France",
			},
			Hobbies: []string{"cinema", "reading"},
		},
		{
			Name: "David",
			Age:  30,
			Address: Address{
				City:    "San Francisco",
				Country: "USA",
			},
			Hobbies: []string{"hiking", "photography"},
		},
	}

	// Example 1: Filter people older than 18
	adults := filter.Apply(people, filter.Gt[Person]("Age", 18))
	fmt.Println("Adults:")
	for _, person := range adults {
		fmt.Printf("- %s (%d)\n", person.Name, person.Age)
	}

	// Example 2: Filter people with names containing "a"
	withA := filter.Apply(people, filter.Contains[Person]("Name", "a"))
	fmt.Println("\nNames containing 'a':")
	for _, person := range withA {
		fmt.Printf("- %s\n", person.Name)
	}

	// Example 3: Composed filter (adults from USA)
	adultsFromUSA := filter.Apply(people,
		filter.And[Person](
			filter.Gt[Person]("Age", 18),
			filter.Eq[Person]("Address.Country", "USA"),
		),
	)
	fmt.Println("\nAdults from USA:")
	for _, person := range adultsFromUSA {
		fmt.Printf("- %s (%d) from %s\n", person.Name, person.Age, person.Address.City)
	}

	// Example 4: People who like reading or are teenagers
	readingOrTeens := filter.Apply(people,
		filter.Or[Person](
			filter.Contains[Person]("Hobbies", "reading"),
			filter.Lt[Person]("Age", 20),
		),
	)
	fmt.Println("\nPeople who like reading or are teenagers:")
	for _, person := range readingOrTeens {
		fmt.Printf("- %s (%d, hobbies: %v)\n", person.Name, person.Age, person.Hobbies)
	}
}
