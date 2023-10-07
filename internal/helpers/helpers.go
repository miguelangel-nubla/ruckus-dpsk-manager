package helpers

import (
	"fmt"
	"strconv"
	"time"
)

func ParseTimestamp(input string) (time.Time, error) {
	// Try to parse as a Unix timestamp first
	if timestamp, err := strconv.ParseInt(input, 10, 32); err == nil {
		return time.Unix(timestamp, 0), nil
	}

	var dateTimeFormats = []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339 with offset
		"2006-01-02T15:04:05Z",      // RFC3339 without offset
		"2006-01-02 15:04:05",       // YYYY-MM-DD HH:MM:SS
	}

	// If not a Unix timestamp, try to parse as a date-time string
	for _, format := range dateTimeFormats {
		t, err := time.Parse(format, input)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid input format")
}
