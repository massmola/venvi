// Package recommendations provides logic for scoring and sorting events based on user context.
package recommendations

import (
	"math"
	"sort"
	"time"

	"venvi/providers"
)

// RecommendationService handles the logic for scoring and sorting events
// based on user preferences and context.
type RecommendationService struct{}

// NewRecommendationService creates a new instance of RecommendationService.
func NewRecommendationService() *RecommendationService {
	return &RecommendationService{}
}

// UserContext represents the context of the user for recommendations.
// It includes the user's location and potential preferences (even if empty for now).
type UserContext struct {
	Latitude  float64
	Longitude float64
}

// ScoredEvent wraps an event with its calculated score.
type ScoredEvent struct {
	Event *providers.Event
	Score float64
}

const (
	// WeightDistance is the relative weight given to the distance score.
	WeightDistance = 0.6
	// WeightTime is the relative weight given to the time score.
	WeightTime = 0.3
	// WeightNew is the relative weight given to the newness score.
	WeightNew = 0.1

	// DistanceDecayConstant controls how quickly score drops with distance (in km).
	DistanceDecayConstant = 0.05
	// TimeDecayConstant controls how quickly score drops with time (hours).
	TimeDecayConstant = 0.01
)

// Recommend sorts the given events based on the user's context.
// It returns a new slice of events sorted by score (descending).
func (s *RecommendationService) Recommend(userCtx UserContext, events []providers.Event) []providers.Event {
	scoredEvents := make([]ScoredEvent, 0, len(events))

	for i := range events {
		score := s.Score(userCtx, &events[i])
		scoredEvents = append(scoredEvents, ScoredEvent{
			Event: &events[i],
			Score: score,
		})
	}

	// Sort by score descending
	sort.Slice(scoredEvents, func(i, j int) bool {
		return scoredEvents[i].Score > scoredEvents[j].Score
	})

	// Unwrap
	result := make([]providers.Event, len(scoredEvents))
	for i, se := range scoredEvents {
		result[i] = *se.Event
	}

	return result
}

// Score calculates a relevance score for a single event based on the user context.
// Higher score means more relevant.
func (s *RecommendationService) Score(userCtx UserContext, event *providers.Event) float64 {
	score := 0.0

	// 1. Distance Score
	// Calculate distance between user and event
	if event.Latitude != 0 && event.Longitude != 0 && userCtx.Latitude != 0 && userCtx.Longitude != 0 {
		dist := haversine(userCtx.Latitude, userCtx.Longitude, event.Latitude, event.Longitude)
		// Exponential decay based on distance: e^(-k * d)
		// At d=0, score=1. At d=large, score -> 0.
		distScore := math.Exp(-DistanceDecayConstant * dist)
		score += WeightDistance * distScore
	}

	// 2. Time Score
	// Events starting sooner are prioritized? or events starting soon but not too soon?
	// Let's assume we want events in the near future.
	now := time.Now()
	timeUntil := event.DateStart.Sub(now)
	if timeUntil > 0 {
		// Simple inverse or decay. Let's use a similar decay for time.
		// Convert to hours for easier tuning? Or just keep in duration.
		// Let's say preference is for events within next few days.
		hoursUntil := timeUntil.Hours()
		// Decay over days.
		// e^(-0.05 * hours) -> drops to ~0.3 after 24 hours? 0.05*24 = 1.2, e^-1.2 = 0.3. Maybe too fast.
		// e^(-0.01 * hours) -> 24h: e^-0.24 = 0.78. 48h: e^-0.48 = 0.61. 7 days (168h): e^-1.68 = 0.18.
		// This seems reasonable for "upcoming" events.
		timeScore := math.Exp(-TimeDecayConstant * hoursUntil)
		score += WeightTime * timeScore
	} else if event.DateEnd.After(now) {
		score += WeightTime * 0.5 // Ongoing events get flat medium score
	}

	// 3. Newness Score
	if event.IsNew {
		score += WeightNew
	}

	return score
}

// haversine calculates the distance in kilometers between two points.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in km
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	lat1Rad := lat1 * (math.Pi / 180.0)
	lat2Rad := lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
