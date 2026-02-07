package providers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// MuseionProvider fetches events from the Museion website via HTML scraping.
type MuseionProvider struct {
	BaseURL string
	Client  *http.Client
}

func NewMuseionProvider() *MuseionProvider {
	return &MuseionProvider{
		BaseURL: "https://www.museion.it/en/events",
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *MuseionProvider) SourceName() string {
	return "museion"
}

func (p *MuseionProvider) FetchEvents(ctx context.Context) ([]RawEvent, error) {
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

	// Selector based on analysis: classes like preview-item__title
	// I will iterate over the container items.

	doc.Find(".preview-item").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".preview-item__title").Text())
		if title == "" {
			return
		}

		link, _ := s.Attr("href")
		if !strings.HasPrefix(link, "http") {
			link = "https://www.museion.it" + link
		}

		// Extract image
		imgSrc, _ := s.Find("img").Attr("src")

		// Extract date/category
		meta := strings.TrimSpace(s.Find(".preview-item__meta").Text())

		raw := map[string]any{
			"title": title,
			"link":  link,
			"image": imgSrc,
			"meta":  meta,
		}
		events = append(events, RawEvent(raw))
	})

	return events, nil
}

func (p *MuseionProvider) MapEvent(raw RawEvent) *Event {
	title := fmt.Sprintf("%v", raw["title"])
	link := fmt.Sprintf("%v", raw["link"])
	image := fmt.Sprintf("%v", raw["image"])
	meta := fmt.Sprintf("%v", raw["meta"])

	// Generate ID
	id := fmt.Sprintf("museion-%d", time.Now().UnixNano())
	if link != "" {
		parts := strings.Split(link, "/")
		if len(parts) > 0 {
			id = parts[len(parts)-1]
		}
	}

	// Attempt to parse date from meta (e.g. "10.02.2026 | Event")
	// Simple heuristic
	dateStart := time.Now()
	parts := strings.Split(meta, "|")
	if len(parts) > 0 {
		dateStr := strings.TrimSpace(parts[0])
		// format dd.mm.yyyy
		if parsed, err := time.Parse("02.01.2006", dateStr); err == nil {
			dateStart = parsed
		}
	}

	return &Event{
		ID:          id,
		Title:       title,
		Description: meta, // Using meta as description for now
		DateStart:   dateStart,
		DateEnd:     dateStart.Add(2 * time.Hour),
		Location:    "Museion, Bolzano",
		URL:         link,
		ImageURL:    image,
		SourceName:  p.SourceName(),
		SourceID:    id,
		IsNew:       true,
		Category:    "Art & Culture",
	}
}
