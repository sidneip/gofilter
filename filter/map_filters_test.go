package filter

import (
	"testing"
)

type Product struct {
	Name       string
	Price      float64
	Attributes map[string]string
	Metadata   map[string]interface{}
	Counts     map[string]int
}

func TestHasKey(t *testing.T) {
	products := []Product{
		{
			Name:  "Laptop",
			Price: 1200.00,
			Attributes: map[string]string{
				"brand":  "TechBrand",
				"color":  "silver",
				"weight": "2kg",
			},
		},
		{
			Name:  "Phone",
			Price: 800.00,
			Attributes: map[string]string{
				"brand": "MobileX",
				"color": "black",
			},
		},
		{
			Name:  "Headphones",
			Price: 150.00,
			Attributes: map[string]string{
				"color":   "white",
				"type":    "wireless",
				"battery": "20h",
			},
		},
	}

	// Test HasKey with existing key
	result := Apply(products, HasKey[Product]("Attributes", "brand"))
	if len(result) != 2 {
		t.Errorf("Expected 2 products with 'brand' attribute, got %d", len(result))
	}

	// Test HasKey with non-existing key
	result = Apply(products, HasKey[Product]("Attributes", "battery"))
	if len(result) != 1 || result[0].Name != "Headphones" {
		t.Errorf("Expected only Headphones to have 'battery' attribute")
	}
}

func TestHasValue(t *testing.T) {
	products := []Product{
		{
			Name:  "Laptop",
			Price: 1200.00,
			Attributes: map[string]string{
				"brand":  "TechBrand",
				"color":  "silver",
				"weight": "2kg",
			},
		},
		{
			Name:  "Phone",
			Price: 800.00,
			Attributes: map[string]string{
				"brand": "MobileX",
				"color": "black",
			},
		},
		{
			Name:  "Headphones",
			Price: 150.00,
			Attributes: map[string]string{
				"color":   "white",
				"type":    "wireless",
				"battery": "20h",
			},
		},
	}

	// Test HasValue with common value
	result := Apply(products, HasValue[Product]("Attributes", "wireless"))
	if len(result) != 1 || result[0].Name != "Headphones" {
		t.Errorf("Expected only Headphones to have 'wireless' value")
	}

	// Test HasValue with common value in multiple products
	result = Apply(products, HasValue[Product]("Attributes", "black"))
	if len(result) != 1 || result[0].Name != "Phone" {
		t.Errorf("Expected only Phone to have 'black' value")
	}
}

func TestKeyValueEquals(t *testing.T) {
	products := []Product{
		{
			Name:  "Laptop",
			Price: 1200.00,
			Attributes: map[string]string{
				"brand":  "TechBrand",
				"color":  "silver",
				"weight": "2kg",
			},
		},
		{
			Name:  "Phone",
			Price: 800.00,
			Attributes: map[string]string{
				"brand": "MobileX",
				"color": "black",
			},
		},
	}

	// Test KeyValueEquals
	result := Apply(products, KeyValueEquals[Product]("Attributes", "brand", "TechBrand"))
	if len(result) != 1 || result[0].Name != "Laptop" {
		t.Errorf("Expected only Laptop to have brand=TechBrand")
	}

	// Test KeyValueEquals with non-matching value
	result = Apply(products, KeyValueEquals[Product]("Attributes", "color", "red"))
	if len(result) != 0 {
		t.Errorf("Expected no products with color=red")
	}
}

func TestMapContainsAll(t *testing.T) {
	products := []Product{
		{
			Name:  "Laptop",
			Price: 1200.00,
			Attributes: map[string]string{
				"brand":  "TechBrand",
				"color":  "silver",
				"weight": "2kg",
			},
		},
		{
			Name:  "Phone",
			Price: 800.00,
			Attributes: map[string]string{
				"brand": "MobileX",
				"color": "black",
			},
		},
	}

	// Test MapContainsAll
	kvPairs := map[interface{}]interface{}{
		"brand": "TechBrand",
		"color": "silver",
	}
	result := Apply(products, MapContainsAll[Product]("Attributes", kvPairs))
	if len(result) != 1 || result[0].Name != "Laptop" {
		t.Errorf("Expected only Laptop to have all the required attributes")
	}

	// Test MapContainsAll with no matches
	kvPairs = map[interface{}]interface{}{
		"brand": "TechBrand",
		"color": "black",
	}
	result = Apply(products, MapContainsAll[Product]("Attributes", kvPairs))
	if len(result) != 0 {
		t.Errorf("Expected no products to have all the required attributes")
	}
}

