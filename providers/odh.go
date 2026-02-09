// Package providers implements event data source providers for Venvi.
package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ODHProvider fetches events from the Open Data Hub tourism API.
// The Open Data Hub (opendatahub.com) provides tourism events from South Tyrol.
type ODHProvider struct {
	// BaseURL allows overriding the API endpoint for testing.
	BaseURL string
	// Client is the HTTP client used for requests.
	Client *http.Client
}

// NewODHProvider creates a new ODH provider with default settings.
func NewODHProvider() *ODHProvider {
	return &ODHProvider{
		BaseURL: "https://tourism.opendatahub.com/v1/Event",
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// SourceName returns "odh" as the identifier for Open Data Hub events.
func (p *ODHProvider) SourceName() string {
	return "odh"
}

// ODHResponse represents the API response structure from Open Data Hub.
type ODHResponse struct {
	TotalResults int              `json:"TotalResults"`
	Items        []map[string]any `json:"Items"`
}

// FetchEvents retrieves events from the Open Data Hub API.
func (p *ODHProvider) FetchEvents(ctx context.Context) ([]RawEvent, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	q := req.URL.Query()
	q.Set("pagenumber", "1")
	q.Set("pagesize", "50") // Increased page size
	q.Set("active", "true")
	q.Set("odalactive", "true")
	q.Set("datefrom", time.Now().Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching events: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result ODHResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	events := make([]RawEvent, len(result.Items))
	for i, item := range result.Items {
		events[i] = RawEvent(item)
	}
	return events, nil
}

// MapEvent transforms raw ODH event data into a unified Event structure.
func (p *ODHProvider) MapEvent(raw RawEvent) *Event {
	// ODH events might not have a city in ContactInfos, so we fallback to "Unknown"
	// and let the helper try to find it.
	event := buildEventFromRaw(raw, p.SourceName(), "Unknown", "")

	// ODH-specific overrides if needed (currently none as buildEventFromRaw covers it)
	// But we need to ensure URL is correct if ID is present
	// ODH-specific overrides
	if event.ID != "" && event.URL == "" {
		event.URL = "https://opendatahub.com/events/" + event.ID
	}

	// Filter out bad quality events
	// 1. Title looks like a UUID (32 chars, hex)
	if len(event.Title) == 32 && isHex(event.Title) {
		return nil
	}
	// 2. No description (or very short)
	if len(event.Description) < 10 {
		return nil
	}

	return event
}

// isHex checks if a string is hexadecimal.
func isHex(s string) bool {
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}
