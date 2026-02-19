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
	fmt.Println("Query: ?city=SP&age_gt=18&sort=-age")
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
	fmt.Println("Query: ?age_gte=18&sort=name&page=1&limit=2")
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
