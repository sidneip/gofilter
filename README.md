# gofilter

[![Go](https://github.com/sidneip/gofilter/actions/workflows/ci.yml/badge.svg)](https://github.com/sidneip/gofilter/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sidneip/gofilter.svg)](https://pkg.go.dev/github.com/sidneip/gofilter)
[![Go Report Card](https://goreportcard.com/badge/github.com/sidneip/gofilter)](https://goreportcard.com/report/github.com/sidneip/gofilter)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Zero Dependencies](https://img.shields.io/badge/dependencies-zero-brightgreen)

**The missing query engine for Go structs.** Filter, sort, and paginate slices of structs using HTTP query parameters — no database required.

```go
// Before gofilter: manual parsing, if/else chains, boilerplate for every endpoint
city := r.URL.Query().Get("city")
minAge := r.URL.Query().Get("age_gt")
sortBy := r.URL.Query().Get("sort")
// ... 50 lines of manual filtering logic per endpoint

// After gofilter: one line, any endpoint
result, err := query.ApplyPaginated(users, r.URL.Query())
```

```
┌──────────────────┐     ┌───────────┐     ┌──────────────┐     ┌──────────────┐
│  ?city=SP        │     │  Parse &  │     │  In-memory   │     │  Paginated   │
│  &age_gt=25      │────▶│  Validate │────▶│  Filter      │────▶│  JSON        │
│  &sort=-name     │     │  & Coerce │     │  & Sort      │     │  Response    │
│  &page=1&limit=10│     └───────────┘     └──────────────┘     └──────────────┘
└──────────────────┘
```

## Why gofilter?

Every Go library that parses query parameters generates SQL. Every library that filters in-memory requires manual closures. **Nothing connects the two.**

| Library | Parses query params? | Filters in-memory? | Zero deps? |
|---|---|---|---|
| [samber/lo](https://github.com/samber/lo) | No | Yes (manual closures) | Yes |
| [a8m/rql](https://github.com/a8m/rql) | Yes (JSON body) | No (generates SQL) | No |
| [go-goyave/filter](https://github.com/go-goyave/filter) | Yes | No (requires GORM) | No |
| [cbrand/go-filterparams](https://github.com/cbrand/go-filterparams) | Yes | No (parse only) | Yes |
| **gofilter** | **Yes** | **Yes** | **Yes** |

## When to use gofilter

**Use it when you have data in memory and need to expose filtering via API:**

- REST APIs serving cached or preloaded data
- Microservices with in-memory stores (config, feature flags, catalogs)
- Prototyping APIs without setting up a database
- Admin dashboards with client-side filterable tables
- Static datasets (countries, currencies, product catalogs)
- Any `[]struct` that needs `?field=value` filtering

**Don't use it when:**

- You're querying a database directly (use your ORM's filtering)
- You have millions of records (use a proper database or search engine)

## Installation

```bash
go get github.com/sidneip/gofilter
```

Requires Go 1.22+. **Zero external dependencies** — only the Go standard library.

## Quick Start

### HTTP API in 15 lines

```go
type User struct {
    Name  string  `json:"name" gofilter:"filterable,sortable"`
    Age   int     `json:"age" gofilter:"filterable,sortable"`
    City  string  `json:"city" gofilter:"filterable,sortable"`
    Email string  `json:"email"` // not exposed — safe by default
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    page, err := query.ApplyPaginated(users, r.URL.Query(),
        query.WithMaxLimit(100),
        query.WithDefaultSort("Name", true),
    )
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
    json.NewEncoder(w).Encode(page)
}
```

Your API now supports all of these — with zero additional code:

```
GET /users                                     → all users, sorted by name
GET /users?city=SP                             → users from São Paulo
GET /users?age_gt=25&sort=-name                → age > 25, sorted by name desc
GET /users?city_in=SP,RJ,MG&age_between=18,30 → multiple filters combined
GET /users?name_contains=ana&page=2&limit=10   → search + pagination
GET /users?email=test                          → 400: field "email" is not filterable
```

### Response format

```json
{
  "items": [
    {"name": "Ana", "age": 20, "city": "SP"},
    {"name": "Carla", "age": 25, "city": "SP"}
  ],
  "total": 2,
  "page": 1,
  "limit": 10,
  "has_next": false
}
```

Works with **any Go HTTP router**: net/http, Chi, Gin, Echo, Fiber — gofilter just needs `url.Values`.

### Framework Examples

<details>
<summary><strong>Gin</strong></summary>

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/sidneip/gofilter/query"
)

func main() {
    r := gin.Default()

    r.GET("/users", func(c *gin.Context) {
        page, err := query.ApplyPaginated(users, c.Request.URL.Query(),
            query.WithMaxLimit(100),
        )
        if err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, page)
    })

    r.Run(":8080")
}
```

</details>

<details>
<summary><strong>Echo</strong></summary>

```go
import (
    "github.com/labstack/echo/v4"
    "github.com/sidneip/gofilter/query"
)

func main() {
    e := echo.New()

    e.GET("/users", func(c echo.Context) error {
        page, err := query.ApplyPaginated(users, c.QueryParams(),
            query.WithMaxLimit(100),
        )
        if err != nil {
            return c.JSON(400, map[string]string{"error": err.Error()})
        }
        return c.JSON(200, page)
    })

    e.Start(":8080")
}
```

</details>

<details>
<summary><strong>Chi</strong></summary>

```go
import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/sidneip/gofilter/query"
)

func main() {
    r := chi.NewRouter()

    r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
        page, err := query.ApplyPaginated(users, r.URL.Query(),
            query.WithMaxLimit(100),
        )
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
            return
        }
        json.NewEncoder(w).Encode(page)
    })

    http.ListenAndServe(":8080", r)
}
```

</details>

<details>
<summary><strong>Fiber</strong></summary>

```go
import (
    "github.com/gofiber/fiber/v2"
    "github.com/sidneip/gofilter/query"
)

func main() {
    app := fiber.New()

    app.Get("/users", func(c *fiber.Ctx) error {
        // Convert Fiber's query params to url.Values
        params := make(map[string][]string)
        c.Context().QueryArgs().VisitAll(func(key, value []byte) {
            params[string(key)] = append(params[string(key)], string(value))
        })

        page, err := query.ApplyPaginated(users, params,
            query.WithMaxLimit(100),
        )
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(page)
    })

    app.Listen(":8080")
}
```

</details>

## Query Syntax

Django-style suffixes on field names — intuitive for anyone who's used Django REST Framework, Rails, or Strapi:

| Query param | Operator | Example |
|---|---|---|
| `field` | equals | `?city=SP` |
| `field_gt` | greater than | `?age_gt=25` |
| `field_gte` | greater or equal | `?age_gte=18` |
| `field_lt` | less than | `?age_lt=30` |
| `field_lte` | less or equal | `?age_lte=30` |
| `field_ne` | not equal | `?city_ne=SP` |
| `field_contains` | substring match | `?name_contains=ana` |
| `field_in` | in list (comma-separated) | `?city_in=SP,RJ,MG` |
| `field_between` | range inclusive (comma-separated) | `?age_between=18,30` |

Reserved parameters:

| Param | Description | Example |
|---|---|---|
| `sort` | Sort field, `-` prefix for descending | `?sort=-age` |
| `page` | Page number (1-based) | `?page=2` |
| `limit` | Items per page | `?limit=10` |

Multiple filters are combined with AND logic.

## Struct Tags

Control which fields are exposed to filtering — **secure by default**:

```go
type Product struct {
    Name     string  `gofilter:"filterable,sortable"`       // filter + sort
    Price    float64 `gofilter:"filterable,sortable"`       // filter + sort
    Category string  `gofilter:"filterable,column=cat"`     // custom param: ?cat=electronics
    SKU      string                                         // not exposed
    Secret   string                                         // not exposed
}
```

| Tag | Description |
|---|---|
| `filterable` | Field can be used in query filters |
| `sortable` | Field can be used with `sort=` |
| `column=<name>` | Custom query parameter name (default: snake_case of field) |

Fields without the `gofilter` tag are **never** exposed — you can't accidentally leak sensitive data.

## Options

```go
query.Apply(items, params,
    query.WithMaxLimit(100),              // reject requests with limit > 100
    query.WithDefaultLimit(20),           // default items per page
    query.WithDefaultSort("Name", true),  // fallback sort when none specified
)
```

## Type Coercion

Values from query strings are **automatically converted** based on the struct field type:

| Struct field type | Query value | Parsed as |
|---|---|---|
| `string` | `?name=Ana` | `"Ana"` |
| `int`, `int64` | `?age=25` | `25` |
| `float64` | `?score=8.5` | `8.5` |
| `bool` | `?active=true` | `true` |
| `time.Time` | `?date=2024-01-15` | `time.Time` |
| `time.Time` | `?date=2024-01-15T10:30:00Z` | `time.Time` (RFC3339) |

Invalid values return typed errors (no panics, no silent failures).

## Error Handling

Typed errors designed for clean HTTP 400 responses:

```go
page, err := query.ApplyPaginated(users, r.URL.Query())
if err != nil {
    switch err.(type) {
    case *query.ErrFieldNotFilterable:  // field "email" is not filterable
    case *query.ErrFieldNotSortable:    // field "email" is not sortable
    case *query.ErrInvalidValue:        // invalid value "abc" for field "Age": expected int
    case *query.ErrLimitExceeded:       // requested limit 500 exceeds maximum 100
    }
}
```

## Performance

Benchmarks on Apple M4, Go 1.23 (see `query/bench_test.go`):

| Scenario | Time | Allocs |
|---|---|---|
| 1K items, single filter | ~80μs | 2K |
| 1K items, 3 filters | ~125μs | 3.4K |
| 10K items, filter + sort + pagination | ~7.5ms | 252K |
| 100K items, single filter | ~8.5ms | 200K |

gofilter is designed for collections up to ~100K items. For larger datasets, use a database.

## Programmatic API

For building filters in code without HTTP (the `filter/` package):

```go
import "github.com/sidneip/gofilter/filter"

// Compose filters
f := filter.And[User](
    filter.Gt[User]("Age", 18),
    filter.Contains[User]("Name", "ana"),
    filter.In[User]("City", []interface{}{"SP", "RJ"}),
)

result := filter.Apply(users, f)
sorted := filter.Sort(result, "Age", true)
```

<details>
<summary><strong>Advanced filters</strong></summary>

```go
// Case-insensitive string matching
filter.StringMatch[User]("Name", "ana", filter.StringMatchOptions{
    Mode: filter.ContainsMatch, IgnoreCase: true,
})

// Regex
filter.RegexMatch[User]("Email", `^[a-z]+@gmail\.com$`)

// Date ranges
filter.DateBetween[User]("CreatedAt", startDate, endDate)

// Nil/zero checks
filter.IsNil[User]("DeletedAt")
filter.IsNotZero[User]("Score")

// Custom function
filter.Custom[User](func(u User) bool {
    return u.Age > 18 && strings.HasPrefix(u.Email, "admin")
})
```

</details>

<details>
<summary><strong>Geospatial filters</strong></summary>

```go
center := filter.Point{Lat: -23.5505, Lng: -46.6333} // São Paulo

// Within radius (Haversine distance)
nearby := filter.Apply(places,
    filter.WithinRadius[Place]("Lat", "Lng", center, 50.0)) // 50km

// Bounding box
box := filter.BoundingBox{
    SouthWest: filter.Point{Lat: -24.0, Lng: -47.0},
    NorthEast: filter.Point{Lat: -23.0, Lng: -46.0},
}
inBox := filter.Apply(places,
    filter.WithinBoundingBox[Place]("Lat", "Lng", box))

// Sort by distance
sorted := filter.SortByDistance(places, "Lat", "Lng", center)
```

</details>

<details>
<summary><strong>Map field filters</strong></summary>

```go
filter.HasKey[Product]("Attrs", "color")
filter.KeyValueEquals[Product]("Attrs", "brand", "Nike")
filter.MapContainsAll[Product]("Attrs", map[interface{}]interface{}{
    "color": "red", "size": "M",
})
```

</details>

## Examples

See the [examples/](examples/) directory:

| Example | Description |
|---|---|
| [examples/http](examples/http/) | Full HTTP server with query filtering |
| [examples/query](examples/query/) | Query parser usage without HTTP |
| [examples/simple](examples/simple/) | Basic programmatic filtering |
| [examples/geo](examples/geo/) | Geospatial filtering |
| [examples/map](examples/map/) | Map field filtering |

## Roadmap

gofilter is actively maintained. Here's what's coming next:

- [ ] **Framework middleware** — Drop-in middleware for Gin, Echo, Chi, and Fiber
- [ ] **Nested struct queries** — Filter by nested fields: `?address.city=SP`
- [ ] **OR logic via query params** — Support `?or=city:SP,city:RJ` syntax
- [ ] **Full-text search operator** — `?name_search=ana` with fuzzy matching
- [ ] **OpenAPI schema generation** — Auto-generate filter documentation from struct tags
- [ ] **Cached field registry** — Pre-compute struct metadata for zero-alloc parsing
- [ ] **Benchmarks suite** — Comparative benchmarks against manual filtering

Have an idea? [Open an issue](https://github.com/sidneip/gofilter/issues) — we'd love to hear it.

## Contributing

Contributions are welcome! gofilter is designed to be easy to contribute to:

```
gofilter/
├── filter/    # Core filter engine (operators, composition, geo, maps)
├── query/     # Query string parser (parsing, coercion, pagination)
└── examples/  # Usage examples
```

**Good first issues:**

- Add a new operator (e.g., `starts_with`, `ends_with`)
- Improve test coverage for `filter/` package (currently at 42%)
- Add middleware for your favorite Go framework
- Write benchmarks comparing gofilter vs manual filtering

**How to contribute:**

1. Fork the repo
2. Create your branch (`git checkout -b feat/my-feature`)
3. Write tests first, then implementation
4. Run `go test ./... -race` to make sure everything passes
5. Open a PR

## License

MIT
