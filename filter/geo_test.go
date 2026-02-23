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

// LocationWithIntCoords tests integer coordinate types
type LocationWithIntCoords struct {
	Name string
	Lat  int
	Lng  int
}

func TestWithinRadiusIntCoords(t *testing.T) {
	locations := []LocationWithIntCoords{
		{Name: "Origin", Lat: 0, Lng: 0},
		{Name: "Near", Lat: 1, Lng: 1},
		{Name: "Far", Lat: 45, Lng: 90},
	}

	origin := Point{Lat: 0, Lng: 0}
	result := Apply(locations, WithinRadius[LocationWithIntCoords]("Lat", "Lng", origin, 200))
	if len(result) != 2 {
		t.Errorf("Expected 2 locations within 200km of origin, got %d", len(result))
	}
}

func TestWithinBoundingBoxIntCoords(t *testing.T) {
	locations := []LocationWithIntCoords{
		{Name: "Inside", Lat: 10, Lng: 10},
		{Name: "Outside", Lat: 50, Lng: 50},
	}

	box := BoundingBox{
		SouthWest: Point{Lat: 0, Lng: 0},
		NorthEast: Point{Lat: 20, Lng: 20},
	}

	result := Apply(locations, WithinBoundingBox[LocationWithIntCoords]("Lat", "Lng", box))
	if len(result) != 1 {
		t.Errorf("Expected 1 location inside box, got %d", len(result))
	}
	if len(result) > 0 && result[0].Name != "Inside" {
		t.Errorf("Expected 'Inside' location, got %s", result[0].Name)
	}
}

func TestSortByDistanceIntCoords(t *testing.T) {
	locations := []LocationWithIntCoords{
		{Name: "Far", Lat: 45, Lng: 90},
		{Name: "Near", Lat: 1, Lng: 1},
		{Name: "Origin", Lat: 0, Lng: 0},
	}

	origin := Point{Lat: 0, Lng: 0}
	sorted := SortByDistance(locations, "Lat", "Lng", origin)
	if len(sorted) != 3 {
		t.Errorf("Expected 3 sorted locations, got %d", len(sorted))
	}
	if sorted[0].Name != "Origin" {
		t.Errorf("Expected Origin first, got %s", sorted[0].Name)
	}
	if sorted[1].Name != "Near" {
		t.Errorf("Expected Near second, got %s", sorted[1].Name)
	}
}

// LocationWithUnsupportedType tests unsupported coordinate types
type LocationWithUnsupportedType struct {
	Name string
	Lat  string
	Lng  string
}

func TestWithinRadiusUnsupportedType(t *testing.T) {
	locations := []LocationWithUnsupportedType{
		{Name: "Test", Lat: "40.7128", Lng: "-74.0060"},
	}

	center := Point{Lat: 40.7128, Lng: -74.0060}
	result := Apply(locations, WithinRadius[LocationWithUnsupportedType]("Lat", "Lng", center, 100))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for unsupported type, got %d", len(result))
	}
}

func TestWithinBoundingBoxUnsupportedType(t *testing.T) {
	locations := []LocationWithUnsupportedType{
		{Name: "Test", Lat: "10", Lng: "10"},
	}

	box := BoundingBox{
		SouthWest: Point{Lat: 0, Lng: 0},
		NorthEast: Point{Lat: 20, Lng: 20},
	}

	result := Apply(locations, WithinBoundingBox[LocationWithUnsupportedType]("Lat", "Lng", box))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for unsupported type, got %d", len(result))
	}
}

func TestSortByDistanceUnsupportedType(t *testing.T) {
	locations := []LocationWithUnsupportedType{
		{Name: "Test", Lat: "40.7128", Lng: "-74.0060"},
	}

	center := Point{Lat: 40.7128, Lng: -74.0060}
	sorted := SortByDistance(locations, "Lat", "Lng", center)
	if len(sorted) != 0 {
		t.Errorf("Expected 0 results for unsupported type, got %d", len(sorted))
	}
}

// Test invalid field names
func TestWithinRadiusInvalidField(t *testing.T) {
	locations := []Location{
		{Name: "Test", Latitude: 40.7128, Longitude: -74.0060},
	}

	center := Point{Lat: 40.7128, Lng: -74.0060}

	// Invalid lat field
	result := Apply(locations, WithinRadius[Location]("InvalidLat", "Longitude", center, 100))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for invalid lat field, got %d", len(result))
	}

	// Invalid lng field
	result = Apply(locations, WithinRadius[Location]("Latitude", "InvalidLng", center, 100))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for invalid lng field, got %d", len(result))
	}
}

