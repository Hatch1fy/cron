# Cron
Cron is a simple cronjob-like library

## Usage
```go
package main

import (
	"fmt"
	"time"

	"github.com/hatchify/cron"
)

func main() {
	// Initialize new instance of cron
	// After three seconds, display a success message
	cron.New(printSuccess).After(time.Second * 3)

	// Every one minute, print the current time
	cron.New(printTime).Every(time.Minute)

	// Every day at lunchtime, print that it's time to eat
	cron.New(printLunch).EveryAt(getLunchtime())

	// Call empty select so we can keep the service open indefinitely
	select {}
}

func getLunchtime() (lunch time.Time) {
	now := time.Now()
	lunch = time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	return
}

func printSuccess() {
	fmt.Println("The server has been successfully started!")
}

func printTime() {
	fmt.Println("The current time is:", time.Now())
}

func printLunch() {
	fmt.Println("It's time to eat lunch!")
}

```