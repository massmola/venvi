// Package routes defines HTTP routes for the Venvi application.
package routes

import (
	"log"
	"net/http"

	"github.com/pocketbase/pocketbase/core"

	"venvi/providers"
	"venvi/recommendations"
)

// RegisterAPIRoutes registers API endpoints for programmatic access.
func RegisterAPIRoutes(se *core.ServeEvent, app core.App) {
	// List events with optional filters
	se.Router.GET("/api/venvi/events", func(e *core.RequestEvent) error {
		collection, err := app.FindCollectionByNameOrId("events")
		if err != nil {
			return e.NotFoundError("Events collection not found", err)
		}

		// Build filter from query params
		category := e.Request.URL.Query().Get("category")
		source := e.Request.URL.Query().Get("source")
		lat := e.Request.URL.Query().Get("lat")
		long := e.Request.URL.Query().Get("long")

		filter := ""
		if category != "" {
			filter = "category = {:category}"
		}
		if source != "" {
			if filter != "" {
				filter += " && "
			}
			filter += "source_name = {:source}"
		}

		// If location is provided, we might want to fetch more events to sort them effectively
		// We fetch more events to allow the recommendation engine to re-rank them
		limit := 500

		// Default filter: future events only
		if filter == "" {
			filter = "date_end >= @now"
		} else {
			filter += " && date_end >= @now"
		}

		records, err := app.FindRecordsByFilter(
			collection,
			filter,
			"+date_start", // Default sort: ascending (soonest first)
			limit,
			0,
			map[string]any{
				"category": category,
				"source":   source,
			},
		)
		if err != nil {
			log.Printf("Error fetching events API: %v", err)
			return e.InternalServerError("Failed to fetch events", err)
		}

		// Convert to internal events
		internalEvents := make([]providers.Event, len(records))
		for i, r := range records {
			internalEvents[i] = recordToEvent(r)
		}

		// Apply recommendations (Always!)
		// If location is missing, it will rely on time and newness scores
		userLat := castToFloat(lat)
		userLong := castToFloat(long)

		svc := recommendations.NewRecommendationService()
		userCtx := recommendations.UserContext{
			Latitude:  userLat,
			Longitude: userLong,
		}
		internalEvents = svc.Recommend(userCtx, internalEvents)

		// Convert records to JSON-friendly format
		// We use the helper which returns []map[string]any
		return e.JSON(http.StatusOK, eventsToMaps(internalEvents))
	})

	// Trigger manual sync
	se.Router.POST("/api/venvi/sync", func(e *core.RequestEvent) error {
		stats, err := providers.SyncAllEvents(app)
		if err != nil {
			return e.InternalServerError("Sync failed", err)
		}

		// Calculate totals
		totalNew := 0
		totalUpdated := 0
		for _, s := range stats {
			totalNew += s.New
			totalUpdated += s.Updated
		}

		return e.JSON(http.StatusOK, map[string]any{
			"message":       "Sync complete",
			"providers":     stats,
			"total_new":     totalNew,
			"total_updated": totalUpdated,
		})
	})
}
