package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseAge parses an age string like "24h" or "7d" into a time.Time
// Returns nil if age is empty
func ParseAge(age string) (*time.Time, error) {
	if age == "" {
		return nil, nil
	}

	// Parse the age string
	var duration time.Duration
	if strings.HasSuffix(age, "h") {
		// Hours
		hours, err := strconv.Atoi(strings.TrimSuffix(age, "h"))
		if err != nil {
			return nil, fmt.Errorf("invalid hour format: %s", age)
		}
		duration = time.Duration(hours) * time.Hour
	} else if strings.HasSuffix(age, "d") {
		// Days
		days, err := strconv.Atoi(strings.TrimSuffix(age, "d"))
		if err != nil {
			return nil, fmt.Errorf("invalid day format: %s", age)
		}
		duration = time.Duration(days) * 24 * time.Hour
	} else if strings.HasSuffix(age, "m") {
		// Minutes
		minutes, err := strconv.Atoi(strings.TrimSuffix(age, "m"))
		if err != nil {
			return nil, fmt.Errorf("invalid minute format: %s", age)
		}
		duration = time.Duration(minutes) * time.Minute
	} else {
		return nil, fmt.Errorf("unsupported age format: %s (use h for hours, d for days, m for minutes)", age)
	}

	// Calculate the cutoff time
	cutoff := time.Now().Add(-duration)
	return &cutoff, nil
}
