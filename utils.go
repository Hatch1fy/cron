package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const timeFmt = "15:04"

func parseQuery(query string) (descriptor, reference string, value int64, loc *time.Location, err error) {
	// Split query by space
	spl := strings.Split(query, " ")
	// Ensure query has two or three parts
	if len(spl) != 2 && len(spl) != 3 {
		err = fmt.Errorf("invalid query, expected %d statements and received %d", 2, len(spl))
		return
	}

	// Validate the descriptor
	if err = validateDescriptor(spl[0]); err != nil {
		return
	}

	// Validate the reference and reference value
	if reference, value, err = validateReference(spl[1]); err != nil {
		return
	}

	if reference == "time" && len(spl) == 3 {
		// Validate the timezone
		if loc, err = validateTimezone(spl[2]); err != nil {
			return
		}
	}

	descriptor = spl[0]
	return
}

func validateDescriptor(descriptor string) (err error) {
	// Ensure our descriptor is either "every" or "@""
	switch descriptor {
	case "every":
	case "@":

	default:
		return fmt.Errorf("invalid descriptor, expected \"@\" or \"every\" and received \"%s\"", descriptor)
	}

	return
}

func validateReference(reference string) (set string, value int64, err error) {
	refBytes := []byte(reference)

	switch {
	case reference == "midnight":
		// Midnight is just shorthand for 00:00
		reference = "00:00"
		set = "time"
	case isValidTime(reference):
		set = "time"
	case millisecondRegEx.Match(refBytes):
		set = "ms"
	case secondRegEx.Match(refBytes):
		set = "s"
	case minuteRegEx.Match(refBytes):
		set = "m"
	case hourRegEx.Match(refBytes):
		set = "h"

	default:
		err = fmt.Errorf("invalid referenced, expected a time format or \"midnight\" and received \"%s\"", reference)
		return
	}

	if set != "time" {
		// Set is not a time-based value, parse the int and return
		value, err = strconv.ParseInt(reference[:len(reference)-len(set)], 10, 64)
		return
	}

	var t time.Time
	// Parse the reference as our time format (24 hour HH:MM)
	if t, err = time.Parse(timeFmt, reference); err != nil {
		return
	}

	// Set initial minutes value as the number of hours for the time multiplied by 60
	minutes := t.Hour() * 60
	// Increment minutes by the number of minutes for the time
	minutes += t.Minute()
	// Set value as the minutes value converted to int64
	value = int64(minutes)
	return
}

func validateTimezone(timezone string) (loc *time.Location, err error) {
	if len(timezone) == 0 {
		return
	}

	var offset int64
	// Parse the offset as hours
	if offset, err = strconv.ParseInt(timezone, 10, 64); err != nil {
		return
	}

	// Hours to minutes
	offset *= 60

	// Minutes to seconds
	offset *= 60

	// Set location as the timezone
	loc = time.FixedZone(timezone, int(offset))
	return
}

func isValidTime(reference string) (ok bool) {
	if _, err := time.Parse(timeFmt, reference); err == nil {
		// No errors parsing, reference is valid
		ok = true
	}

	return
}

func getTimeDuration(value int64) (dur time.Duration) {
	now := time.Now()
	start := GetStartOfDay(now)
	minutes := time.Minute * time.Duration(value)
	// Set the target as the start of te
	target := start.Add(minutes)

	if target.Before(now) {
		// Target already occurred today, set target for tomorrow
		target = target.AddDate(0, 0, 1)
	}

	// Set duration to the delta between the target and now
	dur = time.Duration(target.Unix() - now.Unix())
	// Convert duration to seconds
	dur *= time.Second
	return
}

// GetNextDay will get the next day at 00:00
func GetNextDay(current time.Time) (next time.Time) {
	year := current.Year()
	month := current.Month()
	day := current.Day() + 1
	loc := current.Location()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

// GetStartOfDay will get the current day at 00:00
func GetStartOfDay(current time.Time) (start time.Time) {
	year := current.Year()
	month := current.Month()
	day := current.Day()
	loc := current.Location()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}
