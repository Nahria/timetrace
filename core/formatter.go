package core

import (
	"errors"
	"fmt"
	"time"
)

// Formatter represents a date- and time formatter. It provides all displayed
// date- and time layouts and is capable of parsing those layouts.
type Formatter struct {
	use12Hours      bool
	useDecimalHours string
}

const dateLayout = "2006-01-02"

// ParseInterval returns the start and end time of an interval using aliases.
func (f *Formatter) ParseInterval(interval string) (time.Time, time.Time, error) {
	now := time.Now()

	if interval == "last-month" {
		start := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local)
		end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, -1, time.Local)
		return start, end, nil
	}

	if interval == "current-month" {
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		end := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, -1, time.Local)
		return start, end, nil
	}

	if interval == "current-year" {
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
		end := time.Date(now.Year()+1, 1, 1, 0, 0, 0, -1, time.Local)
		return start, end, nil
	}

	if interval == "last-year" {
		start := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.Local)
		end := time.Date(now.Year(), 1, 1, 0, 0, 0, -1, time.Local)
		return start, end, nil
	}

	return time.Time{}, time.Time{}, errors.New("interval value not supported")
}

// ParseDate parses a date from an input string in the form YYYY-MM-DD. It also
// supports the `today` and `yesterday` aliases for convenience.
func (f *Formatter) ParseDate(input string) (time.Time, error) {
	if input == "today" {
		return time.Now(), nil
	}
	if input == "yesterday" {
		yesterday := time.Now().AddDate(0, 0, -1)
		return yesterday, nil
	}

	if input == "current-month" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local), nil
	}

	if input == "last-month" {
		now := time.Now()
		return time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local), nil
	}

	date, err := time.Parse(dateLayout, input)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

// ParseTime parses a time from an input string in the configured timeLayout
func (f *Formatter) ParseTime(input string) (time.Time, error) {
	date, err := time.Parse(f.timeLayout(), input)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

// CombineDateAndTime takes a date and a time and combines them to the time
// struct that represents the given time on the given day
func (f *Formatter) CombineDateAndTime(d, t time.Time) time.Time {
	year, month, day := d.Date()
	return time.Date(year, month, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
}

const (
	defaultTimeLayout        = "15:04"
	default12HoursTimeLayout = "03:04PM"
	defaultDatelayoutPretty  = "Monday, 02. January 2006"
	defaultDurationLayout    = time.Minute
)

func (f *Formatter) timeLayout() string {
	if f.use12Hours {
		return default12HoursTimeLayout
	}
	return defaultTimeLayout
}

func (f *Formatter) durationLayout() time.Duration {
	return defaultDurationLayout
}

func (f *Formatter) TimeString(input time.Time) string {
	return input.Format(f.timeLayout())
}

func (f *Formatter) DurationString(input time.Duration) string {
	return input.Round(f.durationLayout()).String()
}

// PrettyDateString returns a nice representation of a given time
// example: Mon 31. May 2021
func (f *Formatter) PrettyDateString(input time.Time) string {
	return input.Format(defaultDatelayoutPretty)
}

const (
	defaultRecordKeyLayout        = "2006-01-02-15-04"
	default12HoursRecordKeyLayout = "2006-01-02-03-04PM"
)

func (f *Formatter) RecordKeyLayout() string {
	if f.use12Hours {
		return default12HoursRecordKeyLayout
	}
	return defaultRecordKeyLayout
}

// ParseRecordKey parses an input string in the form 2006-01-02-15-04 or
// 2006-01-02-03-04PM depending on the use12hours setting.
func (f *Formatter) ParseRecordKey(key string) (time.Time, error) {
	return time.Parse(f.RecordKeyLayout(), key)
}

func (f *Formatter) RecordKey(record *Record) string {
	return record.Start.Format(f.RecordKeyLayout())
}

// formatDuration formats the passed duration into a string.
// The format is determined by value of UseDecimalHours in config file and
// will be "8h 24min", "8.4h, or "8h 24min 8.4h".
// seconds information is ignored.
func (f *Formatter) FormatDuration(duration time.Duration) string {

	hours := int64(duration.Hours())
	minutes := int64(duration.Minutes()) % 60
	dec := duration.Minutes() / 60
	var response string
	switch f.useDecimalHours {
	case "Both":
		response = fmt.Sprintf("%dh %dmin %.1fh", hours, minutes, dec)
	case "On":
		response = fmt.Sprintf("%.1fh", dec)
	case "Off":
		response = fmt.Sprintf("%dh %dmin", hours, minutes)
	default:
		response = fmt.Sprintf("%dh %dmin", hours, minutes)
	}
	return response
}

func (f *Formatter) FormatTags(tags []string) string {
	var result string

	for i, t := range tags {
		result += t
		if i < len(tags)-1 {
			result += ", "
		}
	}

	return result
}
