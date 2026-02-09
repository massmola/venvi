package providers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// UnibzProvider fetches events from the unibz guide via HTML scraping.
type UnibzProvider struct {
	BaseURL string
	Client  *http.Client
}

// NewUnibzProvider creates a new instance of UnibzProvider.
func NewUnibzProvider() *UnibzProvider {
	return &UnibzProvider{
		BaseURL: "https://guide.unibz.it/en/events/",
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// SourceName returns the unique identifier for this provider.
func (p *UnibzProvider) SourceName() string {
	return "unibz"
}

// FetchEvents retrieves raw event data from the unibz guide website.
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

	doc.Find(".mediaItem").Each(func(_ int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".mediaItem_title a").Text())
		if title == "" {
			return
		}

		link, _ := s.Find(".mediaItem_title a").Attr("href")
		if !strings.HasPrefix(link, "http") {
			link = "https://guide.unibz.it" + link
		}

		// Date format: "10 Feb 2026 16:00-17:00"
		dateStrRaw := strings.TrimSpace(s.Find(".mediaItem_content .u-fw-bold").First().Text())

		var dateStart, dateEnd time.Time
		parts := strings.Split(dateStrRaw, " ")
		if len(parts) >= 4 {
			// date part: "10 Feb 2026"
			datePart := strings.Join(parts[:3], " ")
			// time part: "16:00-17:00"
			timeRange := parts[3]
			times := strings.Split(timeRange, "-")

			if len(times) >= 1 {
				startStr := datePart + " " + times[0]
				dateStart, _ = time.Parse("02 Jan 2006 15:04", startStr)
			}
			if len(times) >= 2 {
				endStr := datePart + " " + times[1]
				dateEnd, _ = time.Parse("02 Jan 2006 15:04", endStr)
			}
		}

		raw := map[string]any{
			"title":       title,
			"link":        link,
			"dateStart":   dateStart,
			"dateEnd":     dateEnd,
			"description": strings.TrimSpace(s.Find(".mediaItem_content .typography").Text()),
		}
		events = append(events, RawEvent(raw))
	})

	if len(events) == 0 {
		// Scrapers are fragile, log when potential structure change occurs.
		log.Printf("Warning: Unibz scraper found 0 events at %s (check for DOM structure changes)", p.BaseURL)
	}

	return events, nil
}

// MapEvent converts a RawEvent into the internal Event structure.
func (p *UnibzProvider) MapEvent(raw RawEvent) *Event {
	title, _ := raw["title"].(string)
	link, _ := raw["link"].(string)
	description, _ := raw["description"].(string)

	// Generate ID from link hash or title
	idBase := strconv.FormatInt(time.Now().UnixNano(), 10)
	if link != "" {
		normalizedLink := strings.TrimRight(link, "/")
		base := path.Base(normalizedLink)
		if base != "" && base != "." && base != "/" {
			idBase = base
		}
	}
	id := "unibz-" + idBase

	dateStart, _ := raw["dateStart"].(time.Time)
	dateEnd, _ := raw["dateEnd"].(time.Time)

	if dateStart.IsZero() {
		dateStart = time.Now()
	}
	if dateEnd.IsZero() {
		dateEnd = dateStart.Add(2 * time.Hour)
	}

	return &Event{
		ID:          id,
		Title:       title,
		Description: description,
		DateStart:   dateStart,
		DateEnd:     dateEnd,
		Location:    "unibz Bolzano",
		URL:         link,
		Category:    "Education",
		SourceName:  p.SourceName(),
		SourceID:    id,
		IsNew:       true,
		Topics:      []string{},
	}
}
