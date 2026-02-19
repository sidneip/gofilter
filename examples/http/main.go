package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sidneip/gofilter/query"
)

type User struct {
	Name   string  `json:"name" gofilter:"filterable,sortable,column=name"`
	Age    int     `json:"age" gofilter:"filterable,sortable"`
	City   string  `json:"city" gofilter:"filterable,sortable"`
	Score  float64 `json:"score" gofilter:"filterable,sortable"`
	Active bool    `json:"active" gofilter:"filterable"`
	Email  string  `json:"email"`
}

var users = []User{
	{Name: "Ana", Age: 20, City: "SP", Score: 8.5, Active: true, Email: "ana@mail.com"},
	{Name: "Bruno", Age: 17, City: "RJ", Score: 7.2, Active: true, Email: "bruno@mail.com"},
	{Name: "Carla", Age: 25, City: "SP", Score: 9.1, Active: false, Email: "carla@mail.com"},
	{Name: "Daniel", Age: 30, City: "MG", Score: 6.8, Active: true, Email: "daniel@mail.com"},
	{Name: "Elena", Age: 22, City: "RJ", Score: 8.9, Active: true, Email: "elena@mail.com"},
	{Name: "Felipe", Age: 28, City: "SP", Score: 7.5, Active: false, Email: "felipe@mail.com"},
	{Name: "Gabi", Age: 19, City: "MG", Score: 9.3, Active: true, Email: "gabi@mail.com"},
	{Name: "Hugo", Age: 35, City: "RJ", Score: 5.4, Active: true, Email: "hugo@mail.com"},
}

func main() {
	http.HandleFunc("GET /users", handleUsers)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println()
	fmt.Println("Try these URLs:")
	fmt.Println("  http://localhost:8080/users")
	fmt.Println("  http://localhost:8080/users?city=SP")
	fmt.Println("  http://localhost:8080/users?age_gt=20&sort=-score")
	fmt.Println("  http://localhost:8080/users?city_in=SP,RJ&active=true&sort=name")
	fmt.Println("  http://localhost:8080/users?score_gte=8.0&sort=-score&page=1&limit=3")
	fmt.Println("  http://localhost:8080/users?age_between=18,25&city_ne=MG")
	fmt.Println("  http://localhost:8080/users?name_contains=a&sort=age")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	page, err := query.ApplyPaginated(users, r.URL.Query(),
		query.WithMaxLimit(50),
		query.WithDefaultLimit(10),
		query.WithDefaultSort("Name", true),
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(page)
}
