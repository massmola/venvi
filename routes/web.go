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

// renderEventPartial fetches recommended events and renders the given partial template.
func renderEventPartial(e *core.RequestEvent, registry *template.Registry, partialPath string) error {
	app := e.App

	collection, err := app.FindCollectionByNameOrId("events")
	if err != nil {
		return e.InternalServerError("Collection not found", err)
	}

	var userLat, userLon float64
	if e.Auth != nil {
		userLat = e.Auth.GetFloat("latitude")
		userLon = e.Auth.GetFloat("longitude")
	}

	records, err := app.FindRecordsByFilter(
		collection,
		"date_end >= @now",
		"+date_start",
		500,
		0,
	)
	if err != nil {
		log.Printf("Error fetching events for partial: %v", err)
		return e.InternalServerError("Failed to fetch events", err)
	}

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

	var sortedRecords []*core.Record
	for _, ev := range sortedEvents {
		if r, ok := recordMap[ev.ID]; ok {
			sortedRecords = append(sortedRecords, r)
		}
	}

	html, err := registry.LoadFiles(partialPath).Render(map[string]any{
		"events": sortedRecords,
	})
	if err != nil {
		return e.InternalServerError("Template error", err)
	}
	return e.HTML(http.StatusOK, html)
}

// RegisterWebRoutes registers routes for serving HTMX-powered web pages.
func RegisterWebRoutes(se *core.ServeEvent, registry *template.Registry) {
	// Homepage (existing design)
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

	// HTMX partial for event list (original)
	se.Router.GET("/partials/events", func(e *core.RequestEvent) error {
		return renderEventPartial(e, registry, "views/partials/event_list.html")
	})

	// ── Design Previews ──

	// Design A: Vintage Journal
	se.Router.GET("/design/a", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design_a_layout.html",
			"views/design_a_index.html",
		).Render(map[string]any{
			"title": "Venvi — Vintage Journal",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	se.Router.GET("/partials/events/a", func(e *core.RequestEvent) error {
		return renderEventPartial(e, registry, "views/partials/event_list_a.html")
	})

	// Design B: Origami Fold
	se.Router.GET("/design/b", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design_b_layout.html",
			"views/design_b_index.html",
		).Render(map[string]any{
			"title": "VENVI — Origami Fold",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	se.Router.GET("/partials/events/b", func(e *core.RequestEvent) error {
		return renderEventPartial(e, registry, "views/partials/event_list_b.html")
	})

	// Design C: Paper Collage
	se.Router.GET("/design/c", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design_c_layout.html",
			"views/design_c_index.html",
		).Render(map[string]any{
			"title": "venvi — Paper Collage",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	se.Router.GET("/partials/events/c", func(e *core.RequestEvent) error {
		return renderEventPartial(e, registry, "views/partials/event_list_c.html")
	})
}
