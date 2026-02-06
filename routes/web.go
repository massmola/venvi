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

	// Design 1: Neon/Cyberpunk
	se.Router.GET("/design/1", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design1.html",
		).Render(map[string]any{
			"title": "Venvi - Cyberpunk",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 2: Minimal/Swiss
	se.Router.GET("/design/2", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design2.html",
		).Render(map[string]any{
			"title": "Venvi - Minimal",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 3: Organic/Nature
	se.Router.GET("/design/3", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design3.html",
		).Render(map[string]any{
			"title": "Venvi - Organic",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 4: Brutalist
	se.Router.GET("/design/4", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design4.html",
		).Render(map[string]any{
			"title": "Venvi - Brutalist",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 5: Luxury/Glassmorphism
	se.Router.GET("/design/5", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design5.html",
		).Render(map[string]any{
			"title": "Venvi - Luxury",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 6: Pixel Art
	se.Router.GET("/design/6", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design6.html",
		).Render(map[string]any{
			"title": "Venvi - Pixel Art",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 7: Papercut
	se.Router.GET("/design/7", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design7.html",
		).Render(map[string]any{
			"title": "Venvi - Papercut",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 8: Art Nouveau
	se.Router.GET("/design/8", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design8.html",
		).Render(map[string]any{
			"title": "Venvi - Art Nouveau",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 9: Art Deco
	se.Router.GET("/design/9", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design9.html",
		).Render(map[string]any{
			"title": "Venvi - Art Deco",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 10: Pointillism
	se.Router.GET("/design/10", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design10.html",
		).Render(map[string]any{
			"title": "Venvi - Pointillism",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 11: Neumorphism (Restored)
	se.Router.GET("/design/11", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design11.html",
		).Render(map[string]any{
			"title": "Venvi - Soft UI",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 12: Papercut Monochrome
	se.Router.GET("/design/12", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design12.html",
		).Render(map[string]any{
			"title": "Venvi - Papercut Mono",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 13: Papercut Vibrant
	se.Router.GET("/design/13", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design13.html",
		).Render(map[string]any{
			"title": "Venvi - Papercut Vibrant",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 14: Papercut Noir
	se.Router.GET("/design/14", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design14.html",
		).Render(map[string]any{
			"title": "Venvi - Papercut Noir",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 15: Cardstock (Middleground)
	se.Router.GET("/design/15", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design15.html",
		).Render(map[string]any{
			"title": "Venvi - Cardstock",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 16: Golden Age (Klimt)
	se.Router.GET("/design/16", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design16.html",
		).Render(map[string]any{
			"title": "Venvi - Klimt",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 17: Metro (Ironwork)
	se.Router.GET("/design/17", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design17.html",
		).Render(map[string]any{
			"title": "Venvi - Metro",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 18: La Belle Époque (Poster)
	se.Router.GET("/design/18", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design18.html",
		).Render(map[string]any{
			"title": "Venvi - Belle Époque",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 19: Tiffany (Stained Glass)
	se.Router.GET("/design/19", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design19.html",
		).Render(map[string]any{
			"title": "Venvi - Tiffany",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 20: Whiplash (Botanical)
	se.Router.GET("/design/20", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design20.html",
		).Render(map[string]any{
			"title": "Venvi - Whiplash",
		})
		if err != nil {
			return e.InternalServerError("Template error", err)
		}
		return e.HTML(http.StatusOK, html)
	})

	// Design 21: High-Def Papercut
	se.Router.GET("/design/21", func(e *core.RequestEvent) error {
		html, err := registry.LoadFiles(
			"views/design21.html",
		).Render(map[string]any{
			"title": "Venvi - High-Def Papercut",
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
