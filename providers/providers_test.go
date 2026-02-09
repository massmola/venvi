// Package providers implements event data source providers for Venvi.
package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockODHResponse is a sample response from the Open Data Hub API.
var mockODHResponse = map[string]any{
	"TotalResults": 1,
	"Items": []any{
		map[string]any{
			"Id": "test-event-1",
			"Detail": map[string]any{
				"en": map[string]any{
					"Title":    "Test Event",
					"BaseText": "Description of Test Event",
				},
			},
			"DateBegin":    "2024-01-01T10:00:00",
			"DateEnd":      "2024-01-01T12:00:00",
			"ContactInfos": map[string]any{"en": map[string]any{"City": "Bozen"}},
			"ImageGallery": []any{map[string]any{"ImageUrl": "http://example.com/image.jpg"}},
		},
	},
}

func TestODHProvider_FetchEvents(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(mockODHResponse)
		require.NoError(t, err)
	}))
	defer server.Close()

	// Create provider with mock URL
	provider := &ODHProvider{
		BaseURL: server.URL,
		Client:  &http.Client{Timeout: 5 * time.Second},
	}

	// Fetch events
	events, err := provider.FetchEvents(context.Background())
	require.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "test-event-1", events[0]["Id"])
}

func TestODHProvider_MapEvent(t *testing.T) {
	provider := NewODHProvider()

	raw := RawEvent{
		"Id": "test-123",
		"Detail": map[string]any{
			"en": map[string]any{
				"Title":    "My Event",
				"BaseText": "Event description",
			},
		},
		"DateBegin":    "2024-06-15T09:00:00",
		"DateEnd":      "2024-06-15T17:00:00",
		"ContactInfos": map[string]any{"en": map[string]any{"City": "Bolzano"}},
	}

	event := provider.MapEvent(raw)

	assert.Equal(t, "test-123", event.ID)
	assert.Equal(t, "My Event", event.Title)
	assert.Equal(t, "Event description", event.Description)
	assert.Equal(t, "Bolzano", event.Location)
	assert.Equal(t, "odh", event.SourceName)
	assert.Equal(t, "general", event.Category)
}

func TestODHProvider_MapEvent_EdgeCases(t *testing.T) {
	provider := NewODHProvider()

	// Minimal data - no title, no dates, no description
	// This should now be filtered out because description is too short/missing
	raw := RawEvent{
		"Detail": map[string]any{},
	}

	event := provider.MapEvent(raw)

	assert.Nil(t, event, "Event should be filtered out due to missing description")
}

func TestODHProvider_MapEvent_Valid(t *testing.T) {
	provider := NewODHProvider()

	raw := RawEvent{
		"Id": "test-valid",
		"Detail": map[string]any{
			"en": map[string]any{
				"Title":    "Valid Title of Event",
				"BaseText": "This is a valid description with enough length.",
			},
		},
	}

	event := provider.MapEvent(raw)
	assert.NotNil(t, event)
	assert.Equal(t, "test-valid", event.ID)
}

// mockHackathonsResponse is a sample response from Euro Hackathons API.
var mockHackathonsResponse = map[string]any{
	"data": []any{
		map[string]any{
			"id":           "hack-123",
			"name":         "EuroHack 2024",
			"notes":        "A great hackathon",
			"date_start":   "2024-03-15",
			"date_end":     "2024-03-17",
			"city":         "Berlin",
			"country_code": "DE",
			"url":          "https://eurohack.example.com",
			"topics":       []any{"ai", "blockchain"},
		},
	},
}

func TestEuroHackathonsProvider_FetchEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(mockHackathonsResponse)
		require.NoError(t, err)
	}))
	defer server.Close()

	provider := &EuroHackathonsProvider{
		BaseURL: server.URL,
		Client:  &http.Client{Timeout: 5 * time.Second},
	}

	events, err := provider.FetchEvents(context.Background())
	require.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, "hack-123", events[0]["id"])
}

func TestEuroHackathonsProvider_MapEvent(t *testing.T) {
	provider := NewEuroHackathonsProvider()

	raw := RawEvent{
		"id":           "hack-456",
		"name":         "DevConf",
		"notes":        "Developer conference",
		"date_start":   "2024-07-01",
		"date_end":     "2024-07-03",
		"city":         "Prague",
		"country_code": "CZ",
		"url":          "https://devconf.example.com",
		"topics":       []any{"devops", "cloud"},
	}

	event := provider.MapEvent(raw)

	assert.Equal(t, "hack-456", event.ID)
	assert.Equal(t, "DevConf", event.Title)
	assert.Equal(t, "Developer conference", event.Description)
	assert.Equal(t, "Prague, CZ", event.Location)
	assert.Equal(t, "euro_hackathons", event.SourceName)
	assert.Equal(t, "hackathon", event.Category)
	assert.Contains(t, event.Topics, "devops")
	assert.Contains(t, event.Topics, "cloud")
}
