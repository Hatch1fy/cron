# Cron
Cron is a simple cronjob-like library

## Usage
```go
package main

import (
	"fmt"
	"time"

	"github.com/Hatch1fy/cron"
)

func main() {
	// Initialize new instance of cron
	// After three seconds, display a success message
	cron.New(printSuccess).After(time.Second * 3)

	// Every one minute, print the current time
	cron.New(printTime).Every(time.Minute)

	// Call empty select so we can keep the service open indefinitely
	select {}
}

func printSuccess() {
	fmt.Println("The server has been successfully started!")
}

func printTime() {
	fmt.Println("The current time is:", time.Now())
}

```