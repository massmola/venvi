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

func TestNOIProvider_FetchEvents(t *testing.T) {
	// Read fixture
	fixture, err := os.ReadFile("fixtures/noi.json")
	require.NoError(t, err)

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/v1/Event")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	// Initialize provider
	p := NewNOIProvider()
	p.BaseURL = server.URL + "/v1/Event"

	// Test Fetch
	events, err := p.FetchEvents(context.Background())
	require.NoError(t, err)
	assert.Len(t, events, 1)

	// Test Map
	mapped := p.MapEvent(events[0])
	assert.Equal(t, "noi-123", mapped.ID)
	assert.Contains(t, mapped.Title, "NOI Techpark Summit")
	assert.Equal(t, "noi", mapped.SourceName)
}
