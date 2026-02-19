package query

import (
	"net/url"
	"testing"
)

type BenchUser struct {
	Name   string  `gofilter:"filterable,sortable"`
	Age    int     `gofilter:"filterable,sortable"`
	City   string  `gofilter:"filterable,sortable"`
	Score  float64 `gofilter:"filterable,sortable"`
	Active bool    `gofilter:"filterable"`
}

func generateUsers(n int) []BenchUser {
	cities := []string{"SP", "RJ", "MG", "BA", "RS"}
	users := make([]BenchUser, n)
	for i := range users {
		users[i] = BenchUser{
			Name:   "User" + string(rune('A'+i%26)),
			Age:    18 + i%50,
			City:   cities[i%len(cities)],
			Score:  float64(i%100) / 10.0,
			Active: i%3 != 0,
		}
	}
	return users
}

func BenchmarkApply_1K_SingleFilter(b *testing.B) {
	users := generateUsers(1_000)
	params := url.Values{"city": {"SP"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(users, params)
	}
}

func BenchmarkApply_1K_MultipleFilters(b *testing.B) {
	users := generateUsers(1_000)
	params := url.Values{"city": {"SP"}, "age_gt": {"25"}, "active": {"true"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(users, params)
	}
}

func BenchmarkApply_10K_WithSortAndPagination(b *testing.B) {
	users := generateUsers(10_000)
	params := url.Values{"city": {"SP"}, "age_gt": {"25"}, "sort": {"-score"}, "page": {"1"}, "limit": {"20"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ApplyPaginated(users, params)
	}
}

func BenchmarkApply_100K_SingleFilter(b *testing.B) {
	users := generateUsers(100_000)
	params := url.Values{"city": {"SP"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(users, params)
	}
}
