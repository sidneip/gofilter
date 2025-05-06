package filter

import (
	"fmt"
	"reflect"
	"testing"
)

// We don't need to redefine Product, as it's already defined in map_filters_test.go

func TestDebugKeyValueEquals(t *testing.T) {
	products := []Product{
		{
			Name:  "Laptop",
			Price: 1200.00,
			Metadata: map[string]interface{}{
				"inStock":    true,
				"popularity": 4.5,
				"tags":       []string{"electronics", "computing"},
			},
		},
		{
			Name:  "Phone",
			Price: 800.00,
			Metadata: map[string]interface{}{
				"inStock":    true,
				"popularity": 4.8,
				"tags":       []string{"electronics", "mobile"},
			},
		},
	}

	// Manually check the map field and value
	for i, p := range products {
		fmt.Printf("Product %d: %s\n", i, p.Name)

		// Get the Metadata field
		fieldValue := reflect.ValueOf(p).FieldByName("Metadata")
		fmt.Printf("  Metadata kind: %v\n", fieldValue.Kind())

		// Check for the "inStock" key
		keyValue := reflect.ValueOf("inStock")
		fmt.Printf("  inStock key kind: %v\n", keyValue.Kind())
		fmt.Printf("  Key assignable to map key: %v\n", keyValue.Type().AssignableTo(fieldValue.Type().Key()))

		// Get the map value for the key
		mapValueForKey := fieldValue.MapIndex(keyValue)
		fmt.Printf("  Map value valid: %v\n", mapValueForKey.IsValid())
		if mapValueForKey.IsValid() {
			fmt.Printf("  Map value kind: %v\n", mapValueForKey.Kind())
			fmt.Printf("  Map value: %v\n", mapValueForKey.Interface())

			// Compare with true value
			targetValue := reflect.ValueOf(true)
			fmt.Printf("  Target true value kind: %v\n", targetValue.Kind())

			// Can we convert?
			fmt.Printf("  Target convertible to map value type: %v\n", targetValue.Type().ConvertibleTo(mapValueForKey.Type()))

			// Are they equal by direct comparison?
			fmt.Printf("  Direct comparison: %v == %v: %v\n", mapValueForKey.Interface(), targetValue.Interface(), mapValueForKey.Interface() == targetValue.Interface())
		}
		fmt.Println()
	}

	// Test the filter
	result := Apply(products, KeyValueEquals[Product]("Metadata", "inStock", true))
	fmt.Printf("Filter result count: %d\n", len(result))
	for _, p := range result {
		fmt.Printf("  - %s matched\n", p.Name)
	}
}
