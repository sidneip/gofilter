package query_test

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/sidneip/gofilter/query"
)

type User struct {
	Name   string `json:"name" gofilter:"filterable,sortable"`
	Age    int    `json:"age" gofilter:"filterable,sortable"`
	City   string `json:"city" gofilter:"filterable,sortable"`
	Active bool   `json:"active" gofilter:"filterable"`
}

var users = []User{
	{Name: "Ana", Age: 25, City: "SP", Active: true},
	{Name: "Bob", Age: 30, City: "RJ", Active: false},
	{Name: "Carlos", Age: 35, City: "SP", Active: true},
	{Name: "Diana", Age: 28, City: "MG", Active: true},
}

func ExampleApply() {
	// Simulate query: ?city=SP&sort=-age
	params := url.Values{
		"city": []string{"SP"},
		"sort": []string{"-age"},
	}

	result, err := query.Apply(users, params)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, u := range result {
		fmt.Printf("%s: %d (%s)\n", u.Name, u.Age, u.City)
	}
	// Output:
	// Carlos: 35 (SP)
	// Ana: 25 (SP)
}

func ExampleApplyPaginated() {
	// Simulate query: ?age_gte=25&sort=name&page=1&limit=2
	params := url.Values{
		"age_gte": []string{"25"},
		"sort":    []string{"name"},
		"page":    []string{"1"},
		"limit":   []string{"2"},
	}

	result, err := query.ApplyPaginated(users, params)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Page %d of %d items (has_next: %v)\n", result.Page, result.Total, result.HasNext)
	for _, u := range result.Items {
		fmt.Printf("- %s: %d\n", u.Name, u.Age)
	}
	// Output:
	// Page 1 of 4 items (has_next: true)
	// - Ana: 25
	// - Bob: 30
}

func ExampleApplyPaginated_withOptions() {
	// Empty query with default options
	params := url.Values{}

	result, err := query.ApplyPaginated(users, params,
		query.WithDefaultSort("Name", true),
		query.WithDefaultLimit(10),
		query.WithMaxLimit(100),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, u := range result.Items {
		fmt.Println(u.Name)
	}
	// Output:
	// Ana
	// Bob
	// Carlos
	// Diana
}

func ExampleApplyPaginated_jsonResponse() {
	params := url.Values{
		"city":  []string{"SP"},
		"limit": []string{"2"},
	}

	result, _ := query.ApplyPaginated(users, params)

	jsonBytes, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "items": [
	//     {
	//       "name": "Ana",
	//       "age": 25,
	//       "city": "SP",
	//       "active": true
	//     },
	//     {
	//       "name": "Carlos",
	//       "age": 35,
	//       "city": "SP",
	//       "active": true
	//     }
	//   ],
	//   "total": 2,
	//   "page": 1,
	//   "limit": 2,
	//   "has_next": false
	// }
}
