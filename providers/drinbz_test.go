package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDrinbzProvider_FetchEvents(t *testing.T) {
	// Read fixture
	fixture, err := os.ReadFile("fixtures/drinbz.json")
	require.NoError(t, err)

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/wp-json/wp/v2/posts", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	// Initialize provider with mock server URL
	p := NewDrinbzProvider()
	p.BaseURL = server.URL + "/wp-json/wp/v2/posts"

	// Test Fetch
	events, err := p.FetchEvents(context.Background())
	require.NoError(t, err)
	assert.Len(t, events, 1)

	// Test Map
	mapped := p.MapEvent(events[0])
	assert.Equal(t, "12345", mapped.SourceID)
	assert.Equal(t, "Valentine's Party – Live Music", mapped.Title)
	assert.Equal(t, "drinbz", mapped.SourceName)
}
