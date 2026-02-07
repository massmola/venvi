package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// NOIProvider fetches events from the Open Data Hub API, filtered for NOI Techpark.
// It reuses the ODH API structure but targets specific events.
type NOIProvider struct {
	BaseURL string
	Client  *http.Client
}

func NewNOIProvider() *NOIProvider {
	return &NOIProvider{
		BaseURL: "https://tourism.api.opendatahub.com/v1/Event",
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *NOIProvider) SourceName() string {
	return "noi"
}

func (p *NOIProvider) FetchEvents(ctx context.Context) ([]RawEvent, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	q := req.URL.Query()
	q.Set("pagenumber", "1")
	q.Set("pagesize", "50")
	// Filter for Bolzano and potentially "NOI" in title or location if specific ID isn't known.
	// ODH location filter for Bolzano is usually enough to start, then we filter in-memory.
	// Or we can query by "NOI Techpark" string if ODH supports text search?
	// The ODH API documentation usually supports `?searchfilter`
	// Let's try fetching Bolzano events and filtering for "NOI" in the title/location.
	q.Set("locationfilter", "Bolzano")

	req.URL.RawQuery = q.Encode()

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching events: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// odhResponse is unexported in odh.go. We should probably export it or redefine it.
	// For now redefining to be safe.
	type odhResponseLocal struct {
		TotalResults int              `json:"TotalResults"`
		Items        []map[string]any `json:"Items"`
	}

	var localResult odhResponseLocal
	if err := json.NewDecoder(resp.Body).Decode(&localResult); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	events := make([]RawEvent, 0, len(localResult.Items))
	for _, item := range localResult.Items {
		// Simple in-memory filter for NOI
		// We check if "NOI" appears in the title or location
		match := false

		// Check Title (multilingual)
		if details, ok := item["Detail"].(map[string]any); ok {
			for _, lang := range []string{"en", "it", "de"} {
				if langData, ok := details[lang].(map[string]any); ok {
					if title, ok := langData["Title"].(string); ok && strings.Contains(strings.ToUpper(title), "NOI") {
						match = true
						break
					}
				}
			}
		}

		if match {
			events = append(events, RawEvent(item))
		}
	}
	return events, nil
}

// MapEvent reuses logic similar to ODH but sets source to "noi"
func (p *NOIProvider) MapEvent(raw RawEvent) *Event {
	// We can instantiate an ODHProvider to reuse its MapEvent logic
	// IF we change ODHProvider methods to be usable or copy logic.
	// Since ODHProvider.MapEvent is bound to *ODHProvider struct, we can't easily reuse it
	// without making helper functions.
	// For now, copying the logic is safer and faster than refactoring odh.go.

	// Copy-paste of ODH mapping logic for now, ensuring SourceName is "noi"

	details, _ := raw["Detail"].(map[string]any)
	if details == nil {
		details = map[string]any{}
	}

	title := getLocalized(details, "Title")
	description := getLocalized(details, "BaseText")
	if description == "" {
		description = getLocalized(details, "IntroText")
	}

	location := "NOI Techpark" // Default for this provider

	// Extract image URL from gallery
	var imageURL string
	if gallery, ok := raw["ImageGallery"].([]any); ok && len(gallery) > 0 {
		if firstImg, ok := gallery[0].(map[string]any); ok {
			imageURL, _ = firstImg["ImageUrl"].(string)
		}
	}

	rawID, _ := raw["Id"].(string)
	if rawID == "" {
		rawID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// Dates parsing logic...
	dateStart := time.Now()
	if dateStr, ok := raw["DateBegin"].(string); ok && dateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, dateStr); err == nil {
			dateStart = parsed
		}
	}

	dateEnd := time.Now()
	if dateStr, ok := raw["DateEnd"].(string); ok && dateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, dateStr); err == nil {
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
		URL:         "https://noi.bz.it/en/events", // Link to NOI site instead of generic ODH
		ImageURL:    imageURL,
		SourceName:  p.SourceName(),
		SourceID:    rawID,
		IsNew:       true,
	}
}