func TestWithinBoundingBoxInvalidField(t *testing.T) {
	locations := []Location{
		{Name: "Test", Latitude: 10, Longitude: 10},
	}

	box := BoundingBox{
		SouthWest: Point{Lat: 0, Lng: 0},
		NorthEast: Point{Lat: 20, Lng: 20},
	}

	// Invalid lat field
	result := Apply(locations, WithinBoundingBox[Location]("InvalidLat", "Longitude", box))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for invalid lat field, got %d", len(result))
	}

	// Invalid lng field
	result = Apply(locations, WithinBoundingBox[Location]("Latitude", "InvalidLng", box))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for invalid lng field, got %d", len(result))
	}
}

func TestSortByDistanceInvalidField(t *testing.T) {
	locations := []Location{
		{Name: "Test", Latitude: 40.7128, Longitude: -74.0060},
	}

	center := Point{Lat: 40.7128, Lng: -74.0060}

	// Invalid lat field
	sorted := SortByDistance(locations, "InvalidLat", "Longitude", center)
	if len(sorted) != 0 {
		t.Errorf("Expected 0 results for invalid lat field, got %d", len(sorted))
	}

	// Invalid lng field
	sorted = SortByDistance(locations, "Latitude", "InvalidLng", center)
	if len(sorted) != 0 {
		t.Errorf("Expected 0 results for invalid lng field, got %d", len(sorted))
	}
}

// LocationWithMixedTypes tests unsupported lng type with valid lat
type LocationWithMixedTypes struct {
	Name string
	Lat  float64
	Lng  bool
}

func TestWithinRadiusMixedTypes(t *testing.T) {
	locations := []LocationWithMixedTypes{
		{Name: "Test", Lat: 40.7128, Lng: true},
	}

	center := Point{Lat: 40.7128, Lng: -74.0060}
	result := Apply(locations, WithinRadius[LocationWithMixedTypes]("Lat", "Lng", center, 100))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for unsupported lng type, got %d", len(result))
	}
}

func TestWithinBoundingBoxMixedTypes(t *testing.T) {
	locations := []LocationWithMixedTypes{
		{Name: "Test", Lat: 10, Lng: true},
	}

	box := BoundingBox{
		SouthWest: Point{Lat: 0, Lng: 0},
		NorthEast: Point{Lat: 20, Lng: 20},
	}

	result := Apply(locations, WithinBoundingBox[LocationWithMixedTypes]("Lat", "Lng", box))
	if len(result) != 0 {
		t.Errorf("Expected 0 results for unsupported lng type, got %d", len(result))
	}
}

func TestSortByDistanceMixedTypes(t *testing.T) {
	locations := []LocationWithMixedTypes{
		{Name: "Test", Lat: 40.7128, Lng: true},
	}

	center := Point{Lat: 40.7128, Lng: -74.0060}
	sorted := SortByDistance(locations, "Lat", "Lng", center)
	if len(sorted) != 0 {
		t.Errorf("Expected 0 results for unsupported lng type, got %d", len(sorted))
	}
}

// Test float32 coordinates
type LocationFloat32 struct {
	Name string
	Lat  float32
	Lng  float32
}

func TestWithinRadiusFloat32(t *testing.T) {
	locations := []LocationFloat32{
		{Name: "Origin", Lat: 0, Lng: 0},
		{Name: "Near", Lat: 1, Lng: 1},
	}

	origin := Point{Lat: 0, Lng: 0}
	result := Apply(locations, WithinRadius[LocationFloat32]("Lat", "Lng", origin, 200))
	if len(result) != 2 {
		t.Errorf("Expected 2 locations within 200km of origin, got %d", len(result))
	}
}

func TestWithinBoundingBoxFloat32(t *testing.T) {
	locations := []LocationFloat32{
		{Name: "Inside", Lat: 10, Lng: 10},
	}

	box := BoundingBox{
		SouthWest: Point{Lat: 0, Lng: 0},
		NorthEast: Point{Lat: 20, Lng: 20},
	}

	result := Apply(locations, WithinBoundingBox[LocationFloat32]("Lat", "Lng", box))
	if len(result) != 1 {
		t.Errorf("Expected 1 location inside box, got %d", len(result))
	}
}

func TestSortByDistanceFloat32(t *testing.T) {
	locations := []LocationFloat32{
		{Name: "Far", Lat: 45, Lng: 90},
		{Name: "Origin", Lat: 0, Lng: 0},
	}

	origin := Point{Lat: 0, Lng: 0}
	sorted := SortByDistance(locations, "Lat", "Lng", origin)
	if len(sorted) != 2 {
		t.Errorf("Expected 2 sorted locations, got %d", len(sorted))
	}
	if sorted[0].Name != "Origin" {
		t.Errorf("Expected Origin first, got %s", sorted[0].Name)
	}
}