func TestMapContainsAny(t *testing.T) {
	products := []Product{
		{
			Name:  "Laptop",
			Price: 1200.00,
			Attributes: map[string]string{
				"brand":  "TechBrand",
				"color":  "silver",
				"weight": "2kg",
			},
		},
		{
			Name:  "Phone",
			Price: 800.00,
			Attributes: map[string]string{
				"brand": "MobileX",
				"color": "black",
			},
		},
	}

	// Test MapContainsAny
	kvPairs := map[interface{}]interface{}{
		"brand": "TechBrand",
		"color": "black",
	}
	result := Apply(products, MapContainsAny[Product]("Attributes", kvPairs))
	if len(result) != 2 {
		t.Errorf("Expected both products to have at least one of the attributes")
	}

	// Test MapContainsAny with no matches
	kvPairs = map[interface{}]interface{}{
		"brand": "Unknown",
		"color": "red",
	}
	result = Apply(products, MapContainsAny[Product]("Attributes", kvPairs))
	if len(result) != 0 {
		t.Errorf("Expected no products to have any of the attributes")
	}
}

func TestMapSize(t *testing.T) {
	products := []Product{
		{
			Name:  "Laptop",
			Price: 1200.00,
			Attributes: map[string]string{
				"brand":  "TechBrand",
				"color":  "silver",
				"weight": "2kg",
			},
		},
		{
			Name:  "Phone",
			Price: 800.00,
			Attributes: map[string]string{
				"brand": "MobileX",
				"color": "black",
			},
		},
		{
			Name:  "Keyboard",
			Price: 100.00,
			Attributes: map[string]string{
				"type": "mechanical",
			},
		},
	}

	// Test MapSizeEquals
	result := Apply(products, MapSizeEquals[Product]("Attributes", 2))
	if len(result) != 1 || result[0].Name != "Phone" {
		t.Errorf("Expected only Phone to have exactly 2 attributes")
	}

	// Test MapSizeGreaterThan
	result = Apply(products, MapSizeGreaterThan[Product]("Attributes", 2))
	if len(result) != 1 || result[0].Name != "Laptop" {
		t.Errorf("Expected only Laptop to have more than 2 attributes")
	}

	// Test MapSizeLessThan
	result = Apply(products, MapSizeLessThan[Product]("Attributes", 2))
	if len(result) != 1 || result[0].Name != "Keyboard" {
		t.Errorf("Expected only Keyboard to have less than 2 attributes")
	}
}

func TestComplexMapFilters(t *testing.T) {
	products := []Product{
		{
			Name:  "Laptop",
			Price: 1200.00,
			Attributes: map[string]string{
				"brand":  "TechBrand",
				"color":  "silver",
				"weight": "2kg",
			},
			Metadata: map[string]interface{}{
				"inStock":    true,
				"popularity": 4.5,
				"tags":       []string{"electronics", "computing"},
			},
			Counts: map[string]int{
				"views":     1200,
				"purchases": 45,
			},
		},
		{
			Name:  "Phone",
			Price: 800.00,
			Attributes: map[string]string{
				"brand": "MobileX",
				"color": "black",
			},
			Metadata: map[string]interface{}{
				"inStock":    true,
				"popularity": 4.8,
				"tags":       []string{"electronics", "mobile"},
			},
			Counts: map[string]int{
				"views":     2500,
				"purchases": 120,
			},
		},
	}

	// Debug: Check each product directly
	t.Logf("Before filter - Products count: %d", len(products))
	for i, p := range products {
		inStock, _ := p.Metadata["inStock"]
		t.Logf("Product %d: %s, inStock: %v, views: %d", i, p.Name, inStock, p.Counts["views"])
	}

	// Check inStock filter alone
	inStockFilter := KeyValueEquals[Product]("Metadata", "inStock", true)
	inStockResult := Apply(products, inStockFilter)
	t.Logf("In stock filter alone - Count: %d", len(inStockResult))
	for _, p := range inStockResult {
		t.Logf("  - %s", p.Name)
	}

	// Check views filter alone
	viewsFilter := FilterFunc[Product](func(p Product) bool {
		views, ok := p.Counts["views"]
		return ok && views > 2000
	})
	viewsResult := Apply(products, viewsFilter)
	t.Logf("Views filter alone - Count: %d", len(viewsResult))
	for _, p := range viewsResult {
		t.Logf("  - %s, views: %d", p.Name, p.Counts["views"])
	}

	// Complex filter: Products that are in stock, have more than 2000 views
	result := Apply(products, And[Product](
		KeyValueEquals[Product]("Metadata", "inStock", true),
		FilterFunc[Product](func(p Product) bool {
			views, ok := p.Counts["views"]
			return ok && views > 2000
		}),
	))

	// Debug: Log the results
	t.Logf("After filter - Results count: %d", len(result))
	for _, p := range result {
		t.Logf("  - %s matched the filter", p.Name)
	}

	if len(result) != 1 || result[0].Name != "Phone" {
		t.Errorf("Expected only Phone to match the complex filter criteria")
	}
}
