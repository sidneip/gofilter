package filter

import (
	"math"
	"reflect"
	"sort"
)

// Point represents a geographic coordinate with latitude and longitude
type Point struct {
	Lat float64
	Lng float64
}

// earthRadiusKm is the radius of Earth in kilometers
const earthRadiusKm = 6371.0

// degreesToRadians converts degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// calculateDistance calculates the Haversine distance between two points
// Returns distance in kilometers
func calculateDistance(p1, p2 Point) float64 {
	lat1 := degreesToRadians(p1.Lat)
	lng1 := degreesToRadians(p1.Lng)
	lat2 := degreesToRadians(p2.Lat)
	lng2 := degreesToRadians(p2.Lng)

	dLat := lat2 - lat1
	dLng := lng2 - lng1

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// WithinRadius returns a filter that checks if a location is within a specified radius from a center point
// latField and lngField are the struct field names for latitude and longitude
// centerPoint is the center point to compare against
// radiusKm is the radius in kilometers
func WithinRadius[T any](latField, lngField string, centerPoint Point, radiusKm float64) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		latValue, err := getFieldValue(item, latField)
		if err != nil {
			return false
		}

		lngValue, err := getFieldValue(item, lngField)
		if err != nil {
			return false
		}

		// Convert field values to float64
		var lat, lng float64

		switch latValue.Kind() {
		case reflect.Float64, reflect.Float32:
			lat = latValue.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			lat = float64(latValue.Int())
		default:
			return false // Unsupported type
		}

		switch lngValue.Kind() {
		case reflect.Float64, reflect.Float32:
			lng = lngValue.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			lng = float64(lngValue.Int())
		default:
			return false // Unsupported type
		}

		point := Point{Lat: lat, Lng: lng}
		distance := calculateDistance(centerPoint, point)

		return distance <= radiusKm
	})
}

// OutsideRadius returns a filter that checks if a location is outside a specified radius from a center point
func OutsideRadius[T any](latField, lngField string, centerPoint Point, radiusKm float64) Filter[T] {
	return Not(WithinRadius[T](latField, lngField, centerPoint, radiusKm))
}

// BoundingBox represents a geographic rectangle defined by southwest and northeast corners
type BoundingBox struct {
	SouthWest Point
	NorthEast Point
}

// WithinBoundingBox returns a filter that checks if a location is within a bounding box
func WithinBoundingBox[T any](latField, lngField string, box BoundingBox) Filter[T] {
	return FilterFunc[T](func(item T) bool {
		latValue, err := getFieldValue(item, latField)
		if err != nil {
			return false
		}

		lngValue, err := getFieldValue(item, lngField)
		if err != nil {
			return false
		}

		// Convert field values to float64
		var lat, lng float64

		switch latValue.Kind() {
		case reflect.Float64, reflect.Float32:
			lat = latValue.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			lat = float64(latValue.Int())
		default:
			return false // Unsupported type
		}

		switch lngValue.Kind() {
		case reflect.Float64, reflect.Float32:
			lng = lngValue.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			lng = float64(lngValue.Int())
		default:
			return false // Unsupported type
		}

		return lat >= box.SouthWest.Lat &&
			lat <= box.NorthEast.Lat &&
			lng >= box.SouthWest.Lng &&
			lng <= box.NorthEast.Lng
	})
}

// SortByDistance sorts a slice of items by distance from a center point
// This is a standalone function rather than a filter, as it returns a sorted slice
func SortByDistance[T any](items []T, latField, lngField string, centerPoint Point) []T {
	type itemWithDistance struct {
		item     T
		distance float64
	}

	itemsWithDistance := make([]itemWithDistance, 0, len(items))

	for _, item := range items {
		latValue, err := getFieldValue(item, latField)
		if err != nil {
			continue
		}

		lngValue, err := getFieldValue(item, lngField)
		if err != nil {
			continue
		}

		// Convert field values to float64
		var lat, lng float64

		switch latValue.Kind() {
		case reflect.Float64, reflect.Float32:
			lat = latValue.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			lat = float64(latValue.Int())
		default:
			continue // Unsupported type
		}

		switch lngValue.Kind() {
		case reflect.Float64, reflect.Float32:
			lng = lngValue.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			lng = float64(lngValue.Int())
		default:
			continue // Unsupported type
		}

		point := Point{Lat: lat, Lng: lng}
		distance := calculateDistance(centerPoint, point)

		itemsWithDistance = append(itemsWithDistance, itemWithDistance{
			item:     item,
			distance: distance,
		})
	}

	// Sort by distance
	sort.Slice(itemsWithDistance, func(i, j int) bool {
		return itemsWithDistance[i].distance < itemsWithDistance[j].distance
	})

	// Convert back to just items
	result := make([]T, 0, len(itemsWithDistance))
	for _, itemDist := range itemsWithDistance {
		result = append(result, itemDist.item)
	}

	return result
}
