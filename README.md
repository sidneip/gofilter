# gofilter

A generic library for dynamically and flexibly filtering slices of structs in Go.

## Motivation

Go does not provide a native solution for filtering slices of structs by dynamic fields, especially when you need to apply multiple conditions or access nested fields. `gofilter` fills this gap by offering a simple and powerful API to create reusable and composable filters.

## Installation

```bash
go get github.com/sidneip/gofilter
```

Make sure you're using Go 1.18+ for generics support.

## Main Features

- Generic filtering by any struct field (including nested fields)
- Support for operators: equals, not equals, greater than, less than, contains, etc.
- Filter composition with AND, OR, NOT
- Easy integration with any struct

## Usage Example

```go
package main

import (
    "fmt"
    "github.com/sidneip/gofilter/filter"
)

type Person struct {
    Name    string
    Age     int
    Hobbies []string
}

func main() {
    people := []Person{
        {Name: "Ana", Age: 20, Hobbies: []string{"reading", "swimming"}},
        {Name: "Bruno", Age: 17, Hobbies: []string{"soccer"}},
        {Name: "Carla", Age: 25, Hobbies: []string{"cinema", "reading"}},
    }

    result := filter.Apply(people,
        filter.And[Person](
            filter.Gt[Person]("Age", 18),
            filter.Contains[Person]("Name", "a"),
        ),
    )

    fmt.Println(result)
}
```

## Available Operators

### Comparison Operators

- `Eq(field, value)` - Equal to
- `Ne(field, value)` - Not equal to
- `Gt(field, value)` - Greater than
- `Lt(field, value)` - Less than
- `Gte(field, value)` - Greater than or equal to
- `Lte(field, value)` - Less than or equal to
- `Contains(field, value)` - Field contains value (for strings, slices, arrays)
- `In(field, []value)` - Field is in a list of values

### Logical Operators

- `And(filter1, filter2, ...)` - All filters must match
- `Or(filter1, filter2, ...)` - At least one filter must match
- `Not(filter)` - Negates the result of a filter

### Special Operators

- `IsNil(field)` - Field is nil (for pointers, slices, maps)
- `IsNotNil(field)` - Field is not nil
- `IsZero(field)` - Field has its zero value
- `IsNotZero(field)` - Field does not have its zero value
