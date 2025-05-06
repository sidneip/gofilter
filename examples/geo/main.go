package main

import (
	"fmt"

	"github.com/sidneip/gofilter/filter"
)

// Location represents a place with geographic coordinates
type Location struct {
	Name      string
	Latitude  float64
	Longitude float64
	Type      string
}

func main() {
	// Create some sample locations
	locations := []Location{
		{Name: "New York", Latitude: 40.7128, Longitude: -74.0060, Type: "City"},
		{Name: "Los Angeles", Latitude: 34.0522, Longitude: -118.2437, Type: "City"},
		{Name: "Chicago", Latitude: 41.8781, Longitude: -87.6298, Type: "City"},
		{Name: "Houston", Latitude: 29.7604, Longitude: -95.3698, Type: "City"},
		{Name: "Paris", Latitude: 48.8566, Longitude: 2.3522, Type: "City"},
		{Name: "London", Latitude: 51.5074, Longitude: -0.1278, Type: "City"},
		{Name: "Tokyo", Latitude: 35.6762, Longitude: 139.6503, Type: "City"},
	}

	// Define a center point (San Francisco)
	sanFrancisco := filter.Point{Lat: 37.7749, Lng: -122.4194}

	// Example 1: Find locations within 1000km of San Francisco
	nearSF := filter.Apply(locations,
		filter.WithinRadius[Location]("Latitude", "Longitude", sanFrancisco, 1000))

	fmt.Println("Locations within 1000km of San Francisco:")
	for _, loc := range nearSF {
		fmt.Printf("- %s (%.2f, %.2f)\n", loc.Name, loc.Latitude, loc.Longitude)
	}

	// Example 2: Combining geo filter with regular filters
	citiesNearSF := filter.Apply(locations,
		filter.And[Location](
			filter.WithinRadius[Location]("Latitude", "Longitude", sanFrancisco, 1000),
			filter.Eq[Location]("Type", "City"),
		),
	)

	fmt.Println("\nCities within 1000km of San Francisco:")
	for _, loc := range citiesNearSF {
		fmt.Printf("- %s\n", loc.Name)
	}

	// Example 3: Define a bounding box for the United States (approximate)
	usBox := filter.BoundingBox{
		SouthWest: filter.Point{Lat: 24.396308, Lng: -125.000000},
		NorthEast: filter.Point{Lat: 49.384358, Lng: -66.934570},
	}

	// Find locations within the US
	locationsInUS := filter.Apply(locations,
		filter.WithinBoundingBox[Location]("Latitude", "Longitude", usBox),
	)

	fmt.Println("\nLocations within the US bounding box:")
	for _, loc := range locationsInUS {
		fmt.Printf("- %s\n", loc.Name)
	}

	// Example 4: Sort locations by distance from Tokyo
	tokyo := filter.Point{Lat: 35.6762, Lng: 139.6503}
	sortedLocations := filter.SortByDistance(locations, "Latitude", "Longitude", tokyo)

	fmt.Println("\nLocations sorted by distance from Tokyo:")
	for i, loc := range sortedLocations {
		fmt.Printf("%d. %s\n", i+1, loc.Name)
	}
}
