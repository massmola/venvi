---
name: create_endpoint
description: Guide for adding a new API endpoint to the PocketBase application.
---
# Create Endpoint Skill

Follow these steps to add a new endpoint to the Venvi PocketBase application.

## 1. Plan the Endpoint
Define the following:
- **Path**: e.g., `/api/venvi/users/{id}`
- **Method**: GET, POST, PUT, DELETE
- **Request Body**: Struct for request parsing (if needed)
- **Response**: JSON response structure

## 2. Register the Route (`routes/api.go` or `routes/web.go`)

For API endpoints, edit `routes/api.go`:

```go
// In RegisterAPIRoutes function
se.Router.GET("/api/venvi/example/{id}", func(e *core.RequestEvent) error {
    id := e.Request.PathValue("id")
    
    // Access PocketBase app
    collection, err := app.FindCollectionByNameOrId("example")
    if err != nil {
        return e.NotFoundError("Collection not found", err)
    }
    
    record, err := app.FindRecordById(collection, id)
    if err != nil {
        return e.NotFoundError("Record not found", err)
    }
    
    return e.JSON(http.StatusOK, record)
})
```

For web/HTMX endpoints, edit `routes/web.go`.

## 3. Add Authentication (Optional)

Use PocketBase's built-in auth middleware:

```go
se.Router.POST("/api/venvi/protected", func(e *core.RequestEvent) error {
    // Handler code
    return e.JSON(http.StatusOK, map[string]any{"message": "success"})
}).Bind(apis.RequireAuth())
```

## 4. Write Tests (`tests/routes_test.go`)

```go
func TestExampleEndpoint(t *testing.T) {
    app := setupTestApp(t)
    
    resp, err := http.Get(app.BaseURL + "/api/venvi/example/123")
    require.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```

## 5. Verify
Run `./scripts/validate.sh` to ensure everything is correct.
