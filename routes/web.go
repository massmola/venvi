// Package routes defines HTTP routes for the Venvi application.
package routes

import (
	"log"
	"net/http"

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

		records, err := app.FindRecordsByFilter(
			collection,
			"",            // no filter
			"-date_start", // sort by date descending
			100,           // limit
			0,             // offset
		)
		if err != nil {
			log.Printf("Error fetching events for partial: %v", err)
			return e.InternalServerError("Failed to fetch events", err)
		}

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
