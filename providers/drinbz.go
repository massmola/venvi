package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

// DrinbzProvider fetches events from the Drinbz WordPress API.
type DrinbzProvider struct {
	BaseURL string
	Client  *http.Client
}

// NewDrinbzProvider creates a new Drinbz provider.
func NewDrinbzProvider() *DrinbzProvider {
	return &DrinbzProvider{
		BaseURL: "https://drinbz.it/wp-json/wp/v2/posts",
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *DrinbzProvider) SourceName() string {
	return "drinbz"
}

// wpPost represents the WordPress post structure.
type wpPost struct {
	ID    int    `json:"id"`
	Date  string `json:"date"`
	Link  string `json:"link"`
	Title struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Content struct {
		Rendered string `json:"rendered"`
	} `json:"content"`
	Embedded map[string]any `json:"_embedded,omitempty"` // For images if needed
}

func (p *DrinbzProvider) FetchEvents(ctx context.Context) ([]RawEvent, error) {
	// Fetch recent posts. Drinbz posts are often events.
	// API: https://drinbz.it/wp-json/wp/v2/posts?per_page=20
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	q := req.URL.Query()
	q.Set("per_page", "20")
	req.URL.RawQuery = q.Encode()

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching posts: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var posts []wpPost
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	events := make([]RawEvent, 0, len(posts))
	log.Printf("Drinbz: fetched %d posts\n", len(posts))
	for _, post := range posts {
		// Convert struct to map for RawEvent interface
		// In a real scenario we might just use the struct directly if the interface allowed it,
		// but RawEvent is map[string]any.
		raw := make(map[string]any)
		raw["id"] = post.ID
		raw["date"] = post.Date
		raw["link"] = post.Link
		raw["title"] = post.Title.Rendered
		raw["content"] = post.Content.Rendered

		// Filter: Only consider posts that look like events (e.g. have a date in title or content, or category)
		// For now, we assume all posts on Drinbz "Next Week's Events" are relevant or we filter later.
		events = append(events, RawEvent(raw))
	}

	return events, nil
}

func (p *DrinbzProvider) MapEvent(raw RawEvent) *Event {
	id := fmt.Sprintf("%v", raw["id"])
	title := fmt.Sprintf("%v", raw["title"])
	link := fmt.Sprintf("%v", raw["link"])
	content := fmt.Sprintf("%v", raw["content"])
	dateStr := fmt.Sprintf("%v", raw["date"])

	// Parse date from WP format "2023-10-27T10:00:00"
	// Parse date from WP format "2023-10-27T10:00:00"
	dateStart, err := time.Parse("2006-01-02T15:04:05", dateStr)
	if err != nil {
		log.Printf("Drinbz: failed to parse date %q: %v\n", dateStr, err)
		dateStart = time.Now()
	}

	// Clean HTML from title (simple replacement)
	title = html.UnescapeString(title)

	return &Event{
		// Let PocketBase generate ID
		Title:       title,
		Description: content, // This contains HTML, might need stripping in future
		DateStart:   dateStart,
		DateEnd:     dateStart.Add(2 * time.Hour), // Default duration
		Location:    "Bolzano",                    // Default
		URL:         link,
		Category:    "Other", // Make sure category is set
		SourceName:  p.SourceName(),
		SourceID:    id,
		IsNew:       true,
		Topics:      []string{},
	}
}
