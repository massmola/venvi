// Package providers implements event data source providers for Venvi.
package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

// Providers is the list of all registered event providers.
// Add new providers here to include them in sync operations.
var Providers = []EventProvider{
	NewODHProvider(),
	NewEuroHackathonsProvider(),
	NewDrinbzProvider(),
	NewNOIProvider(),
	NewUnibzProvider(),
	NewMuseionProvider(),
}

// SyncStats contains statistics about a sync operation.
type SyncStats struct {
	Provider string `json:"provider"`
	New      int    `json:"new"`
	Updated  int    `json:"updated"`
	Errors   int    `json:"errors"`
}

// SyncAllEvents synchronizes events from all registered providers.
// It fetches events from each provider, maps them to the unified format,
// and upserts them into the PocketBase events collection.
func SyncAllEvents(app core.App) (map[string]SyncStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	stats := make(map[string]SyncStats)

	for _, provider := range Providers {
		providerStats, err := syncProvider(ctx, app, provider)
		if err != nil {
			log.Printf("Error syncing %s: %v", provider.SourceName(), err)
			stats[provider.SourceName()] = SyncStats{
				Provider: provider.SourceName(),
				Errors:   1,
			}
			continue
		}
		stats[provider.SourceName()] = providerStats
	}

	return stats, nil
}

// syncProvider syncs events from a single provider.
func syncProvider(ctx context.Context, app core.App, provider EventProvider) (SyncStats, error) {
	stats := SyncStats{Provider: provider.SourceName()}

	// Fetch raw events
	rawEvents, err := provider.FetchEvents(ctx)
	if err != nil {
		return stats, fmt.Errorf("fetching events: %w", err)
	}

	// Get or create events collection
	collection, err := app.FindCollectionByNameOrId("events")
	if err != nil {
		return stats, fmt.Errorf("finding events collection: %w", err)
	}

	// Process each event
	for _, raw := range rawEvents {
		event := provider.MapEvent(raw)

		// Find existing record by source_name and source_id
		records, err := app.FindRecordsByFilter(
			collection,
			"source_name = {:source_name} && source_id = {:source_id}",
			"",
			1,
			0,
			map[string]any{
				"source_name": event.SourceName,
				"source_id":   event.SourceID,
			},
		)

		if err != nil || len(records) == 0 {
			// Event doesn't exist, create new
			record := core.NewRecord(collection)
			if err := populateRecord(record, event); err != nil {
				log.Printf("Error populating record: %v", err)
				stats.Errors++
				continue
			}

			if err := app.Save(record); err != nil {
				log.Printf("Error saving new event %s/%s: %v", event.SourceName, event.SourceID, err)
				stats.Errors++
				continue
			}
			stats.New++
		} else {
			// Event exists, update it (but preserve is_new status)
			existing := records[0]
			event.IsNew = existing.GetBool("is_new")
			if err := populateRecord(existing, event); err != nil {
				log.Printf("Error updating record: %v", err)
				stats.Errors++
				continue
			}

			if err := app.Save(existing); err != nil {
				log.Printf("Error updating event %s/%s: %v", event.SourceName, event.SourceID, err)
				stats.Errors++
				continue
			}
			stats.Updated++
		}
	}

	return stats, nil
}

// populateRecord fills a PocketBase record with event data.
func populateRecord(record *core.Record, event *Event) error {
	// Note: We don't set 'id' manually to allow PocketBase to generate a valid 15-char ID.
	// We use 'source_name' and 'source_id' for uniqueness.
	record.Set("title", event.Title)
	record.Set("description", event.Description)
	record.Set("date_start", event.DateStart)
	record.Set("date_end", event.DateEnd)
	record.Set("location", event.Location)
	record.Set("url", event.URL)
	record.Set("image_url", event.ImageURL)
	record.Set("source_name", event.SourceName)
	record.Set("source_id", event.SourceID)

	// Serialize topics as JSON
	topicsJSON, err := json.Marshal(event.Topics)
	if err != nil {
		return fmt.Errorf("marshaling topics: %w", err)
	}
	record.Set("topics", string(topicsJSON))

	record.Set("category", event.Category)
	record.Set("is_new", event.IsNew)

	return nil
}
