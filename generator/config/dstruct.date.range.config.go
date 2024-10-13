package config

import "time"

// DstructDateRangeConfig implements the DateRangeConfig interface.
type DstructDateRangeConfig struct {
	dateFormat string
	startDate  time.Time
	endDate    time.Time
}

var _ DateRangeConfig = &DstructDateRangeConfig{}

// SetStartDate implements DateConfig.SetStartDate.
func (dc *DstructDateRangeConfig) SetStartDate(start time.Time) {
	dc.SetDateRange(start, dc.endDate)
}

// SetEndDate implements DateConfig.SetEndDate.
func (dc *DstructDateRangeConfig) SetEndDate(end time.Time) {
	dc.SetDateRange(dc.startDate, end)
}

// SetDateFormat implements DateConfig.SetDateFormat.
func (dc *DstructDateRangeConfig) SetDateFormat(format string) {
	dc.dateFormat = format
}

// SetDateRange implements DateConfig.SetDateRange.
func (dc *DstructDateRangeConfig) SetDateRange(start time.Time, end time.Time) {
	if start.Before(end) {
		dc.startDate = start
		dc.endDate = end
	}
}

// GetStartDate implements DateConfig.GetStartDate.
func (dc *DstructDateRangeConfig) GetStartDate() time.Time {
	return dc.startDate
}

// GetEndDate implements DateConfig.GetEndDate.
func (dc *DstructDateRangeConfig) GetEndDate() time.Time {
	return dc.endDate
}

// GetDateFormat implements DateConfig.GetDateFormat.
func (dc *DstructDateRangeConfig) GetDateFormat() string {
	return dc.dateFormat
}

// NewDstructDateRangeConfig is a constructor for DstructDateRangeConfig.
// It initializes the date format to RFC3339, the start date to the current date truncated to the day and the end date to the next day.
func NewDstructDateRangeConfig() *DstructDateRangeConfig {
	return &DstructDateRangeConfig{
		dateFormat: time.RFC3339,
		startDate:  time.Now().Truncate(24 * time.Hour),
		endDate:    time.Now().Truncate(24 * time.Hour).Add(24*time.Hour - time.Nanosecond),
	}
}
