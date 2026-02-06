// Package main is the entry point for the Venvi EU Event Suggestion Platform.
// Venvi is built on top of PocketBase and provides event aggregation from
// multiple sources including Open Data Hub and EuroHackathons.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/template"

	"venvi/providers"
	"venvi/routes"
)

func main() {
	app := pocketbase.New()

	// Register JSVM plugin
	jsvm.MustRegister(app, jsvm.Config{})

	// Register migrate command plugin
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: true, // auto run migrations on serve
	})

	// Register routes and jobs on serve
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// Serve static files from pb_public
		se.Router.GET("/static/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		// Register web routes (HTMX pages)
		routes.RegisterWebRoutes(se, template.NewRegistry())

		// Register API routes
		routes.RegisterAPIRoutes(se, app)

		return se.Next()
	})

	// Register scheduled job for event sync
	app.Cron().MustAdd("sync_events", "0 */6 * * *", func() {
		log.Println("Running scheduled event sync...")
		stats, err := providers.SyncAllEvents(app)
		if err != nil {
			log.Printf("Sync failed: %v", err)
			return
		}
		log.Printf("Sync complete: %v", stats)
	})

	// Custom admin dashboard message
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/venvi/health", func(e *core.RequestEvent) error {
			return e.JSON(http.StatusOK, map[string]string{
				"status":  "healthy",
				"version": "0.1.0",
			})
		})
		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
