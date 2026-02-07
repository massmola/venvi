// Package routes defines HTTP routes for the Venvi application.
package routes

import (
	"log"
	"math"
	"net/http"
	"sort"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

// haversine calculates the distance between two points in kilometers.
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
		hasLocation := false
		if e.Auth != nil {
			userLat = e.Auth.GetFloat("latitude")
			userLon = e.Auth.GetFloat("longitude")
			if userLat != 0 || userLon != 0 {
				hasLocation = true
			}
		}

		sortExpr := "-date_start"
		limit := 100
		if hasLocation {
			// Fetch more events if we need to find the closest ones
			limit = 500
		}

		records, err := app.FindRecordsByFilter(
			collection,
			"",       // no filter
			sortExpr, // sort by date descending initially
			limit,    // limit
			0,        // offset
		)
		if err != nil {
			log.Printf("Error fetching events for partial: %v", err)
			return e.InternalServerError("Failed to fetch events", err)
		}

		// If user has location, perform in-memory sort by distance
		if hasLocation {
			sort.SliceStable(records, func(i, j int) bool {
				lat1 := records[i].GetFloat("latitude")
				lon1 := records[i].GetFloat("longitude")
				dist1 := haversine(userLat, userLon, lat1, lon1)

				lat2 := records[j].GetFloat("latitude")
				lon2 := records[j].GetFloat("longitude")
				dist2 := haversine(userLat, userLon, lat2, lon2)

				return dist1 < dist2
			})
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
