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
	// Filter for Bolzano events and narrowing down to NOI Techpark in-memory.
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

	var localResult ODHResponse
	if err := json.NewDecoder(resp.Body).Decode(&localResult); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	hasNOI := func(s string) bool {
		s = strings.ToUpper(s)
		return strings.Contains(s, "NOI") || strings.Contains(s, "VOLTA") || strings.Contains(s, "TECHPARK")
	}

	events := make([]RawEvent, 0, len(localResult.Items))

	for _, item := range localResult.Items {
		match := false

		// Check Title (multilingual)
		if details, ok := item["Detail"].(map[string]any); ok {
			for _, lang := range []string{"en", "it", "de"} {
				if langData, ok := details[lang].(map[string]any); ok {
					if title, ok := langData["Title"].(string); ok && hasNOI(title) {
						match = true
						break
					}
				}
			}
		}

		// Check Location
		if !match {
			if locInfo, ok := item["LocationInfo"].(map[string]any); ok {
				// Check District/Municipality names if available
				if district, ok := locInfo["DistrictInfo"].(map[string]any); ok {
					if nameMap, ok := district["Name"].(map[string]any); ok {
						for _, lang := range []string{"en", "it", "de"} {
							if name, ok := nameMap[lang].(string); ok && hasNOI(name) {
								match = true
								break
							}
						}
					}
				}
			}
		}

		// Also check if any ContactInfo address contains Volta/NOI
		if !match {
			if contacts, ok := item["ContactInfos"].(map[string]any); ok {
				for _, lang := range []string{"en", "it", "de"} {
					if contact, ok := contacts[lang].(map[string]any); ok {
						if addr, ok := contact["Address"].(string); ok && hasNOI(addr) {
							match = true
							break
						}
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

// MapEvent reuses logic from helpers but sets source to "noi"
func (p *NOIProvider) MapEvent(raw RawEvent) *Event {
	// NOI events are always in NOI Techpark
	event := buildEventFromRaw(raw, p.SourceName(), "NOI Techpark", "https://noi.bz.it/en/events")

	// Ensure URL is specific to NOI if not already set by buildEventFromRaw (which uses default)
	if event.URL == "" || event.URL == "https://opendatahub.com/events/"+event.ID {
		event.URL = "https://noi.bz.it/en/events"
	}

	return event
}
