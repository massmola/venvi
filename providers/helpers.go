package providers

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

// extractLocalized extracts a localized value from ODH detail objects.
// It tries English first, then Italian, then German.
func extractLocalized(obj map[string]any, key string) string {
	for _, lang := range []string{"en", "it", "de"} {
		if langData, ok := obj[lang].(map[string]any); ok {
			if val, ok := langData[key].(string); ok && val != "" {
				return val
			}
		}
	}
	return ""
}

// extractLocalizedDetails extracts title and description from raw event data.
func extractLocalizedDetails(raw RawEvent) (string, string) {
	details, _ := raw["Detail"].(map[string]any)
	if details == nil {
		details = map[string]any{}
	}

	title := extractLocalized(details, "Title")
	if title == "" {
		title = "Untitled Event"
	}

	description := extractLocalized(details, "BaseText")
	if description == "" {
		description = extractLocalized(details, "IntroText")
	}
	return title, description
}

// extractImageURL extracts the first image URL from the event gallery.
func extractImageURL(raw RawEvent) string {
	if gallery, ok := raw["ImageGallery"].([]any); ok && len(gallery) > 0 {
		if firstImg, ok := gallery[0].(map[string]any); ok {
			if url, ok := firstImg["ImageUrl"].(string); ok {
				return url
			}
		}
	}
	return ""
}

// parseDates extracts start and end dates. Returns error if date strings present but malformed.
func parseDates(raw RawEvent) (time.Time, time.Time, error) {
	parse := func(key string) (time.Time, error) {
		dateStr, ok := raw[key].(string)
		if !ok || dateStr == "" {
			return time.Now(), nil
		}
		if parsed, err := time.Parse(time.RFC3339, dateStr); err == nil {
			return parsed, nil
		}
		if parsed, err := time.Parse("2006-01-02T15:04:05", dateStr); err == nil {
			return parsed, nil
		}
		return time.Now(), fmt.Errorf("malformed date %q", dateStr)
	}

	start, errStart := parse("DateBegin")
	end, errEnd := parse("DateEnd")

	var err error
	if errStart != nil {
		err = fmt.Errorf("parsing start date: %w", errStart)
	} else if errEnd != nil {
		err = fmt.Errorf("parsing end date: %w", errEnd)
	}

	// If end is before start (or equal/default), make it at least start time
	if end.Before(start) {
		end = start
	}

	return start, end, err
}

// extractLocation extracts city name from contact info.
func extractLocation(raw RawEvent) string {
	if contactInfos, ok := raw["ContactInfos"].(map[string]any); ok {
		if enContact, ok := contactInfos["en"].(map[string]any); ok {
			if city, ok := enContact["City"].(string); ok && city != "" {
				return city
			}
		}
	}
	return "Unknown"
}

// extractGPS extracts latitude and longitude.
func extractGPS(raw RawEvent) (float64, float64) {
	var lat, long float64
	if gpsInfo, ok := raw["GpsInfo"].([]any); ok && len(gpsInfo) > 0 {
		if firstGps, ok := gpsInfo[0].(map[string]any); ok {
			lat, _ = firstGps["Latitude"].(float64)
			long, _ = firstGps["Longitude"].(float64)
		}
	}
	if lat == 0 && long == 0 {
		lat, _ = raw["Latitude"].(float64)
		long, _ = raw["Longitude"].(float64)
	}
	return lat, long
}

// buildEventFromRaw maps common fields.
func buildEventFromRaw(raw RawEvent, sourceName, defaultLocation, defaultURL string) *Event {
	title, description := extractLocalizedDetails(raw)
	imageURL := extractImageURL(raw)
	dateStart, dateEnd, err := parseDates(raw)
	if err != nil {
		// Log warning but continue with fallback to time.Now()
		// (handled by parseDates returning time.Now() alongside the error)
		// Note: We use raw["Id"] for context if available
		rawID, _ := raw["Id"].(string)
		log.Printf("Warning: failed to parse dates for event %s from %s: %v", rawID, sourceName, err)
	}

	location := extractLocation(raw)
	if location == "Unknown" && defaultLocation != "" {
		location = defaultLocation
	}

	lat, long := extractGPS(raw)

	// Get raw ID or generate one
	rawID, _ := raw["Id"].(string)
	if rawID == "" {
		rawID = strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	url := defaultURL
	if url == "" {
		url = "https://opendatahub.com/events/" + rawID
	}

	return &Event{
		ID:          rawID,
		Title:       title,
		Description: description,
		DateStart:   dateStart,
		DateEnd:     dateEnd,
		Location:    location,
		URL:         url,
		ImageURL:    imageURL,
		SourceName:  sourceName,
		SourceID:    rawID,
		Topics:      []string{},
		Category:    "general",
		IsNew:       true,
		Latitude:    lat,
		Longitude:   long,
	}
}
