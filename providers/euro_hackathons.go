// Package providers implements event data source providers for Venvi.
package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// EuroHackathonsProvider fetches hackathon events from the Euro Hackathons API.
type EuroHackathonsProvider struct {
	// BaseURL allows overriding the API endpoint for testing.
	BaseURL string
	// Client is the HTTP client used for requests.
	Client *http.Client
}

// NewEuroHackathonsProvider creates a new EuroHackathons provider with default settings.
func NewEuroHackathonsProvider() *EuroHackathonsProvider {
	return &EuroHackathonsProvider{
		BaseURL: "https://euro-hackathons.vercel.app/api/hackathons",
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// SourceName returns "euro_hackathons" as the identifier for hackathon events.
func (p *EuroHackathonsProvider) SourceName() string {
	return "euro_hackathons"
}

// hackathonsResponse represents the API response structure.
type hackathonsResponse struct {
	Data []map[string]any `json:"data"`
}

// FetchEvents retrieves hackathon events from the Euro Hackathons API.
func (p *EuroHackathonsProvider) FetchEvents(ctx context.Context) ([]RawEvent, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	q := req.URL.Query()
	q.Set("status", "upcoming")
	req.URL.RawQuery = q.Encode()

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching hackathons: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result hackathonsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	events := make([]RawEvent, len(result.Data))
	for i, item := range result.Data {
		events[i] = RawEvent(item)
	}
	return events, nil
}

// MapEvent transforms raw hackathon data into a unified Event structure.
func (p *EuroHackathonsProvider) MapEvent(raw RawEvent) *Event {
	// Extract required fields with defaults
	id, _ := raw["id"].(string)
	name, _ := raw["name"].(string)
	if name == "" {
		name = "Untitled Hackathon"
	}
	notes, _ := raw["notes"].(string)
	url, _ := raw["url"].(string)

	// Build location from city and country code
	city, _ := raw["city"].(string)
	countryCode, _ := raw["country_code"].(string)
	location := ""
	if city != "" && countryCode != "" {
		location = fmt.Sprintf("%s, %s", city, countryCode)
	} else if city != "" {
		location = city
	} else if countryCode != "" {
		location = countryCode
	}

	// Parse dates
	dateStart := time.Now()
	if dateStr, ok := raw["date_start"].(string); ok && dateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, dateStr); err == nil {
			dateStart = parsed
		} else if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			dateStart = parsed
		}
	}

	dateEnd := time.Now()
	if dateStr, ok := raw["date_end"].(string); ok && dateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, dateStr); err == nil {
			dateEnd = parsed
		} else if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			dateEnd = parsed
		}
	}

	// Extract topics
	var topics []string
	if rawTopics, ok := raw["topics"].([]any); ok {
		for _, t := range rawTopics {
			if topic, ok := t.(string); ok {
				topics = append(topics, topic)
			}
		}
	}

	return &Event{
		ID:          id,
		Title:       name,
		Description: notes,
		DateStart:   dateStart,
		DateEnd:     dateEnd,
		Location:    location,
		URL:         url,
		ImageURL:    "",
		SourceName:  p.SourceName(),
		SourceID:    id,
		Topics:      topics,
		Category:    "hackathon",
		IsNew:       true,
	}
}
