# Query Parser Design

## Summary

Add a `query/` package that parses HTTP query string parameters into gofilter's existing in-memory filters and executes them against struct slices. This makes gofilter the only Go library that does end-to-end: `query string -> parse -> in-memory filter -> paginated result`.

## Decisions

- **Syntax**: Django-style suffix (`age_gt=25`, `city=SP`, `name_contains=ana`)
- **Security**: Struct tags (`gofilter:"filterable"`) control which fields are exposed
- **Architecture**: New `query/` package that imports `filter/` — keeps the motor clean

## Struct Tags

```go
type User struct {
    Name  string `gofilter:"filterable,sortable,column=name"`
    Age   int    `gofilter:"filterable"`
    City  string `gofilter:"filterable,sortable"`
    Email string // not filterable
}
```

- `filterable` — field can be used in query filters
- `sortable` — field can be used in `sort=` param
- `column=<name>` — custom query param name (default: snake_case of field name)

## API

```go
// Simple: query params -> filtered result
result, err := query.Apply(users, r.URL.Query())

// With options
result, err := query.Apply(users, r.URL.Query(),
    query.WithMaxLimit(100),
    query.WithDefaultSort("Name", true),
)

// Paginated result with metadata
page, err := query.ApplyPaginated(users, r.URL.Query())
// page.Items   []User
// page.Total   int
// page.Page    int
// page.Limit   int
// page.HasNext bool
```

## Operator Mapping

| Query param | Operator | Filter generated |
|---|---|---|
| `name=Ana` | eq (default) | `Eq("Name", "Ana")` |
| `age_gt=25` | gt | `Gt("Age", 25)` |
| `age_gte=18` | gte | `Gte("Age", 18)` |
| `age_lt=30` | lt | `Lt("Age", 30)` |
| `age_lte=30` | lte | `Lte("Age", 30)` |
| `city_ne=SP` | ne | `Ne("City", "SP")` |
| `name_contains=an` | contains | `Contains("Name", "an")` |
| `city_in=SP,RJ,MG` | in | `In("City", ["SP","RJ","MG"])` |
| `age_between=18,30` | between | `Between("Age", 18, 30)` |

## Reserved Parameters

- `sort` — field name, prefix `-` for descending (`sort=-age`)
- `page` — page number (1-based)
- `limit` — items per page

## Type Coercion

The parser inspects the struct field type via reflection and converts the string value:

- `string` — used directly
- `int`, `int8`..`int64` — `strconv.ParseInt`
- `uint`, `uint8`..`uint64` — `strconv.ParseUint`
- `float32`, `float64` — `strconv.ParseFloat`
- `bool` — `strconv.ParseBool`
- `time.Time` — tries RFC3339 then `2006-01-02`

## Error Types

```go
type ErrFieldNotFilterable struct{ Field string }
type ErrFieldNotSortable struct{ Field string }
type ErrInvalidOperator struct{ Param, Operator string }
type ErrInvalidValue struct{ Field, Value, ExpectedType string }
type ErrLimitExceeded struct{ Requested, Max int }
```

All errors implement `error` interface with clear messages suitable for HTTP 400 responses.

## File Structure

```
query/
    query.go        — Apply, ApplyPaginated, options
    parser.go       — parseQueryParams, operator detection, field resolution
    tags.go         — struct tag parsing, field registry
    coerce.go       — type coercion (string -> typed value)
    errors.go       — error types
    query_test.go   — tests
```
