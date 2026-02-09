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

// NewNOIProvider creates a new instance of NOIProvider targeted at the NOI Techpark.
func NewNOIProvider() *NOIProvider {
	return &NOIProvider{
		BaseURL: "https://tourism.api.opendatahub.com/v1/Event",
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// SourceName returns the unique identifier for this provider.
func (p *NOIProvider) SourceName() string {
	return "noi"
}

// FetchEvents retrieves raw event data from the Open Data Hub API for NOI Techpark.
func (p *NOIProvider) FetchEvents(ctx context.Context) ([]RawEvent, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	q := req.URL.Query()
	q.Set("pagesize", "200") // Validated safe page size for ODH
	q.Set("locationfilter", "Bolzano")

	var allEvents []RawEvent
	page := 1

	for {
		q.Set("pagenumber", fmt.Sprintf("%d", page))
		req.URL.RawQuery = q.Encode()

		// Clone the request for each iteration to avoid reusing body/context issues if any (though here we reuse client)
		// But wait, we can reuse the request object if we are careful, or just create a new one inside the loop.
		// Actually, creating a new request inside the loop is safer but higher overhead.
		// Given the `req` is created outside, we just update URL.
		// However, `req` body is nil for GET, so it's fine.

		resp, err := p.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetching events page %d: %w", page, err)
		}

		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("unexpected status page %d: %d", page, resp.StatusCode)
		}

		var localResult ODHResponse
		if err := json.NewDecoder(resp.Body).Decode(&localResult); err != nil {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("decoding response page %d: %w", page, err)
		}
		_ = resp.Body.Close()

		if len(localResult.Items) == 0 {
			break
		}

		hasNOI := func(s string) bool {
			s = strings.ToUpper(s)
			return strings.Contains(s, "NOI") || strings.Contains(s, "VOLTA") || strings.Contains(s, "TECHPARK")
		}

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
				allEvents = append(allEvents, RawEvent(item))
			}
		}

		// If we got fewer items than page size, we are done
		if len(localResult.Items) < 200 {
			break
		}
		page++
	}

	return allEvents, nil
}

// MapEvent reuses logic from helpers but sets source to "noi"
// MapEvent converts a RawEvent into the internal Event structure.
func (p *NOIProvider) MapEvent(raw RawEvent) *Event {
	// NOI events are always in NOI Techpark
	event := buildEventFromRaw(raw, p.SourceName(), "NOI Techpark", "https://noi.bz.it/en/events")

	// Ensure URL is specific to NOI if not already set by buildEventFromRaw (which uses default)
	if event.URL == "" {
		event.URL = "https://noi.bz.it/en/events"
	}

	return event
}
