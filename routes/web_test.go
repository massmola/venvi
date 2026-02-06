package routes

import (
	"math"
	"testing"
)

func TestHaversine(t *testing.T) {
	// Berlin (52.5200, 13.4050) to Munich (48.1351, 11.5820)
	// Approx 504 km
	berlinLat, berlinLon := 52.5200, 13.4050
	munichLat, munichLon := 48.1351, 11.5820

	dist := haversine(berlinLat, berlinLon, munichLat, munichLon)

	if math.Abs(dist-504) > 5 { // Allow 5km error margin for different Earth radius constants
		t.Errorf("Expected approx 504km, got %f", dist)
	}

	// 0 distance
	if d := haversine(10, 10, 10, 10); d != 0 {
		t.Errorf("Expected 0 distance, got %f", d)
	}
}
