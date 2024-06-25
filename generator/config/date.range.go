package config

import (
	"time"
)

// DateRangeConfig defines the configuration for a date range for a dynamic struct.
type DateRangeConfig interface {

	// SetDateFormat sets the date format.
	SetDateFormat(format string)

	// GetDateFormat returns the date format.
	GetDateFormat() string

	// SetDateRange sets the date range.
	SetDateRange(start time.Time, end time.Time)

	// SetStartDate sets the start date for the date range.
	SetStartDate(start time.Time)

	// GetStartDate returns the start date for the date range.
	GetStartDate() time.Time

	// SetEndDate sets the end date for the date range.
	SetEndDate(end time.Time)

	// GetEndDate returns the end date for the date range.
	GetEndDate() time.Time
}
