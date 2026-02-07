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

func TestUnibzProvider_FetchEvents(t *testing.T) {
	fixture, err := os.ReadFile("fixtures/unibz.html")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/en/events/", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	p := NewUnibzProvider()
	p.BaseURL = server.URL + "/en/events/"

	events, err := p.FetchEvents(context.Background())
	require.NoError(t, err)

	// Assuming fixture has 2 items but one external link
	// The implementation might fetch all.
	assert.NotEmpty(t, events)

	// Verify mapping
	mapped := p.MapEvent(events[0])
	assert.Contains(t, mapped.Title, "Infosession: GenNext 2026")
	assert.Equal(t, "unibz", mapped.SourceName)
}
