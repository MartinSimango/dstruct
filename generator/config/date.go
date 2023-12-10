package config

import (
	"time"
)

type DateConfig interface {
	SetDateFormat(format string)
	GetDateFormat() string
	SetDateRange(start time.Time, end time.Time)
	SetDateStart(start time.Time)
	GetDateStart() time.Time
	SetDateEnd(end time.Time)
	GetDateEnd() time.Time
}

type DateConfigImpl struct {
	dateFormat string
	dateStart  time.Time
	dateEnd    time.Time
}

var _ DateConfig = &DateConfigImpl{}

// SetDateStart implements DateConfig.
func (dc *DateConfigImpl) SetDateStart(start time.Time) {
	dc.SetDateRange(start, dc.dateEnd)
}

// SetDateEnd implements DateConfig.
func (dc *DateConfigImpl) SetDateEnd(end time.Time) {
	dc.SetDateRange(dc.dateStart, end)
}

// SetDateFormat implements DateConfig.
func (dc *DateConfigImpl) SetDateFormat(format string) {
	dc.dateFormat = format
}

// SetDateRange implements DateConfig.
func (dc *DateConfigImpl) SetDateRange(start time.Time, end time.Time) {
	if start.Before(end) {
		dc.dateStart = start
		dc.dateEnd = end
	}
}

// GetDateStart implements DateConfig.
func (dc *DateConfigImpl) GetDateStart() time.Time {
	return dc.dateStart
}

// GetDateEnd implements DateConfig.
func (dc *DateConfigImpl) GetDateEnd() time.Time {
	return dc.dateEnd
}

// GetDateFormat implements DateConfig.
func (dc *DateConfigImpl) GetDateFormat() string {
	return dc.dateFormat
}

func NewDateConfig() *DateConfigImpl {
	return &DateConfigImpl{
		dateFormat: time.RFC3339,
		dateStart:  time.Now().Truncate(24 * time.Hour),
		dateEnd:    time.Now().Truncate(24 * time.Hour).Add(24*time.Hour - time.Nanosecond),
	}
}
