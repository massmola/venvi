// Package providers implements event data source providers for Venvi.
// Each provider fetches events from an external API and maps them to the
// unified Event model stored in PocketBase.
package providers

import (
	"context"
	"time"
)

// RawEvent represents unprocessed event data from any source.
type RawEvent map[string]any

// Event represents a unified event structure for storage in PocketBase.
type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DateStart   time.Time `json:"date_start"`
	DateEnd     time.Time `json:"date_end"`
	Location    string    `json:"location"`
	URL         string    `json:"url"`
	ImageURL    string    `json:"image_url"`
	SourceName  string    `json:"source_name"`
	SourceID    string    `json:"source_id"`
	Topics      []string  `json:"topics"`
	Category    string    `json:"category"`
	IsNew       bool      `json:"is_new"`
}

// EventProvider defines the interface that all event sources must implement.
// This allows for easy addition of new data sources by implementing this interface.
type EventProvider interface {
	// SourceName returns the unique identifier for this data source.
	SourceName() string

	// FetchEvents retrieves raw event data from the external API.
	FetchEvents(ctx context.Context) ([]RawEvent, error)

	// MapEvent transforms raw event data into a unified Event structure.
	MapEvent(raw RawEvent) *Event
}
