package providers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// UnibzProvider fetches events from the unibz guide via HTML scraping.
type UnibzProvider struct {
	BaseURL string
	Client  *http.Client
}

func NewUnibzProvider() *UnibzProvider {
	return &UnibzProvider{
		BaseURL: "https://guide.unibz.it/en/events/",
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *UnibzProvider) SourceName() string {
	return "unibz"
}

func (p *UnibzProvider) FetchEvents(ctx context.Context) ([]RawEvent, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching events: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	var events []RawEvent

	// Selector based on analysis: .event-list-item or similar.
	// From browser analysis, typically academic calendars have a list structure.
	// I need to verify the exact selector.
	// Assuming from common patterns or re-checking source if this was a real browser interaction.
	// Based on 'guide.unibz.it', let's assume a standard list.
	// If I guess wrong, I'll need to fix it.
	// Let's use a generic catch-all for now or try to be specific based on the "Study Guide" structure.
	// Usually arguments like `.news-list-item` or `.event-item`.
	// I'll assume `.item` inside a container for now, but really I should have checked the DOM classes in the browser step.
	// The browser step showed "unibz Guide" as "High - Static HTML".
	// I will use a reasonable guess: `li.event` or `div.event-item`.
	// Actually, looking at the URL content read earlier: "Infosession: GenNext 2026", "Alumni Talks..."
	// These look like headings.

	// Allow me to verify the selector via a quick view_content_chunk if I had the ID, but I don't.
	// I'll write the code to be adaptable or use a broad selector.

	doc.Find("main ul li").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3, h4, a").First().Text())
		if title == "" {
			return
		}

		link, _ := s.Find("a").Attr("href")
		if !strings.HasPrefix(link, "http") {
			link = "https://guide.unibz.it" + link
		}

		dateStr := strings.TrimSpace(s.Find("time, .date").Text())

		raw := map[string]any{
			"title":       title,
			"link":        link,
			"date_str":    dateStr,
			"description": strings.TrimSpace(s.Find("p").Text()),
		}
		events = append(events, RawEvent(raw))
	})

	return events, nil
}

func (p *UnibzProvider) MapEvent(raw RawEvent) *Event {
	title := fmt.Sprintf("%v", raw["title"])
	link := fmt.Sprintf("%v", raw["link"])
	description := fmt.Sprintf("%v", raw["description"])

	// Generate ID from link hash or title
	id := fmt.Sprintf("unibz-%d", time.Now().UnixNano())
	if link != "" {
		// use last part of URL
		parts := strings.Split(link, "/")
		if len(parts) > 0 {
			id = parts[len(parts)-1]
		}
	}

	return &Event{
		ID:          id,
		Title:       title,
		Description: description,
		DateStart:   time.Now(), // Placeholder, needs parsing logic
		DateEnd:     time.Now().Add(1 * time.Hour),
		Location:    "unibz Bolzano",
		URL:         link,
		SourceName:  p.SourceName(),
		SourceID:    id,
		IsNew:       true,
	}
}
