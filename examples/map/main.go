package main

import (
	"fmt"

	"github.com/sidneip/gofilter/filter"
)

type Product struct {
	Name       string
	Price      float64
	Attributes map[string]string
	Metadata   map[string]interface{}
	Counts     map[string]int
}

func main() {
	// Create some sample products
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
		{
			Name:  "Headphones",
			Price: 150.00,
			Attributes: map[string]string{
				"color":   "white",
				"type":    "wireless",
				"battery": "20h",
			},
			Metadata: map[string]interface{}{
				"inStock":    false,
				"popularity": 3.9,
				"tags":       []string{"electronics", "accessories"},
			},
			Counts: map[string]int{
				"views":     800,
				"purchases": 30,
			},
		},
	}

	// Example 1: Find products that have a "brand" attribute
	brandProducts := filter.Apply(products, filter.HasKey[Product]("Attributes", "brand"))
	fmt.Println("Products with a 'brand' attribute:")
	for _, p := range brandProducts {
		fmt.Printf("- %s (brand: %s)\n", p.Name, p.Attributes["brand"])
	}

	// Example 2: Find products that have "wireless" as a value in their attributes
	wirelessProducts := filter.Apply(products, filter.HasValue[Product]("Attributes", "wireless"))
	fmt.Println("\nProducts with 'wireless' as an attribute value:")
	for _, p := range wirelessProducts {
		fmt.Printf("- %s\n", p.Name)
	}

	// Example 3: Find products where brand equals "TechBrand"
	techBrandProducts := filter.Apply(products, filter.KeyValueEquals[Product]("Attributes", "brand", "TechBrand"))
	fmt.Println("\nProducts with brand = 'TechBrand':")
	for _, p := range techBrandProducts {
		fmt.Printf("- %s\n", p.Name)
	}

	// Example 4: Find products that have both specified attributes
	requiredAttrs := map[interface{}]interface{}{
		"color": "silver",
		"brand": "TechBrand",
	}
	matchingProducts := filter.Apply(products, filter.MapContainsAll[Product]("Attributes", requiredAttrs))
	fmt.Println("\nProducts with all required attributes:")
	for _, p := range matchingProducts {
		fmt.Printf("- %s\n", p.Name)
	}

	// Example 5: Find products that have at least one of the specified attributes
	someAttrs := map[interface{}]interface{}{
		"type":  "wireless",
		"brand": "MobileX",
	}
	someMatchingProducts := filter.Apply(products, filter.MapContainsAny[Product]("Attributes", someAttrs))
	fmt.Println("\nProducts with at least one of the specified attributes:")
	for _, p := range someMatchingProducts {
		fmt.Printf("- %s\n", p.Name)
	}

	// Example 6: Find products with exactly 2 attributes
	twoAttrProducts := filter.Apply(products, filter.MapSizeEquals[Product]("Attributes", 2))
	fmt.Println("\nProducts with exactly 2 attributes:")
	for _, p := range twoAttrProducts {
		fmt.Printf("- %s (attributes: %v)\n", p.Name, p.Attributes)
	}

	// Example 7: Complex filter combining map filters and other filters
	complexFilter := filter.And[Product](
		filter.KeyValueEquals[Product]("Metadata", "inStock", true),
		filter.MapSizeGreaterThan[Product]("Attributes", 1),
		filter.Gt[Product]("Price", 500),
	)

	complexResults := filter.Apply(products, complexFilter)
	fmt.Println("\nProducts matching complex criteria (in stock, >1 attribute, price > $500):")
	for _, p := range complexResults {
		fmt.Printf("- %s ($%.2f)\n", p.Name, p.Price)
	}
}
