package tests

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/pocketbase/pocketbase/tools/template"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"venvi/providers"
	"venvi/routes"
)

func TestIntegration(t *testing.T) {
	// Ensure we run from project root so templates can be found
	if err := os.Chdir(".."); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// 1. Verify Schema and ID Constraints (Direct DB tests)
	t.Run("DatabaseLogic", func(t *testing.T) {
		testApp, err := tests.NewTestApp("pb_data")
		require.NoError(t, err)
		defer testApp.Cleanup()

		// Manually ensure the collection exists and has rules (to prevent migration sync issues in test)
		collection, err := testApp.FindCollectionByNameOrId("events")
		require.NoError(t, err)

		collection.ListRule = types.Pointer("")
		collection.ViewRule = types.Pointer("")
		err = testApp.Save(collection)
		require.NoError(t, err)

		expectedFields := []string{
			"title", "description", "date_start", "date_end",
			"location", "url", "image_url", "source_name",
			"source_id", "topics", "category", "is_new",
		}

		for _, fieldName := range expectedFields {
			field := collection.Fields.GetByName(fieldName)
			assert.NotNil(t, field, "Field %s should exist in events collection", fieldName)
		}

		// Sync ID Constraints Handling
		longID := "this-is-a-very-long-id-that-exceeds-pocketbase-limit"
		event := &providers.Event{
			ID:         longID,
			Title:      "Test Long ID Event",
			DateStart:  time.Now(),
			DateEnd:    time.Now().Add(time.Hour),
			URL:        "http://example.com",
			SourceName: "test_source",
			SourceID:   longID,
			Category:   "test",
		}

		record := core.NewRecord(collection)
		record.Set("title", event.Title)
		record.Set("date_start", event.DateStart)
		record.Set("date_end", event.DateEnd)
		record.Set("url", event.URL)
		record.Set("source_name", event.SourceName)
		record.Set("source_id", event.SourceID)
		record.Set("category", event.Category)

		err = testApp.Save(record)
		require.NoError(t, err, "Should save record even if source ID is long")
		assert.Len(t, record.Id, 15, "PocketBase should generate a 15-character ID")
	})

	// 2. Verify Routes using ApiScenario
	scenarios := []tests.ApiScenario{
		{
			Name:           "HealthCheck",
			Method:         http.MethodGet,
			URL:            "/api/venvi/health",
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"status"`,
				`"healthy"`,
			},
			TestAppFactory: func(t testing.TB) *tests.TestApp {
				app, _ := tests.NewTestApp("pb_data")
				return app
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				// Register health route explicitly since it's in main.go
				e.Router.GET("/api/venvi/health", func(er *core.RequestEvent) error {
					return er.JSON(http.StatusOK, map[string]string{"status": "healthy"})
				})
			},
		},
		{
			Name:            "EventsAPI",
			Method:          http.MethodGet,
			URL:             "/api/venvi/events",
			ExpectedStatus:  http.StatusOK,
			ExpectedContent: []string{`[`}, // Expect a JSON array
			TestAppFactory: func(t testing.TB) *tests.TestApp {
				app, _ := tests.NewTestApp("pb_data")
				return app
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				// Ensure collection has list rule
				collection, _ := app.FindCollectionByNameOrId("events")
				collection.ListRule = types.Pointer("")
				if err := app.Save(collection); err != nil {
					t.Fatalf("failed to save collection rule: %v", err)
				}

				routes.RegisterAPIRoutes(e, app)
			},
		},
		{
			Name:           "WebHome",
			Method:         http.MethodGet,
			URL:            "/",
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				"<title>", // Basic check for HTML
			},
			TestAppFactory: func(t testing.TB) *tests.TestApp {
				app, _ := tests.NewTestApp("pb_data")
				return app
			},
			BeforeTestFunc: func(t testing.TB, app *tests.TestApp, e *core.ServeEvent) {
				// Register web routes.
				// Template registry needs views relative to test run.
				// We'll set CWD or specify absolute paths if needed.
				routes.RegisterWebRoutes(e, template.NewRegistry())
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
