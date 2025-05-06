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
- Geospatial filtering for location-based data
- Map field filtering for key-value data structures

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

### Map Operators

- `HasKey(field, key)` - Map field contains the specified key
- `HasValue(field, value)` - Map field contains the specified value
- `KeyValueEquals(field, key, value)` - Key in map field has the specified value
- `MapContainsAll(field, kvPairs)` - Map field contains all the specified key-value pairs
- `MapContainsAny(field, kvPairs)` - Map field contains at least one of the specified key-value pairs
- `MapSizeEquals(field, size)` - Map field has exactly the specified number of entries
- `MapSizeGreaterThan(field, size)` - Map field has more than the specified number of entries
- `MapSizeLessThan(field, size)` - Map field has fewer than the specified number of entries

#### Map Filter Example

```go
type Product struct {
    Name       string
    Attributes map[string]string
}

products := []Product{
    {
        Name: "Laptop", 
        Attributes: map[string]string{
            "brand": "TechBrand",
            "color": "silver",
        },
    },
    {
        Name: "Phone", 
        Attributes: map[string]string{
            "brand": "MobileX", 
            "color": "black",
        },
    },
}

// Find products with a specific brand
techBrandProducts := filter.Apply(products, 
    filter.KeyValueEquals[Product]("Attributes", "brand", "TechBrand"))

// Find products that have all required attributes
requiredAttrs := map[interface{}]interface{}{
    "color": "silver",
    "brand": "TechBrand",
}
matchingProducts := filter.Apply(products, 
    filter.MapContainsAll[Product]("Attributes", requiredAttrs))
```

### Geospatial Operators

- `WithinRadius(latField, lngField, centerPoint, radiusKm)` - Checks if a location is within a radius from a center point
- `OutsideRadius(latField, lngField, centerPoint, radiusKm)` - Checks if a location is outside a radius from a center point
- `WithinBoundingBox(latField, lngField, box)` - Checks if a location is within a geographic rectangle
- `SortByDistance(items, latField, lngField, centerPoint)` - Sorts items by distance from a center point

#### Geospatial Example

```go
// Define location data
locations := []Location{
    {Name: "New York", Latitude: 40.7128, Longitude: -74.0060},
    {Name: "Los Angeles", Latitude: 34.0522, Longitude: -118.2437},
    {Name: "Chicago", Latitude: 41.8781, Longitude: -87.6298},
}

// Define a center point (San Francisco)
sanFrancisco := filter.Point{Lat: 37.7749, Lng: -122.4194}

// Find locations within 1000km of San Francisco
nearSF := filter.Apply(locations, 
    filter.WithinRadius[Location]("Latitude", "Longitude", sanFrancisco, 1000))

// Define a bounding box for the United States (approximate)
usBox := filter.BoundingBox{
    SouthWest: filter.Point{Lat: 24.396308, Lng: -125.000000},
    NorthEast: filter.Point{Lat: 49.384358, Lng: -66.934570},
}

// Find locations within the US
locationsInUS := filter.Apply(locations,
    filter.WithinBoundingBox[Location]("Latitude", "Longitude", usBox))

// Sort locations by distance from Tokyo
tokyo := filter.Point{Lat: 35.6762, Lng: 139.6503}
sortedLocations := filter.SortByDistance(locations, "Latitude", "Longitude", tokyo)
```
