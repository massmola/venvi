// Package routes defines HTTP routes for the Venvi application.
package routes

import (
	"log"
	"net/http"

	"venvi/providers"
	"venvi/recommendations"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

// RegisterWebRoutes registers routes for serving HTMX-powered web pages.
func RegisterWebRoutes(se *core.ServeEvent, registry *template.Registry) {
	// Homepage
	se.Router.GET("/", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/layout.html",
			"views/index.html",
		).Render(map[string]any{
			"title": "Venvi - EU Event Suggestions",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// HTMX partial for event list
	se.Router.GET("/partials/events", func(e *core.RequestEvent) error {
		app := e.App

		// Fetch events from collection
		collection, err := app.FindCollectionByNameOrId("events")
		if err != nil {
			return e.InternalServerError("Collection not found", err)
		}

		// Check for authenticated user location
		var userLat, userLon float64
		if e.Auth != nil {
			userLat = e.Auth.GetFloat("latitude")
			userLon = e.Auth.GetFloat("longitude")
		}

		sortExpr := "+date_start" // Sort by date ascending (soonest first)
		limit := 500              // Fetch more candidates for re-ranking

		records, err := app.FindRecordsByFilter(
			collection,
			"date_end >= @now", // Only future events
			sortExpr,
			limit,
			0,
		)
		if err != nil {
			log.Printf("Error fetching events for partial: %v", err)
			return e.InternalServerError("Failed to fetch events", err)
		}

		// Apply recommendations (Always!)
		// Map records to internal events for sorting
		internalEvents := make([]providers.Event, len(records))
		recordMap := make(map[string]*core.Record)

		for i, r := range records {
			internalEvents[i] = recordToEvent(r)
			recordMap[r.Id] = r
		}

		svc := recommendations.NewRecommendationService()
		userCtx := recommendations.UserContext{
			Latitude:  userLat,
			Longitude: userLon,
		}
		sortedEvents := svc.Recommend(userCtx, internalEvents)

		// Reconstruct sorted records
		var sortedRecords []*core.Record
		for _, ev := range sortedEvents {
			if r, ok := recordMap[ev.ID]; ok {
				sortedRecords = append(sortedRecords, r)
			}
		}
		records = sortedRecords

		html, err := registry.LoadFiles(
			"views/partials/event_list.html",
		).Render(map[string]any{
			"events": records,
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})
}
