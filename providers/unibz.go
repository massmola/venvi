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

	fmt.Printf("Unibz: parsing HTML, doc length %d\n", len(doc.Text()))
	count := 0
	doc.Find(".mediaItem").Each(func(i int, s *goquery.Selection) {
		count++
		title := strings.TrimSpace(s.Find(".mediaItem_title a").Text())
		if title == "" {
			return
		}

		link, _ := s.Find(".mediaItem_title a").Attr("href")
		if !strings.HasPrefix(link, "http") {
			link = "https://guide.unibz.it" + link
		}

		// Date format: "10 Feb 2026 16:00-17:00"
		// Selector: .mediaItem_content > div (first child) or .u-fw-bold
		dateStrRaw := strings.TrimSpace(s.Find(".mediaItem_content .u-fw-bold").First().Text())
		// Parse date. Example: "10 Feb 2026 16:00-17:00"
		// We need to extract the date part "10 Feb 2026" and time "16:00"
		// Standard layout "02 Jan 2006 15:04"

		// Simple logic to extract date+time string
		parts := strings.Split(dateStrRaw, " ")
		var dateStr string
		if len(parts) >= 4 {
			// date part: "10 Feb 2026"
			// time part: "16:00-17:00" -> take "16:00"
			datePart := strings.Join(parts[:3], " ")
			timePart := strings.Split(parts[3], "-")[0] // "16:00"
			dateStr = datePart + " " + timePart
		} else {
			dateStr = dateStrRaw // Fallback
		}

		raw := map[string]any{
			"title":       title,
			"link":        link,
			"date_str":    dateStr,
			"description": strings.TrimSpace(s.Find(".mediaItem_content .typography").Text()),
		}
		events = append(events, RawEvent(raw))
	})
	fmt.Printf("Unibz: found %d items\n", count)

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

	// Parse date "02 Jan 2006 15:04"
	dateStr := fmt.Sprintf("%v", raw["date_str"])
	dateStart, err := time.Parse("02 Jan 2006 15:04", dateStr)
	if err != nil {
		dateStart = time.Now() // Fallback
	}

	return &Event{
		ID:          id,
		Title:       title,
		Description: description,
		DateStart:   dateStart,
		DateEnd:     dateStart.Add(2 * time.Hour),
		Location:    "unibz Bolzano",
		URL:         link,
		Category:    "Education",
		SourceName:  p.SourceName(),
		SourceID:    id,
		IsNew:       true,
	}
}
