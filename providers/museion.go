package providers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// MuseionProvider fetches events from the Museion website via HTML scraping.
type MuseionProvider struct {
	BaseURL string
	Client  *http.Client
}

// NewMuseionProvider creates a new instance of MuseionProvider.
func NewMuseionProvider() *MuseionProvider {
	return &MuseionProvider{
		BaseURL: "https://www.museion.it/en/events",
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// SourceName returns the unique identifier for this provider.
func (p *MuseionProvider) SourceName() string {
	return "museion"
}

// FetchEvents retrieves raw event data from the Museion website.
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

	doc.Find(".preview-item").Each(func(_ int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".preview-item__title").Text())
		if title == "" {
			return
		}

		link, exists := s.Find("a").Attr("href")
		if !exists || link == "" {
			return
		}
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

// MapEvent converts a RawEvent into the internal Event structure.
func (p *MuseionProvider) MapEvent(raw RawEvent) *Event {
	title, _ := raw["title"].(string)
	link, _ := raw["link"].(string)
	image, _ := raw["image"].(string)
	meta, _ := raw["meta"].(string)

	// Generate ID
	var id string
	normalizedLink := strings.TrimRight(link, "/")
	if normalizedLink != "" {
		id = path.Base(normalizedLink)
	}
	// Fallback if ID is empty or invalid path
	if id == "" || id == "." || id == "/" {
		hash := sha256.Sum256([]byte(link))
		id = hex.EncodeToString(hash[:8])
	}
	id = "museion-" + id

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
		Topics:      []string{},
	}
}
