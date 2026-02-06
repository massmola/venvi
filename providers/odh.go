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

// odhResponse represents the API response structure from Open Data Hub.
type odhResponse struct {
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
	q.Set("pagesize", "20")
	req.URL.RawQuery = q.Encode()

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching events: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result odhResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	events := make([]RawEvent, len(result.Items))
	for i, item := range result.Items {
		events[i] = RawEvent(item)
	}
	return events, nil
}

// getLocalized extracts a localized value from ODH detail objects.
// It tries English first, then Italian, then German.
func getLocalized(obj map[string]any, key string) string {
	for _, lang := range []string{"en", "it", "de"} {
		if langData, ok := obj[lang].(map[string]any); ok {
			if val, ok := langData[key].(string); ok && val != "" {
				return val
			}
		}
	}
	return ""
}

// MapEvent transforms raw ODH event data into a unified Event structure.
func (p *ODHProvider) MapEvent(raw RawEvent) *Event {
	details, _ := raw["Detail"].(map[string]any)
	if details == nil {
		details = map[string]any{}
	}

	title := getLocalized(details, "Title")
	if title == "" {
		title = "Untitled Event"
	}

	description := getLocalized(details, "BaseText")
	if description == "" {
		description = getLocalized(details, "IntroText")
	}

	// Extract location from ContactInfos
	location := "Unknown"
	if contactInfos, ok := raw["ContactInfos"].(map[string]any); ok {
		if enContact, ok := contactInfos["en"].(map[string]any); ok {
			if city, ok := enContact["City"].(string); ok && city != "" {
				location = city
			}
		}
	}

	// Extract image URL from gallery
	var imageURL string
	if gallery, ok := raw["ImageGallery"].([]any); ok && len(gallery) > 0 {
		if firstImg, ok := gallery[0].(map[string]any); ok {
			imageURL, _ = firstImg["ImageUrl"].(string)
		}
	}

	// Get raw ID or generate one
	rawID, _ := raw["Id"].(string)
	if rawID == "" {
		rawID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// Parse dates
	dateStart := time.Now()
	if dateStr, ok := raw["DateBegin"].(string); ok && dateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, dateStr); err == nil {
			dateStart = parsed
		} else if parsed, err := time.Parse("2006-01-02T15:04:05", dateStr); err == nil {
			dateStart = parsed
		}
	}

	dateEnd := time.Now()
	if dateStr, ok := raw["DateEnd"].(string); ok && dateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, dateStr); err == nil {
			dateEnd = parsed
		} else if parsed, err := time.Parse("2006-01-02T15:04:05", dateStr); err == nil {
			dateEnd = parsed
		}
	}

	return &Event{
		ID:          rawID,
		Title:       title,
		Description: description,
		DateStart:   dateStart,
		DateEnd:     dateEnd,
		Location:    location,
		URL:         fmt.Sprintf("https://opendatahub.com/events/%s", rawID),
		ImageURL:    imageURL,
		SourceName:  p.SourceName(),
		SourceID:    rawID,
		Topics:      []string{},
		Category:    "general",
		IsNew:       true,
	}
}
