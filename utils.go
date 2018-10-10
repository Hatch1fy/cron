package cron

import (
	"time"
)

func getDelay(reference time.Time) (delay time.Duration) {
	// Set duration as seconds of the reference time
	duration := time.Second * time.Duration(reference.Second())
	// Increment duration by minutes of the reference time
	duration += time.Minute * time.Duration(reference.Minute())
	// Increment duration by hours of the reference time
	duration += time.Hour * time.Duration(reference.Hour())

	// Get the current time (in the reference time's location)
	now := time.Now().In(reference.Location())
	// Get the start of the current day
	start := GetStartOfDay(now)
	// Get the target time for the current day
	target := start.Add(duration)

	// Ensure target has not already occurred for today
	if target.Before(now) {
		// Target already occurred today, set target for tomorrow
		target = target.AddDate(0, 0, 1)
	}

	// Delta between target and current time (in seconds)
	delta := target.Unix() - now.Unix()

	// Return the delay as our delta converted to time.Duration (seconds)
	return time.Second * time.Duration(delta)
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
