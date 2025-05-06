package filter

import (
	"testing"
)

type Location struct {
	Name      string
	Latitude  float64
	Longitude float64
}

func TestWithinRadius(t *testing.T) {
	locations := []Location{
		{Name: "New York", Latitude: 40.7128, Longitude: -74.0060},
		{Name: "Los Angeles", Latitude: 34.0522, Longitude: -118.2437},
		{Name: "Tokyo", Latitude: 35.6762, Longitude: 139.6503},
	}

	// San Francisco
	sf := Point{Lat: 37.7749, Lng: -122.4194}

	// Test locations within 1000km of San Francisco
	result := Apply(locations, WithinRadius[Location]("Latitude", "Longitude", sf, 1000))
	if len(result) != 1 {
		t.Errorf("Expected 1 location within 1000km of SF, got %d", len(result))
	}

	if len(result) > 0 && result[0].Name != "Los Angeles" {
		t.Errorf("Expected Los Angeles to be within 1000km of SF, got %s", result[0].Name)
	}

	// Test with larger radius
	result = Apply(locations, WithinRadius[Location]("Latitude", "Longitude", sf, 5000))
	if len(result) != 2 {
		t.Errorf("Expected 2 locations within 5000km of SF, got %d", len(result))
	}
}

func TestOutsideRadius(t *testing.T) {
	locations := []Location{
		{Name: "New York", Latitude: 40.7128, Longitude: -74.0060},
		{Name: "Los Angeles", Latitude: 34.0522, Longitude: -118.2437},
		{Name: "Tokyo", Latitude: 35.6762, Longitude: 139.6503},
	}

	// San Francisco
	sf := Point{Lat: 37.7749, Lng: -122.4194}

	// Test locations outside 1000km of San Francisco
	result := Apply(locations, OutsideRadius[Location]("Latitude", "Longitude", sf, 1000))
	if len(result) != 2 {
		t.Errorf("Expected 2 locations outside 1000km of SF, got %d", len(result))
	}

	for _, loc := range result {
		if loc.Name == "Los Angeles" {
			t.Errorf("Expected Los Angeles to be within 1000km of SF, not outside")
		}
	}
}

func TestWithinBoundingBox(t *testing.T) {
	locations := []Location{
		{Name: "New York", Latitude: 40.7128, Longitude: -74.0060},
		{Name: "Los Angeles", Latitude: 34.0522, Longitude: -118.2437},
		{Name: "Tokyo", Latitude: 35.6762, Longitude: 139.6503},
	}

	// Rough US bounding box
	usBox := BoundingBox{
		SouthWest: Point{Lat: 24.396308, Lng: -125.000000},
		NorthEast: Point{Lat: 49.384358, Lng: -66.934570},
	}

	// Test locations within US bounding box
	result := Apply(locations, WithinBoundingBox[Location]("Latitude", "Longitude", usBox))
	if len(result) != 2 {
		t.Errorf("Expected 2 locations within US bounding box, got %d", len(result))
	}

	// Ensure Tokyo is not in the results
	for _, loc := range result {
		if loc.Name == "Tokyo" {
			t.Errorf("Expected Tokyo to be outside US bounding box")
		}
	}
}

func TestSortByDistance(t *testing.T) {
	locations := []Location{
		{Name: "New York", Latitude: 40.7128, Longitude: -74.0060},
		{Name: "Los Angeles", Latitude: 34.0522, Longitude: -118.2437},
		{Name: "Tokyo", Latitude: 35.6762, Longitude: 139.6503},
	}

	// Tokyo
	tokyo := Point{Lat: 35.6762, Lng: 139.6503}

	// Sort by distance from Tokyo
	sorted := SortByDistance(locations, "Latitude", "Longitude", tokyo)
	if len(sorted) != 3 {
		t.Errorf("Expected 3 sorted locations, got %d", len(sorted))
	}

	if sorted[0].Name != "Tokyo" {
		t.Errorf("Expected Tokyo to be closest to itself, got %s", sorted[0].Name)
	}
}
