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

func TestMuseionProvider_FetchEvents(t *testing.T) {
	fixture, err := os.ReadFile("fixtures/museion.html")
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/en/events", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	p := NewMuseionProvider()
	p.BaseURL = server.URL + "/en/events"

	events, err := p.FetchEvents(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, events)

	mapped := p.MapEvent(events[0])
	assert.Equal(t, "Hope – The Exhibition", mapped.Title)
	assert.Equal(t, "museion", mapped.SourceName)
	assert.Equal(t, "Art & Culture", mapped.Category)
}
