// Package routes defines HTTP routes for the Venvi application.
package routes

import (
	"log"
	"net/http"

	"github.com/pocketbase/pocketbase/core"

	"venvi/providers"
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

		records, err := app.FindRecordsByFilter(
			collection,
			filter,
			"-date_start",
			100,
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

		// Convert records to JSON-friendly format
		events := make([]map[string]any, len(records))
		for i, r := range records {
			events[i] = map[string]any{
				"id":          r.Id,
				"title":       r.GetString("title"),
				"description": r.GetString("description"),
				"date_start":  r.GetDateTime("date_start"),
				"date_end":    r.GetDateTime("date_end"),
				"location":    r.GetString("location"),
				"url":         r.GetString("url"),
				"image_url":   r.GetString("image_url"),
				"source_name": r.GetString("source_name"),
				"source_id":   r.GetString("source_id"),
				"topics":      r.Get("topics"),
				"category":    r.GetString("category"),
				"is_new":      r.GetBool("is_new"),
			}
		}

		return e.JSON(http.StatusOK, events)
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
