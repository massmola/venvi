package recommendations

import (
	"testing"
	"time"

	"venvi/providers"

	"github.com/stretchr/testify/assert"
)

func TestScore(t *testing.T) {
	service := NewRecommendationService()
	userCtx := UserContext{
		Latitude:  40.7128,
		Longitude: -74.0060, // NYC
	}

	now := time.Now()

	tests := []struct {
		name     string
		event    providers.Event
		expected float64 // approximate expected score
	}{
		{
			name: "Close and Soon",
			event: providers.Event{
				Latitude:  40.7128,
				Longitude: -74.0060,
				DateStart: now.Add(24 * time.Hour), // Tomorrow
				IsNew:     false,
			},
			// Distance score = 1.0 (dist=0) * 0.6 = 0.6
			// Time score = e^(-0.01 * 24) = 0.78 * 0.3 = ~0.23
			// Total ~ 0.83
			expected: 0.83,
		},
		{
			name: "Far Away",
			event: providers.Event{
				Latitude:  34.0522,
				Longitude: -118.2437, // LA (~4000km away)
				DateStart: now.Add(24 * time.Hour),
				IsNew:     false,
			},
			// Distance score = e^(-0.05 * 4000) ~ 0 * 0.6 = 0
			// Time score ~ 0.23
			// Total ~ 0.23
			expected: 0.23,
		},
		{
			name: "New Event Boost",
			event: providers.Event{
				Latitude:  40.7128,
				Longitude: -74.0060,
				DateStart: now.Add(24 * time.Hour),
				IsNew:     true,
			},
			// Base ~ 0.83 + 0.1
			expected: 0.93,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			score := service.Score(userCtx, &tc.event)
			// Allow for some floating point variance
			assert.InDelta(t, tc.expected, score, 0.05, "Score should match expected value")
		})
	}
}

func TestRecommend(t *testing.T) {
	service := NewRecommendationService()
	userCtx := UserContext{
		Latitude:  40.7128,
		Longitude: -74.0060, // NYC
	}
	now := time.Now()

	events := []providers.Event{
		{
			ID:        "1",
			Title:     "Far Event",
			Latitude:  34.0522,
			Longitude: -118.2437, // LA
			DateStart: now.Add(24 * time.Hour),
		},
		{
			ID:        "2",
			Title:     "Near Event",
			Latitude:  40.7128,
			Longitude: -74.0060, // NYC
			DateStart: now.Add(24 * time.Hour),
		},
	}

	recommended := service.Recommend(userCtx, events)

	assert.Equal(t, 2, len(recommended))
	assert.Equal(t, "2", recommended[0].ID, "Near event should be first")
	assert.Equal(t, "1", recommended[1].ID, "Far event should be second")
}
