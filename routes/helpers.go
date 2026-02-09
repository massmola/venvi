package routes

import (
	"strconv"
	"venvi/providers"

	"github.com/pocketbase/pocketbase/core"
)

// recordToEvent maps a PocketBase record to the internal Event struct.
func recordToEvent(r *core.Record) providers.Event {
	// Safely extract topics
	var topics []string
	rawTopics := r.Get("topics")
	if tSlice, ok := rawTopics.([]string); ok {
		topics = tSlice
	} else if tSlice, ok := rawTopics.([]any); ok {
		for _, v := range tSlice {
			if s, ok := v.(string); ok {
				topics = append(topics, s)
			}
		}
	}

	return providers.Event{
		ID:          r.Id,
		Title:       r.GetString("title"),
		Description: r.GetString("description"),
		DateStart:   r.GetDateTime("date_start").Time(), // PocketBase DateTime to Go Time
		DateEnd:     r.GetDateTime("date_end").Time(),
		Location:    r.GetString("location"),
		URL:         r.GetString("url"),
		ImageURL:    r.GetString("image_url"),
		SourceName:  r.GetString("source_name"),
		SourceID:    r.GetString("source_id"),
		Topics:      topics,
		Category:    r.GetString("category"),
		IsNew:       r.GetBool("is_new"),
		Latitude:    r.GetFloat("latitude"),
		Longitude:   r.GetFloat("longitude"),
	}
}

// eventsToMaps converts a slice of Events to a slice of maps for JSON response.
func eventsToMaps(events []providers.Event) []map[string]any {
	result := make([]map[string]any, len(events))
	for i, e := range events {
		result[i] = map[string]any{
			"id":          e.ID,
			"title":       e.Title,
			"description": e.Description,
			"date_start":  e.DateStart,
			"date_end":    e.DateEnd,
			"location":    e.Location,
			"url":         e.URL,
			"image_url":   e.ImageURL,
			"source_name": e.SourceName,
			"source_id":   e.SourceID,
			"topics":      e.Topics,
			"category":    e.Category,
			"is_new":      e.IsNew,
			"latitude":    e.Latitude,
			"longitude":   e.Longitude,
		}
	}
	return result
}

// castToFloat converts a string to float64, returning 0 if invalid.
func castToFloat(s string) float64 {
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
